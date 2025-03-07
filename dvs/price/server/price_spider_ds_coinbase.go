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

type CoinbaseFetchPriceService struct {
	logger log.Logger
}

func (s *CoinbaseFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
	}

	var url = fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/spot", base, quote)
	s.logger.Info("start fetch price from Coinbase", "base", base, "quote", quote, "url", url)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching price from Coinbase", "error", err)
		return err
	}
	defer resp.Body.Close()

	var coinbaseResp struct {
		Data struct {
			Amount string `json:"amount"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&coinbaseResp); err != nil {
		s.logger.Error("Error decoding Coinbase response", "error", err)
		return err
	}

	price, err := strconv.ParseFloat(coinbaseResp.Data.Amount, 64)
	if err != nil {
		s.logger.Error("Error parsing price from Coinbase failed")
		return err
	}
	s.logger.Info("Fetched price from Coinbase", "base", base, "quote", quote, "price", price)

	dec, err := math.LegacyNewDecFromStr(coinbaseResp.Data.Amount)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceCoinbase, Price: dec}
	return nil
}
