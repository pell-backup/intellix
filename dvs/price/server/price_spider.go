package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"cosmossdk.io/math"

	"github.com/0xPellNetwork/pelldvs/libs/log"
)

type PriceInfo struct {
	DataSource string
	Price      math.LegacyDec
}

const (
	dataSourceCoinbase = "coinbase"
	dataSourceBinance  = "binance"
)

func fetchRawPrices(ctx context.Context, logger log.Logger, baseSymbol, quoteSymbol string) (map[string]math.LegacyDec, error) {
	var wg sync.WaitGroup

	// TODO: configurable
	var fetchPriceIfs = map[string]FetchPriceServiceIF{
		dataSourceCoinbase: &CoinbaseFetchPriceService{logger: logger},
		dataSourceBinance:  &BinanceFetchPriceService{logger: logger},
	}
	priceChan := make(chan *PriceInfo, len(fetchPriceIfs))
	wg.Add(2)
	for _, fetchPriceIf := range fetchPriceIfs {
		go fetchPriceIf.fetchCoinPrice(baseSymbol, quoteSymbol, &wg, priceChan)
	}
	wg.Wait()

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

type FetchPriceServiceIF interface {
	fetchCoinPrice(base, quote string, wg *sync.WaitGroup, priceChan chan<- *PriceInfo)
}

type CoinbaseFetchPriceService struct {
	logger log.Logger
}

func (s *CoinbaseFetchPriceService) fetchCoinPrice(base, quote string, wg *sync.WaitGroup, priceChan chan<- *PriceInfo) {
	defer wg.Done()

	var url = fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/spot", base, quote)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching price from Coinbase", "error", err)
		return
	}
	defer resp.Body.Close()

	var coinbaseResp struct {
		Data struct {
			Amount string `json:"amount"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&coinbaseResp); err != nil {
		s.logger.Error("Error decoding Coinbase response", "error", err)
		return
	}

	price, err := strconv.ParseFloat(coinbaseResp.Data.Amount, 64)
	if err != nil {
		s.logger.Error("Error parsing price from Coinbase failed")
		return
	}
	s.logger.Info("Fetched price from Coinbase", "base", base, "quote", quote, "price", price)

	dec, err := math.LegacyNewDecFromStr(coinbaseResp.Data.Amount)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err)
		return
	}
	priceChan <- &PriceInfo{DataSource: dataSourceCoinbase, Price: dec}
}

type BinanceFetchPriceService struct {
	logger log.Logger
}

func (s *BinanceFetchPriceService) fetchCoinPrice(base, quote string, wg *sync.WaitGroup, priceChan chan<- *PriceInfo) {
	defer wg.Done()

	var url = fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", base, quote)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching price from Binance", "error", err)
		return
	}
	defer resp.Body.Close()

	var binanceResp struct {
		Price string `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		s.logger.Error("Error decoding Binance response", "error", err)
		return
	}

	price, err := strconv.ParseFloat(binanceResp.Price, 64)
	if err != nil {
		s.logger.Error("Error parsing price from Binance failed")
		return
	}

	s.logger.Info("Fetched price from Binance", "base", base, "quote", quote, "price", price)

	dec, err := math.LegacyNewDecFromStr(binanceResp.Price)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err)
		return
	}
	priceChan <- &PriceInfo{DataSource: dataSourceBinance, Price: dec}
}
