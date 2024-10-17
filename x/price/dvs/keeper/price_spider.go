package keeper

import (
	"cosmossdk.io/log"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"net/http"
	"sync"
)

type PriceInfo struct {
	DataSource string
	Price      *big.Int
}

const (
	dataSourceCoinbase = "coinbase"
	dataSourceBinance  = "binance"
)

func fetchRawPrices(ctx sdk.Context, logger log.Logger, baseSymbol, quoteSymbol string) (map[string]*big.Int, error) {
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

	price, ok := new(big.Int).SetString(coinbaseResp.Data.Amount, 10)
	if !ok {
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

	price, ok := new(big.Int).SetString(binanceResp.Price, 10)
	if !ok {
		s.logger.Error("Error parsing price from Binance failed")
		return
	}

	priceChan <- &PriceInfo{DataSource: dataSourceBinance, Price: price}
}
