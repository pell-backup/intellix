package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

// BinanceAdapter 实现 Binance 数据源
type BinanceAdapter struct {
	symbolMap map[string]string
	baseURL   string
	logger    log.Logger
}

func NewBinanceAdapter(logger log.Logger) *BinanceAdapter {
	return &BinanceAdapter{
		symbolMap: map[string]string{
			// 主流代币
			"BTC":  "BTCUSDT",
			"ETH":  "ETHUSDT",
			"BNB":  "BNBUSDT",
			"XRP":  "XRPUSDT",
			"XLM":  "XLMUSDT",
			"USDT": "USDTUSD",
			"USDC": "USDCUSDT",
			"DAI":  "DAIUSDT",

			// 迷因代币
			"DOGE":  "DOGEUSDT",
			"PEPE":  "PEPEUSDT",
			"CAT":   "CATUSDT",
			"GOAT":  "GOATUSDT",
			"DOGS":  "DOGSUSDT",
			"MEME":  "MEMEUSDT",
			"CATS":  "CATSUSDT",
			"LADYS": "LADYSUSDT",
			"PUSS":  "PUSSUSDT",

			// 新兴代币
			"TON":      "TONUSDT",
			"SAVM":     "SAVMUSDT",
			"ORDI":     "ORDIUSDT",
			"NEIROCTO": "NEIROCTOUSDT",
			"MUBI":     "MUBIUSDT",
			"BB":       "BBUSDT",
			"SUNDOG":   "SUNDOGUSDT",
			"BAN":      "BANUSDT",
			"MOODENG":  "MOODENGUSDT",
			"NEIRO":    "NEIROUSDT",
			"1000SATS": "1000SATSUSDT",
			"NEIROETH": "NEIROETHUSDT",
			"AUCTION":  "AUCTIONUSDT",
			"WIF":      "WIFUSDT",
			"BOME":     "BOMEUSDT",
			"ACT":      "ACTUSDT",
		},
		baseURL: "https://api.binance.com",
		logger:  logger,
	}
}

func (b *BinanceAdapter) Name() string {
	return "binance"
}

func (b *BinanceAdapter) SupportsToken(symbol string) bool {
	_, ok := b.symbolMap[symbol]
	return ok
}

func (b *BinanceAdapter) MapSymbol(symbol string) string {
	if mapped, ok := b.symbolMap[symbol]; ok {
		return mapped
	}
	return symbol + "USDT"
}

func (b *BinanceAdapter) BuildRequest(ctx context.Context, symbol string) (*http.Request, error) {
	mappedSymbol := b.MapSymbol(symbol)
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", b.baseURL, mappedSymbol)
	return http.NewRequestWithContext(ctx, "GET", url, nil)
}

func (b *BinanceAdapter) ParseResponse(resp *http.Response, symbol string) (*PriceData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Binance response: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Binance API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var binanceResp struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	if err := json.Unmarshal(body, &binanceResp); err != nil {
		return nil, fmt.Errorf("error decoding Binance response: %w", err)
	}

	price, err := math.LegacyNewDecFromStr(binanceResp.Price)
	if err != nil {
		return nil, fmt.Errorf("error parsing price from Binance: %w", err)
	}

	return &PriceData{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().Unix(),
		Source:    "binance",
		Raw:       binanceResp,
	}, nil
}

// CoinbaseAdapter 实现 Coinbase 数据源
type CoinbaseAdapter struct {
	symbolMap map[string]string
	baseURL   string
	logger    log.Logger
}

func NewCoinbaseAdapter(logger log.Logger) *CoinbaseAdapter {
	return &CoinbaseAdapter{
		symbolMap: map[string]string{
			// 主流代币
			"BTC":  "BTC-USD",
			"ETH":  "ETH-USD",
			"XRP":  "XRP-USD",
			"XLM":  "XLM-USD",
			"USDT": "USDT-USD",
			"USDC": "USDC-USD",
			"DAI":  "DAI-USD",
			"BNB":  "BNB-USD",

			// 迷因代币
			"DOGE": "DOGE-USD",
			"PEPE": "PEPE-USD",
			"MEME": "MEME-USD",

			// 新兴代币
			"TON":     "TON-USD",
			"ORDI":    "ORDI-USD",
			"AUCTION": "AUCTION-USD",
		},
		baseURL: "https://api.coinbase.com",
		logger:  logger,
	}
}

func (c *CoinbaseAdapter) Name() string {
	return "coinbase"
}

func (c *CoinbaseAdapter) SupportsToken(symbol string) bool {
	_, ok := c.symbolMap[symbol]
	return ok
}

func (c *CoinbaseAdapter) MapSymbol(symbol string) string {
	if mapped, ok := c.symbolMap[symbol]; ok {
		return mapped
	}
	return symbol + "-USD"
}

