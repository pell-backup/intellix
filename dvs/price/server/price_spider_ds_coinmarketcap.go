package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

// APIKeyManager is responsible for managing API keys
type CMCAPIKeyManager struct {
	logger log.Logger
	path   string
}

// NewAPIKeyManager creates a new API key manager
func NewCMCAPIKeyManager(logger log.Logger, path string) *CMCAPIKeyManager {
	return &CMCAPIKeyManager{
		logger: logger,
		path:   path,
	}
}

// GetAPIKey retrieves an API key from a file
func (m *CMCAPIKeyManager) GetAPIKey(name string) (string, error) {
	if m.path == "" {
		m.logger.Error("API key path not set")
		return "", fmt.Errorf("API key path not set")
	}

	keyPath := filepath.Join(m.path, fmt.Sprintf("%s.key", name))
	m.logger.Info("Reading API key", "path", keyPath)

	// Check if the file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		m.logger.Error("API key file does not exist", "path", keyPath)
		return "", fmt.Errorf("API key file does not exist: %s", keyPath)
	}

	// Read file contents
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		m.logger.Error("Failed to read API key file", "path", keyPath, "error", err)
		return "", fmt.Errorf("failed to read API key file: %w", err)
	}

	// Remove whitespace
	key := strings.TrimSpace(string(keyBytes))
	if key == "" {
		m.logger.Error("API key is empty", "path", keyPath)
		return "", fmt.Errorf("API key is empty")
	}

	// Remove potential quotes
	key = strings.Trim(key, "\"'")

	m.logger.Info("API key loaded successfully", "name", name)
	return key, nil
}

// CoinMarketCapFetchPriceService implements the CoinMarketCap data source
type CMCFetchPriceService struct {
	logger        log.Logger
	apiKeyManager *CMCAPIKeyManager
}

// NewCoinMarketCapFetchPriceService creates a new CoinMarketCap data source service
func NewCMCFetchPriceService(logger log.Logger, apiKeyPath string) *CMCFetchPriceService {
	return &CMCFetchPriceService{
		logger:        logger,
		apiKeyManager: NewCMCAPIKeyManager(logger, apiKeyPath),
	}
}

func (s *CMCFetchPriceService) fetchCoinPrice(base, quote string, tickConverter PriceTickConverter, priceChan chan<- *PriceInfo) error {
	s.logger.Info("CoinMarketCap data source activated", "status", "initializing")

	// Get API key
	apiKey, err := s.apiKeyManager.GetAPIKey("coinmarketcap")
	if err != nil {
		s.logger.Error("Failed to get CoinMarketCap API key", "error", err)
		return fmt.Errorf("failed to get CoinMarketCap API key: %w", err)
	}

	if apiKey == "" {
		s.logger.Error("Skipping CoinMarketCap data source due to missing API key")
		return fmt.Errorf("CoinMarketCap API key is empty")
	}

	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	if tick, ok := tickConverter[base]; ok {
		base = tick
		s.logger.Info("CoinMarketCap symbol conversion applied", "original_base", strings.ToUpper(base), "converted_base", tick)
	}
	if tick, ok := tickConverter[quote]; ok {
		quote = tick
		s.logger.Info("CoinMarketCap symbol conversion applied", "original_quote", strings.ToUpper(quote), "converted_quote", tick)
	}

	// CoinMarketCap API requires getting the cryptocurrency ID first
	symbolUrl := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s&convert=%s", base, quote)
	s.logger.Info("CoinMarketCap price request", "base", base, "quote", quote, "url", symbolUrl)

	req, err := http.NewRequest("GET", symbolUrl, nil)
	if err != nil {
		s.logger.Error("CoinMarketCap request creation failed", "error", err)
		return err
	}

	// Set API key
	req.Header.Set("X-CMC_PRO_API_KEY", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("CoinMarketCap price fetch failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("CoinMarketCap API returned non-OK status", "status", resp.StatusCode)
		return fmt.Errorf("CoinMarketCap API returned status %d", resp.StatusCode)
	}

	// Parse response
	var cmcResp struct {
		Data map[string]struct {
			Quote map[string]struct {
				Price float64 `json:"price"`
			} `json:"quote"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cmcResp); err != nil {
		s.logger.Error("CoinMarketCap response decode failed", "error", err)
		return err
	}

	// Check if there is data
	if len(cmcResp.Data) == 0 {
		s.logger.Error("CoinMarketCap empty response", "data_length", len(cmcResp.Data))
		return fmt.Errorf("empty response from CoinMarketCap")
	}

	// Get price
	coinData, ok := cmcResp.Data[base]
	if !ok {
		s.logger.Error("CoinMarketCap base symbol not found in response", "base", base)
		return fmt.Errorf("base symbol %s not found in CoinMarketCap response", base)
	}

	quoteData, ok := coinData.Quote[quote]
	if !ok {
		s.logger.Error("CoinMarketCap quote symbol not found in response", "quote", quote)
		return fmt.Errorf("quote symbol %s not found in CoinMarketCap response", quote)
	}

	price := quoteData.Price
	s.logger.Info("CoinMarketCap price fetch successful",
		"base", base,
		"quote", quote,
		"price", price,
		"status", "success")

	// Convert to Dec type
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)

	// Truncate decimal places if needed
	priceStr = TruncatePriceDecimal(priceStr, s.logger)

	dec, err := math.LegacyNewDecFromStr(priceStr)
	if err != nil {
		s.logger.Error("CoinMarketCap price conversion to Dec failed", "error", err, "price", priceStr)
		return err
	}
	priceChan <- &PriceInfo{DataSource: dataSourceCoinMarketCap, Price: dec}

	s.logger.Info("CoinMarketCap data source completed", "status", "success")
	return nil
}
