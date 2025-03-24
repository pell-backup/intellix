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

type OKXFetchPriceService struct {
	logger log.Logger
}

func (s *OKXFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	s.logger.Info("OKX data source activated", "status", "initializing")

	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
		s.logger.Info("OKX symbol conversion applied", "original_base", strings.ToUpper(base), "converted_base", tick)
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
		s.logger.Info("OKX symbol conversion applied", "original_quote", strings.ToUpper(quote), "converted_quote", tick)
	}

	// OKX uses a format like BTC-USDT
	var url = fmt.Sprintf("https://www.okx.com/api/v5/market/ticker?instId=%s-%s", base, quote)
	s.logger.Info("OKX price request", "base", base, "quote", quote, "url", url)

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("OKX price fetch failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	var okxResp struct {
		Code string `json:"code"`
		Data []struct {
			Last string `json:"last"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&okxResp); err != nil {
		s.logger.Error("OKX response decode failed", "error", err)
		return err
	}

	// Check if response is successful and contains data
	if okxResp.Code != "0" || len(okxResp.Data) == 0 {
		s.logger.Error("OKX invalid response", "code", okxResp.Code, "data_length", len(okxResp.Data))
		return fmt.Errorf("invalid response from OKX: code %s", okxResp.Code)
	}

	price, err := strconv.ParseFloat(okxResp.Data[0].Last, 64)
	if err != nil {
		s.logger.Error("OKX price parsing failed", "error", err)
		return err
	}

	s.logger.Info("OKX price fetch successful",
		"base", base,
		"quote", quote,
		"price", price,
		"status", "success")

	// Convert to Dec type
	priceStr := okxResp.Data[0].Last

	// Truncate decimal places if needed
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("OKX price conversion to Dec failed", "error", err, "price", priceStr)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceOKX, Price: dec}

	s.logger.Info("OKX data source completed", "status", "success")
	return nil
}
