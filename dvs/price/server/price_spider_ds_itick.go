package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

type ITickFetchPriceService struct {
	logger log.Logger
}

func (s *ITickFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	// Convert base and quote to uppercase
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	// Optional: If you still need to use a tickConverter
	// (comment out or remove if you no longer need it)
	if tick, ok := tickConverter[base]; ok {
		base = tick
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
	}

	url := fmt.Sprintf("stock/kline?region=US&code=%s&kType=1", base)
	// Construct the request body for Bouncebit
	// This example uses static values; modify as needed.
	// You might incorporate `base`, `quote` into the "url" field if required.
	reqBody := map[string]interface{}{
		"url":         url,
		"method":      "GET",
		"accept":      "application/json",
		"cacheSecond": 10,
		"body":        "",
		"headerName":  "Club_Quanto_Itick",
	}

	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		s.logger.Error("Error marshaling request body for Bouncebit", "error", err)
		return err
	}

	// Create the POST request
	bouncebitURL := "https://api-club.bouncebit.io/api/index/stock"
	s.logger.Info("Start fetching price from Bouncebit", "base", base, "quote", quote, "url", bouncebitURL)

	req, err := http.NewRequest(http.MethodPost, bouncebitURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.logger.Error("Error creating Bouncebit request", "error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second} // Adjust timeout to your needs
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Error fetching price from Bouncebit", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		s.logger.Error("Non-2xx Bouncebit response", "status", resp.StatusCode, "body", string(bodyBytes))
		return fmt.Errorf("bouncebit request failed with status: %d", resp.StatusCode)
	}

	// Define the structure to decode Bouncebit's response
	var bbResp struct {
		Code int `json:"code"`
		Data []struct {
			C  float64 `json:"c"`
			H  float64 `json:"h"`
			L  float64 `json:"l"`
			O  float64 `json:"o"`
			T  int64   `json:"t"`
			Tu float64 `json:"tu"`
			V  int64   `json:"v"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&bbResp); err != nil {
		s.logger.Error("Error decoding Bouncebit response", "error", err)
		return err
	}

	// Check the response code field
	if bbResp.Code != 0 {
		s.logger.Error("Bouncebit returned a non-zero code", "code", bbResp.Code)
		return fmt.Errorf("bouncebit response code != 0: %d", bbResp.Code)
	}

	// Make sure we have data
	if len(bbResp.Data) == 0 {
		s.logger.Error("Bouncebit response data is empty")
		return fmt.Errorf("bouncebit response data is empty")
	}

	// For demonstration, use the last element in Data
	latestEntry := bbResp.Data[len(bbResp.Data)-1]
	priceFloat := latestEntry.C

	s.logger.Info("Fetched price from Bouncebit",
		"base", base,
		"quote", quote,
		"price", priceFloat,
		"dataLen", len(bbResp.Data),
	)

	// Convert float64 to string
	priceStr := strconv.FormatFloat(priceFloat, 'f', -1, 64)

	// (Optional) Truncate decimal places to your desired precision
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	// Convert to a Cosmos Dec
	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("Error converting price to Dec", "error", err, "price", priceStr)
		return err
	}

	// Send the price to the channel
	priceChan <- &PriceInfo{DataSource: "bouncebit", Price: dec}

	return nil
}
