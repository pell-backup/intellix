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

type GateFetchPriceService struct {
	logger log.Logger
}

func (s *GateFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	s.logger.Info("Gate.io data source activated", "status", "initializing")

	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
		s.logger.Info("Gate.io symbol conversion applied", "original_base", strings.ToUpper(base), "converted_base", tick)
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
		s.logger.Info("Gate.io symbol conversion applied", "original_quote", strings.ToUpper(quote), "converted_quote", tick)
	}

	// Gate.io uses a format like BTC_USDT
	var url = fmt.Sprintf("https://api.gateio.ws/api/v4/spot/tickers?currency_pair=%s_%s", base, quote)
	s.logger.Info("Gate.io price request", "base", base, "quote", quote, "url", url)

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Gate.io price fetch failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	var gateResp []struct {
		Last string `json:"last"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gateResp); err != nil {
		s.logger.Error("Gate.io response decode failed", "error", err)
		return err
	}

	// Check if response contains data
	if len(gateResp) == 0 {
		s.logger.Error("Gate.io empty response", "data_length", len(gateResp))
		return fmt.Errorf("empty response from Gate.io")
	}

	price, err := strconv.ParseFloat(gateResp[0].Last, 64)
	if err != nil {
		s.logger.Error("Gate.io price parsing failed", "error", err)
		return err
	}

	s.logger.Info("Gate.io price fetch successful",
		"base", base,
		"quote", quote,
		"price", price,
		"status", "success")

	// Convert to Dec type
	priceStr := gateResp[0].Last

	// Truncate decimal places if needed
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("Gate.io price conversion to Dec failed", "error", err, "price", priceStr)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceGate, Price: dec}

	s.logger.Info("Gate.io data source completed", "status", "success")
	return nil
}
