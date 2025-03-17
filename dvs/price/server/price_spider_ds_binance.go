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

type BinanceFetchPriceService struct {
	logger log.Logger
}

func (s *BinanceFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
	}

	var url = fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", base, quote)
	s.logger.Info("start fetch price from Binance", "base", base, "quote", quote, "url", url)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching price from Binance", "error", err)
		return err
	}
	defer resp.Body.Close()

	var binanceResp struct {
		Price string `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		s.logger.Error("Error decoding Binance response", "error", err)
		return err
	}

	price, err := strconv.ParseFloat(binanceResp.Price, 64)
	if err != nil {
		s.logger.Error("Error parsing price from Binance failed", "error", err)
		return err
	}

	s.logger.Info("Fetched price from Binance", "base", base, "quote", quote, "price", price)

	// Convert to Dec type
	priceStr := binanceResp.Price

	// Truncate decimal places if needed
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err, "price", priceStr)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceBinance, Price: dec}

	return nil
}
