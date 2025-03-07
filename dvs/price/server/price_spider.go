package server

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"golang.org/x/sync/errgroup"
)

type PriceInfo struct {
	DataSource string
	Price      math.LegacyDec
}

const (
	dataSourceCoinbase = "coinbase"
	dataSourceBinance  = "binance"
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

func fetchRawPrices(ctx context.Context, logger log.Logger, baseSymbol, quoteSymbol string, tickConverter PriceTickConverterByDataSource) (map[string]math.LegacyDec, error) {
	// TODO: configurable
	var fetchPriceIfs = map[string]FetchPriceServiceIF{
		dataSourceCoinbase: &CoinbaseFetchPriceService{logger: logger},
		dataSourceBinance:  &BinanceFetchPriceService{logger: logger},
	}

	priceChan := make(chan *PriceInfo, len(fetchPriceIfs))
	g, _ := errgroup.WithContext(ctx)

	for dataSource, fetchPriceIf := range fetchPriceIfs {
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
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}

	close(priceChan)

	var prices = map[string]math.LegacyDec{}
	for price := range priceChan {
		prices[price.DataSource] = price.Price
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("failed to fetch prices from exchanges")
	}

	return prices, nil
}
