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

	url := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/stock/get?invt=2&fltt=2&fields=f43,f71,f170,f169,f47,f48,f168,f50,f44,f45,f46,f60,f52,f49,f86,f161,f292&secid=%s&ut=fa5fd1943c7b386f172d6893dbfba10b&wbp2u=|0|0|0|web",
		base,
	)
	s.logger.Info("start fetch price from EasyMoney", "base", base, "quote", quote, "url", url)

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Error fetching price from EasyMoney", "error", err)
		return err
	}
	defer resp.Body.Close()

	var response struct {
		RC     int                    `json:"rc"`
		RT     int                    `json:"rt"`
		SVR    int64                  `json:"svr"`
		LT     int                    `json:"lt"`
		Full   int                    `json:"full"`
		Dlmkts string                 `json:"dlmkts"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		s.logger.Error("Error decoding response", "error", err)
		return err
	}

	// Parse f43 as price
	price, err := parseToFloat64(response.Data["f43"])
	if err != nil {
		s.logger.Error("Error parsing f43 as float64", "error", err, "value", response.Data["f43"])
		return err
	}

	s.logger.Info("Fetched price from EasyMoney", "base", base, "quote", quote, "price", price)

	// Convert to string then to Dec
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err, "price", priceStr)
		return err
	}

	priceChan <- &PriceInfo{DataSource: dataSourceEasyMoney, Price: dec}
	return nil
}

func parseToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case string:
		if val == "-" || val == "" {
			return 0, nil // or return an error if preferred
		}
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("unsupported type for float64: %T", v)
	}
}
