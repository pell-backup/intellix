package server

import (
	"context"
	"cosmossdk.io/log"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"sync"
)

type PriceInfo struct {
	DataSource string
	Price      float64
}

const (
	dataSourceCoinbase = "coinbase"
	dataSourceBinance  = "binance"
)

func fetchRawPrices(ctx context.Context, logger log.Logger, baseSymbol, quoteSymbol string) (map[string]*big.Int, error) {
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

	var prices = map[string]*big.Int{}
	for price := range priceChan {
		prices[price.DataSource] = big.NewInt(int64(math.Round(price.Price * 100000000)))
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

func (s *CoinbaseFetchPriceService) fetchCoinPrice(baseSymbol, quote string, wg *sync.WaitGroup, priceChan chan<- *PriceInfo) {
	defer wg.Done()
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/spot", baseSymbol, quote)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching data from Coinbase:", err.Error())
		return
	}
	defer resp.Body.Close()

	type CoinbaseResponse struct {
		Data struct {
			Base     string `json:"base"`
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		} `json:"data"`
	}

	var coinbaseResp CoinbaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&coinbaseResp); err != nil {
		s.logger.Error("Error decoding JSON response from Coinbase:", err.Error())
		return
	}

	price, err := strconv.ParseFloat(coinbaseResp.Data.Amount, 64)
	if err != nil {
		s.logger.Error("Error parsing price from Coinbase failed")
		return
	}

	priceChan <- &PriceInfo{DataSource: dataSourceCoinbase, Price: price}
}

type BinanceFetchPriceService struct {
	logger log.Logger
}

func (s *BinanceFetchPriceService) fetchCoinPrice(baseSymbol, quote string, wg *sync.WaitGroup, priceChan chan<- *PriceInfo) {
	defer wg.Done()
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", baseSymbol, quote)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching data from Binance:", err.Error())
		return
	}
	defer resp.Body.Close()

	type BinanceResponse struct {
		Price string `json:"price"`
	}

	var binanceResp BinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&binanceResp); err != nil {
		s.logger.Error("Error decoding JSON response from Binance:", err.Error())
		return
	}

	price, err := strconv.ParseFloat(binanceResp.Price, 64)
	if err != nil {
		fmt.Println("Error parsing price from Binance:", err)
		return
	}

	priceChan <- &PriceInfo{DataSource: dataSourceBinance, Price: price}
}
