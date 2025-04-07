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

type EasyMoneyFetchPriceService struct {
	logger log.Logger
}

func (s *EasyMoneyFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
	}

	var url = fmt.Sprintf("https://push2.eastmoney.com/api/qt/stock/get?invt=2&fltt=2&fields=f43,f71,f170,f169,f47,f48,f168,f50,f44,f45,f46,f60,f52,f49,f86,f161,f292&secid=%s&ut=fa5fd1943c7b386f172d6893dbfba10b&wbp2u=|0|0|0|web", base)
	s.logger.Info("start fetch price from EasyMoney", "base", base, "quote", quote, "url", url)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching price from EasyMoney", "error", err)
		return err
	}
	defer resp.Body.Close()

	var Response struct {
		RC     int    `json:"rc"`
		RT     int    `json:"rt"`
		SVR    int64  `json:"svr"`
		LT     int    `json:"lt"`
		Full   int    `json:"full"`
		Dlmkts string `json:"dlmkts"`
		Data   struct {
			F43  float64 `json:"f43"`
			F44  float64 `json:"f44"`
			F45  float64 `json:"f45"`
			F46  float64 `json:"f46"`
			F47  float64 `json:"f47"`
			F48  float64 `json:"f48"`
			F49  float64 `json:"f49"`
			F50  float64 `json:"f50"`
			F52  string  `json:"f52"`
			F60  float64 `json:"f60"`
			F71  string  `json:"f71"`
			F86  float64 `json:"f86"`
			F161 float64 `json:"f161"`
			F168 float64 `json:"f168"`
			F169 float64 `json:"f169"`
			F170 float64 `json:"f170"`
			F292 int     `json:"f292"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		s.logger.Error("Error decoding response", "error", err)
		return err
	}

	price := Response.Data.F43
	s.logger.Info("Fetched price from EasyMoney", "base", base, "quote", quote, "price", price)

	// Convert to Dec type
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)

	// Truncate decimal places if needed
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err, "price", priceStr)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceEasyMoney, Price: dec}

	return nil
}
