package server

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"golang.org/x/sync/errgroup"

	"intellix/dvs/price/types"
)

type PriceInfo struct {
	DataSource string
	Price      math.LegacyDec
}

const (
	dataSourceCoinbase      = "coinbase"
	dataSourceBinance       = "binance"
	dataSourceCoinMarketCap = "coinmarketcap"
	dataSourceOKX           = "okx"
	dataSourceGate          = "gate"
	dataSourceEasyMoney     = "easymoney"
	dataSourceITick         = "itick"
)

type PriceTickConverter map[string]string

type PriceTickConverterByDataSource map[string]PriceTickConverter

func ToPriceTickConverterByDataSource(tickConverterConfig map[string]map[string]string) PriceTickConverterByDataSource {
	var tickConverterByDataSource = make(PriceTickConverterByDataSource)
	for dataSource, tickConverter := range tickConverterConfig {
		for base, quote := range tickConverter {
			tickConverter[strings.ToUpper(base)] = strings.ToUpper(quote)
		}
		tickConverterByDataSource[dataSource] = tickConverter
	}
	return tickConverterByDataSource
}

type FetchPriceServiceIF interface {
	fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error
}

// Helper function to get all keys from a map and sort them
func getMapKeys(m map[string]FetchPriceServiceIF) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TruncatePriceDecimal ensures price decimal places don't exceed 18 digits
func TruncatePriceDecimal(priceStr string, logger log.Logger) string {
	// Check if precision exceeds 18 decimal places
	parts := strings.Split(priceStr, ".")
	if len(parts) == 2 && len(parts[1]) > 18 {
		logger.Info("Truncating price decimal places to 18 digits",
			"original", priceStr,
			"decimal_places", len(parts[1]))
		// Truncate to 18 decimal places
		priceStr = parts[0] + "." + parts[1][:18]
		logger.Info("Truncated price", "new_price", priceStr)
	}
	return priceStr
}

func fetchRawPrices(ctx context.Context, logger log.Logger, baseSymbol, quoteSymbol string, tickConverter PriceTickConverterByDataSource,
	apiKey map[string]string, symbolList map[string]map[string]string) (map[string]math.LegacyDec, error) {

	// Record enabled data sources
	logger.Info("Initializing price data sources")

	fetchPriceIfs := make(map[string]FetchPriceServiceIF)

	// This is a workaround for the fact that some exchanges use different base symbols
	baseSymbolConfig := strings.ReplaceAll(baseSymbol, ".", "_")
	baseSymbolConfig = strings.ToLower(baseSymbolConfig)

	// Check if the base symbol is a crypto symbol
	if _, ok := symbolList[types.Crypto][baseSymbolConfig]; ok {
		fetchPriceIfs[dataSourceCoinbase] = &CoinbaseFetchPriceService{logger: logger}
		fetchPriceIfs[dataSourceBinance] = &BinanceFetchPriceService{logger: logger}
		fetchPriceIfs[dataSourceOKX] = &OKXFetchPriceService{logger: logger}
		fetchPriceIfs[dataSourceGate] = &GateFetchPriceService{logger: logger}

		if key := apiKey[types.CoinMarketCap]; key != "" {
			fetchPriceIfs[dataSourceCoinMarketCap] = &CMCFetchPriceService{logger: logger, apiKey: key}
		} else {
			logger.Error("CoinMarketCap API key is empty, skipping this data source")
		}
	}

	// Check if the quote symbol is in A-shares
	if _, ok := symbolList[types.AShares][baseSymbolConfig]; ok {
		fetchPriceIfs[dataSourceEasyMoney] = &EasyMoneyFetchPriceService{logger: logger}
	}

	// Check if the quote symbol is in US stocks
	if _, ok := symbolList[types.USStocks][baseSymbolConfig]; ok {
		// This assumes you want to **add** this source, not override the whole map
		fetchPriceIfs[dataSourceCoinbase] = &CoinbaseFetchPriceService{logger: logger}
	}

	// Record enabled data sources
	logger.Info("Enabled price data sources",
		"count", len(fetchPriceIfs),
		"sources", fmt.Sprintf("%v", getMapKeys(fetchPriceIfs)))

	priceChan := make(chan *PriceInfo, len(fetchPriceIfs))
	g, _ := errgroup.WithContext(ctx)

	for dataSource, fetchPriceIf := range fetchPriceIfs {
		dataSource := dataSource // Create a copy to avoid closure issues
		fetchPriceIf := fetchPriceIf
		g.Go(func() error {
			logger.Info("fetching coin price",
				"baseSymbol", baseSymbol,
				"quoteSymbol", quoteSymbol,
				"dataSource", dataSource,
			)
			return fetchPriceIf.fetchCoinPrice(baseSymbol, quoteSymbol, tickConverter[dataSource], priceChan)
		})
	}
	if err := g.Wait(); err != nil {
		logger.Error("fetching coin price failed", "error", err)
	}

	close(priceChan)

	var prices = map[string]math.LegacyDec{}
	for price := range priceChan {
		prices[price.DataSource] = price.Price
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("failed to fetch prices from exchanges")
	}

	// Record prices from each data source
	logger.Info("Price data summary", "total_sources", len(prices))
	for source, price := range prices {
		logger.Info("Price data from source",
			"source", source,
			"price", price.String(),
			"base", baseSymbol,
			"quote", quoteSymbol)
	}

	return prices, nil
}
