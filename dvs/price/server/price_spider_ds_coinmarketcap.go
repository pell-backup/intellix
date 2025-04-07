package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

// CMCFetchPriceService implements the CoinMarketCap data source
type CMCFetchPriceService struct {
	logger log.Logger
	apiKey string
}

func (s *CMCFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	s.logger.Info("CoinMarketCap data source activated", "status", "initializing")

	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
		s.logger.Info("CoinMarketCap symbol conversion applied", "original_base", strings.ToUpper(base), "converted_base", tick)
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
		s.logger.Info("CoinMarketCap symbol conversion applied", "original_quote", strings.ToUpper(quote), "converted_quote", tick)
	}

	// CoinMarketCap API requires getting the cryptocurrency ID first
	symbolUrl := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s&convert=%s", base, quote)
	s.logger.Info("CoinMarketCap price request", "base", base, "quote", quote, "url", symbolUrl)

	req, err := http.NewRequest("GET", symbolUrl, nil)
	if err != nil {
		s.logger.Error("CoinMarketCap request creation failed", "error", err)
		return err
	}

	// Set API key
	req.Header.Set("X-CMC_PRO_API_KEY", s.apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("CoinMarketCap price fetch failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("CoinMarketCap API returned non-OK status", "status", resp.StatusCode)
		return fmt.Errorf("CoinMarketCap API returned status %d", resp.StatusCode)
	}

	// Parse response
	var cmcResp struct {
		Data map[string]struct {
			Quote map[string]struct {
				Price float64 `json:"price"`
			} `json:"quote"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cmcResp); err != nil {
		s.logger.Error("CoinMarketCap response decode failed", "error", err)
		return err
	}

	// Check if there is data
	if len(cmcResp.Data) == 0 {
		s.logger.Error("CoinMarketCap empty response", "data_length", len(cmcResp.Data))
		return fmt.Errorf("empty response from CoinMarketCap")
	}

	// Get price
	coinData, ok := cmcResp.Data[base]
	if !ok {
		s.logger.Error("CoinMarketCap base symbol not found in response", "base", base)
		return fmt.Errorf("base symbol %s not found in CoinMarketCap response", base)
	}

	quoteData, ok := coinData.Quote[quote]
	if !ok {
		s.logger.Error("CoinMarketCap quote symbol not found in response", "quote", quote)
		return fmt.Errorf("quote symbol %s not found in CoinMarketCap response", quote)
	}

	price := quoteData.Price
	s.logger.Info("CoinMarketCap price fetch successful",
		"base", base,
		"quote", quote,
		"price", price,
		"status", "success")

	// Convert to Dec type
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)

	// Truncate decimal places if needed
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("CoinMarketCap price conversion to Dec failed", "error", err, "price", priceStr)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceCoinMarketCap, Price: dec}

	s.logger.Info("CoinMarketCap data source completed", "status", "success")
	return nil
}