func (c *CoinbaseAdapter) BuildRequest(ctx context.Context, symbol string) (*http.Request, error) {
	mappedSymbol := c.MapSymbol(symbol)
	url := fmt.Sprintf("%s/v2/prices/%s/spot", c.baseURL, mappedSymbol)
	return http.NewRequestWithContext(ctx, "GET", url, nil)
}

func (c *CoinbaseAdapter) ParseResponse(resp *http.Response, symbol string) (*PriceData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Coinbase response: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Coinbase API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var coinbaseResp struct {
		Data struct {
			Base     string `json:"base"`
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &coinbaseResp); err != nil {
		return nil, fmt.Errorf("error decoding Coinbase response: %w", err)
	}

	price, err := math.LegacyNewDecFromStr(coinbaseResp.Data.Amount)
	if err != nil {
		return nil, fmt.Errorf("error parsing price from Coinbase: %w", err)
	}

	return &PriceData{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().Unix(),
		Source:    "coinbase",
		Raw:       coinbaseResp,
	}, nil
}

// CoinmarketcapAdapter 实现 Coinmarketcap 数据源
type CoinmarketcapAdapter struct {
	symbolMap map[string]string
	idMap     map[string]string
	baseURL   string
	apiKey    string
	logger    log.Logger
}

func NewCoinmarketcapAdapter(apiKey string, logger log.Logger) *CoinmarketcapAdapter {
	return &CoinmarketcapAdapter{
		symbolMap: map[string]string{
			// 主流代币
			"BTC":  "bitcoin",
			"ETH":  "ethereum",
			"BNB":  "binancecoin",
			"XRP":  "ripple",
			"XLM":  "stellar",
			"USDT": "tether",
			"USDC": "usd-coin",
			"DAI":  "multi-collateral-dai",

			// 迷因代币
			"DOGE":  "dogecoin",
			"PEPE":  "pepe",
			"CAT":   "cat-token",
			"GOAT":  "goat-token",
			"DOGS":  "dogs",
			"MEME":  "meme",
			"CATS":  "cats",
			"LADYS": "ladys",
			"PUSS":  "puss-token",

			// 新兴代币
			"TON":      "toncoin",
			"SAVM":     "savm",
			"ORDI":     "ordinals",
			"NEIROCTO": "neirocto",
			"MUBI":     "mubi",
			"BB":       "bb-token",
			"SUNDOG":   "sundog",
			"BAN":      "banano",
			"MOODENG":  "moodeng",
			"NEIRO":    "neiro",
			"1000SATS": "1000sats",
			"NEIROETH": "neiroeth",
			"AUCTION":  "auction",
			"WIF":      "dogwifhat",
			"BOME":     "book-of-meme",
			"ACT":      "act-token",
		},
		idMap: map[string]string{
			// 主流代币
			"BTC":  "1",
			"ETH":  "1027",
			"BNB":  "1839",
			"XRP":  "52",
			"XLM":  "512",
			"USDT": "825",
			"USDC": "3408",
			"DAI":  "4943",

			// 迷因代币
			"DOGE":  "74",
			"PEPE":  "24478",
			"CAT":   "8236",
			"GOAT":  "25422",
			"DOGS":  "27741",
			"MEME":  "25541",
			"CATS":  "28352",
			"LADYS": "25092",
			"PUSS":  "28282",

			// 新兴代币
			"TON":      "11419",
			"SAVM":     "28620",
			"ORDI":     "23095",
			"NEIROCTO": "28953",
			"MUBI":     "27210",
			"BB":       "28295",
			"SUNDOG":   "28473",
			"BAN":      "4704",
			"MOODENG":  "28601",
			"NEIRO":    "28602",
			"1000SATS": "28295",
			"NEIROETH": "28603",
			"AUCTION":  "7663",
			"WIF":      "24477",
			"BOME":     "28482",
			"ACT":      "28604",
		},
		baseURL: "https://pro-api.coinmarketcap.com",
		apiKey:  apiKey,
		logger:  logger,
	}
}

func (c *CoinmarketcapAdapter) Name() string {
	return "coinmarketcap"
}

func (c *CoinmarketcapAdapter) SupportsToken(symbol string) bool {
	_, ok := c.idMap[symbol]
	return ok
}

func (c *CoinmarketcapAdapter) BuildRequest(ctx context.Context, symbol string) (*http.Request, error) {
	id, ok := c.idMap[symbol]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found in Coinmarketcap ID map", symbol)
	}

	url := fmt.Sprintf("%s/v2/cryptocurrency/quotes/latest?id=%s", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加API密钥
	req.Header.Add("X-CMC_PRO_API_KEY", c.apiKey)
	req.Header.Add("Accept", "application/json")
	return req, nil
}

func (c *CoinmarketcapAdapter) ParseResponse(resp *http.Response, symbol string) (*PriceData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Coinmarketcap response: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Coinmarketcap API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var cmcResp struct {
		Data map[string]struct {
			Quote map[string]struct {
				Price float64 `json:"price"`
			} `json:"quote"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &cmcResp); err != nil {
		return nil, fmt.Errorf("error decoding Coinmarketcap response: %w", err)
	}

	id := c.idMap[symbol]
	data, ok := cmcResp.Data[id]
	if !ok {
		return nil, fmt.Errorf("data for ID %s not found in response", id)
	}

	usdQuote, ok := data.Quote["USD"]
	if !ok {
		return nil, fmt.Errorf("USD quote not found for symbol %s", symbol)
	}

	// 将浮点数转换为 LegacyDec
	priceStr := strconv.FormatFloat(usdQuote.Price, 'f', -1, 64)
	price, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		return nil, fmt.Errorf("error converting price to Dec: %w", err)
	}

	return &PriceData{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().Unix(),
		Source:    "coinmarketcap",
		Raw:       cmcResp,
	}, nil
}

// GateioAdapter 实现 Gate.io 数据源
type GateioAdapter struct {
	symbolMap map[string]string
	baseURL   string
	logger    log.Logger
}

func NewGateioAdapter(logger log.Logger) *GateioAdapter {
	return &GateioAdapter{
		symbolMap: map[string]string{
			// 主流代币
			"BTC":  "BTC_USDT",
			"ETH":  "ETH_USDT",
			"BNB":  "BNB_USDT",
			"XRP":  "XRP_USDT",
			"XLM":  "XLM_USDT",
			"USDT": "USDT_USD",
			"USDC": "USDC_USDT",
			"DAI":  "DAI_USDT",

			// 迷因代币
			"DOGE":  "DOGE_USDT",
			"PEPE":  "PEPE_USDT",
			"CAT":   "CAT_USDT",
			"GOAT":  "GOAT_USDT",
			"DOGS":  "DOGS_USDT",
			"MEME":  "MEME_USDT",
			"CATS":  "CATS_USDT",
			"LADYS": "LADYS_USDT",
			"PUSS":  "PUSS_USDT",

			// 新兴代币
			"TON":      "TON_USDT",
			"SAVM":     "SAVM_USDT",
			"ORDI":     "ORDI_USDT",
			"NEIROCTO": "NEIROCTO_USDT",
			"MUBI":     "MUBI_USDT",
			"BB":       "BB_USDT",
			"SUNDOG":   "SUNDOG_USDT",
			"BAN":      "BAN_USDT",
			"MOODENG":  "MOODENG_USDT",
			"NEIRO":    "NEIRO_USDT",
			"1000SATS": "1000SATS_USDT",
			"NEIROETH": "NEIROETH_USDT",
			"AUCTION":  "AUCTION_USDT",
			"WIF":      "WIF_USDT",
			"BOME":     "BOME_USDT",
			"ACT":      "ACT_USDT",
		},
		baseURL: "https://api.gateio.ws",
		logger:  logger,
	}
}

func (g *GateioAdapter) Name() string {
	return "gateio"
}

func (g *GateioAdapter) SupportsToken(symbol string) bool {
	_, ok := g.symbolMap[symbol]
	return ok
}

func (g *GateioAdapter) MapSymbol(symbol string) string {
	if mapped, ok := g.symbolMap[symbol]; ok {
		return mapped
	}
	return symbol + "_USDT"
}

func (g *GateioAdapter) BuildRequest(ctx context.Context, symbol string) (*http.Request, error) {
	mappedSymbol := g.MapSymbol(symbol)
	url := fmt.Sprintf("%s/api/v4/spot/tickers?currency_pair=%s", g.baseURL, mappedSymbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (g *GateioAdapter) ParseResponse(resp *http.Response, symbol string) (*PriceData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Gate.io response: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gate.io API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var gateioResp []struct {
		CurrencyPair     string `json:"currency_pair"`
		Last             string `json:"last"`
		LowestAsk        string `json:"lowest_ask"`
		HighestBid       string `json:"highest_bid"`
		ChangePercentage string `json:"change_percentage"`
		BaseVolume       string `json:"base_volume"`
		QuoteVolume      string `json:"quote_volume"`
		High24h          string `json:"high_24h"`
		Low24h           string `json:"low_24h"`
	}

	if err := json.Unmarshal(body, &gateioResp); err != nil {
		return nil, fmt.Errorf("error decoding Gate.io response: %w", err)
	}

	if len(gateioResp) == 0 {
		return nil, fmt.Errorf("empty response from Gate.io")
	}

	price, err := math.LegacyNewDecFromStr(gateioResp[0].Last)
	if err != nil {
		return nil, fmt.Errorf("error parsing price from Gate.io: %w", err)
	}

	return &PriceData{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().Unix(),
		Source:    "gateio",
		Raw:       gateioResp[0],
	}, nil
}

// OKXAdapter 实现 OKX 数据源
type OKXAdapter struct {
	symbolMap map[string]string
	baseURL   string
	logger    log.Logger
}

func NewOKXAdapter(logger log.Logger) *OKXAdapter {
	return &OKXAdapter{
		symbolMap: map[string]string{
			// 主流代币
			"BTC":  "BTC-USDT",
			"ETH":  "ETH-USDT",
			"BNB":  "BNB-USDT",
			"XRP":  "XRP-USDT",
			"XLM":  "XLM-USDT",
			"USDT": "USDT-USD",
			"USDC": "USDC-USDT",
			"DAI":  "DAI-USDT",

			// 迷因代币
			"DOGE":  "DOGE-USDT",
			"PEPE":  "PEPE-USDT",
			"CAT":   "CAT-USDT",
			"GOAT":  "GOAT-USDT",
			"DOGS":  "DOGS-USDT",
			"MEME":  "MEME-USDT",
			"CATS":  "CATS-USDT",
			"LADYS": "LADYS-USDT",
			"PUSS":  "PUSS-USDT",

			// 新兴代币
			"TON":      "TON-USDT",
			"SAVM":     "SAVM-USDT",
			"ORDI":     "ORDI-USDT",
			"NEIROCTO": "NEIROCTO-USDT",
			"MUBI":     "MUBI-USDT",
			"BB":       "BB-USDT",
			"SUNDOG":   "SUNDOG-USDT",
			"BAN":      "BAN-USDT",
			"MOODENG":  "MOODENG-USDT",
			"NEIRO":    "NEIRO-USDT",
			"1000SATS": "1000SATS-USDT",
			"NEIROETH": "NEIROETH-USDT",
			"AUCTION":  "AUCTION-USDT",
			"WIF":      "WIF-USDT",
			"BOME":     "BOME-USDT",
			"ACT":      "ACT-USDT",
		},
		baseURL: "https://www.okx.com",
		logger:  logger,
	}
}

func (o *OKXAdapter) Name() string {
	return "okx"
}

func (o *OKXAdapter) SupportsToken(symbol string) bool {
	_, ok := o.symbolMap[symbol]
	return ok
}

func (o *OKXAdapter) MapSymbol(symbol string) string {
	if mapped, ok := o.symbolMap[symbol]; ok {
		return mapped
	}
	return symbol + "-USDT"
}

func (o *OKXAdapter) BuildRequest(ctx context.Context, symbol string) (*http.Request, error) {
	mappedSymbol := o.MapSymbol(symbol)
	url := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", o.baseURL, mappedSymbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (o *OKXAdapter) ParseResponse(resp *http.Response, symbol string) (*PriceData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading OKX response: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OKX API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var okxResp struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstID    string `json:"instId"`
			Last      string `json:"last"`
			AskPx     string `json:"askPx"`
			BidPx     string `json:"bidPx"`
			Open24h   string `json:"open24h"`
			High24h   string `json:"high24h"`
			Low24h    string `json:"low24h"`
			VolCcy24h string `json:"volCcy24h"`
			Vol24h    string `json:"vol24h"`
			Ts        string `json:"ts"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &okxResp); err != nil {
		return nil, fmt.Errorf("error decoding OKX response: %w", err)
	}

	if okxResp.Code != "0" {
		return nil, fmt.Errorf("OKX API error: %s", okxResp.Msg)
	}

	if len(okxResp.Data) == 0 {
		return nil, fmt.Errorf("empty data from OKX")
	}

	price, err := math.LegacyNewDecFromStr(okxResp.Data[0].Last)
	if err != nil {
		return nil, fmt.Errorf("error parsing price from OKX: %w", err)
	}

	return &PriceData{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now().Unix(),
		Source:    "okx",
		Raw:       okxResp.Data[0],
	}, nil
}

// 创建价格服务
func NewPriceService(logger log.Logger, cmcApiKey string) *PriceServiceImpl {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建数据源适配器
	adapters := []DataSourceAdapter{
		NewBinanceAdapter(logger),
		NewCoinbaseAdapter(logger),
		NewCoinmarketcapAdapter(cmcApiKey, logger),
		NewGateioAdapter(logger),
		NewOKXAdapter(logger),
	}

	return &PriceServiceImpl{
		adapters: adapters,
		logger:   logger,
		client:   client,
	}
}

// 创建一个简单的测试函数，用于测试实际的API调用
func TestRealAPIs(symbol string, cmcApiKey string) (*TokenPrice, error) {
	logger := log.NewNopLogger()
	service := NewPriceService(logger, cmcApiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return service.GetTokenPrice(ctx, symbol)
}
