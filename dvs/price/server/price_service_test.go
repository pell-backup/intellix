package server

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 模拟数据源适配器
type MockDataSourceAdapter struct {
	name            string
	supportedTokens map[string]bool
	priceResponse   *PriceData
	shouldFail      bool
}

func (m *MockDataSourceAdapter) Name() string {
	return m.name
}

func (m *MockDataSourceAdapter) SupportsToken(symbol string) bool {
	return m.supportedTokens[symbol]
}

func (m *MockDataSourceAdapter) BuildRequest(ctx context.Context, symbol string) (*http.Request, error) {
	return httptest.NewRequest("GET", "/mock-endpoint", nil), nil
}

func (m *MockDataSourceAdapter) ParseResponse(resp *http.Response, symbol string) (*PriceData, error) {
	if m.shouldFail {
		return nil, assert.AnError
	}
	return m.priceResponse, nil
}

// 模拟HTTP客户端
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// 测试价格服务
func TestPriceService_GetTokenPrice(t *testing.T) {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 创建测试用例
	testCases := []struct {
		name            string
		symbol          string
		adapters        []DataSourceAdapter
		expectedPrice   math.LegacyDec
		expectedError   bool
		expectedSources int
	}{
		{
			name:   "successful price fetch from multiple sources",
			symbol: "BTC",
			adapters: []DataSourceAdapter{
				&MockDataSourceAdapter{
					name:            "binance",
					supportedTokens: map[string]bool{"BTC": true},
					priceResponse: &PriceData{
						Symbol:    "BTC",
						Price:     math.LegacyNewDec(4350000).QuoInt64(100), // 43500.00
						Timestamp: time.Now().Unix(),
						Source:    "binance",
					},
				},
				&MockDataSourceAdapter{
					name:            "coinbase",
					supportedTokens: map[string]bool{"BTC": true},
					priceResponse: &PriceData{
						Symbol:    "BTC",
						Price:     math.LegacyNewDec(4360000).QuoInt64(100), // 43600.00
						Timestamp: time.Now().Unix(),
						Source:    "coinbase",
					},
				},
				&MockDataSourceAdapter{
					name:            "coinmarketcap",
					supportedTokens: map[string]bool{"BTC": true},
					priceResponse: &PriceData{
						Symbol:    "BTC",
						Price:     math.LegacyNewDec(4370000).QuoInt64(100), // 43700.00
						Timestamp: time.Now().Unix(),
						Source:    "coinmarketcap",
					},
				},
			},
			expectedPrice:   math.LegacyNewDec(4360000).QuoInt64(100), // 中位数: 43600.00
			expectedError:   false,
			expectedSources: 3,
		},
		{
			name:   "some sources fail",
			symbol: "ETH",
			adapters: []DataSourceAdapter{
				&MockDataSourceAdapter{
					name:            "binance",
					supportedTokens: map[string]bool{"ETH": true},
					priceResponse: &PriceData{
						Symbol:    "ETH",
						Price:     math.LegacyNewDec(2250000).QuoInt64(100), // 22500.00
						Timestamp: time.Now().Unix(),
						Source:    "binance",
					},
				},
				&MockDataSourceAdapter{
					name:            "coinbase",
					supportedTokens: map[string]bool{"ETH": true},
					shouldFail:      true,
				},
				&MockDataSourceAdapter{
					name:            "coinmarketcap",
					supportedTokens: map[string]bool{"ETH": true},
					priceResponse: &PriceData{
						Symbol:    "ETH",
						Price:     math.LegacyNewDec(2270000).QuoInt64(100), // 22700.00
						Timestamp: time.Now().Unix(),
						Source:    "coinmarketcap",
					},
				},
			},
			expectedPrice:   math.LegacyNewDec(2260000).QuoInt64(100), // 中位数: 22600.00 (平均值)
			expectedError:   false,
			expectedSources: 2,
		},
		{
			name:   "token not supported by any source",
			symbol: "UNKNOWN",
			adapters: []DataSourceAdapter{
				&MockDataSourceAdapter{
					name:            "binance",
					supportedTokens: map[string]bool{"BTC": true},
				},
				&MockDataSourceAdapter{
					name:            "coinbase",
					supportedTokens: map[string]bool{"BTC": true},
				},
			},
			expectedError: true,
		},
		{
			name:   "all sources fail",
			symbol: "DOGE",
			adapters: []DataSourceAdapter{
				&MockDataSourceAdapter{
					name:            "binance",
					supportedTokens: map[string]bool{"DOGE": true},
					shouldFail:      true,
				},
				&MockDataSourceAdapter{
					name:            "coinbase",
					supportedTokens: map[string]bool{"DOGE": true},
					shouldFail:      true,
				},
			},
			expectedError: true,
		},
		{
			name:   "even number of sources for median calculation",
			symbol: "TON",
			adapters: []DataSourceAdapter{
				&MockDataSourceAdapter{
					name:            "binance",
					supportedTokens: map[string]bool{"TON": true},
					priceResponse: &PriceData{
						Symbol:    "TON",
						Price:     math.LegacyNewDec(310000).QuoInt64(100), // 3100.00
						Timestamp: time.Now().Unix(),
						Source:    "binance",
					},
				},
				&MockDataSourceAdapter{
					name:            "coinbase",
					supportedTokens: map[string]bool{"TON": true},
					priceResponse: &PriceData{
						Symbol:    "TON",
						Price:     math.LegacyNewDec(320000).QuoInt64(100), // 3200.00
						Timestamp: time.Now().Unix(),
						Source:    "coinbase",
					},
				},
				&MockDataSourceAdapter{
					name:            "gateio",
					supportedTokens: map[string]bool{"TON": true},
					priceResponse: &PriceData{
						Symbol:    "TON",
						Price:     math.LegacyNewDec(330000).QuoInt64(100), // 3300.00
						Timestamp: time.Now().Unix(),
						Source:    "gateio",
					},
				},
				&MockDataSourceAdapter{
					name:            "okx",
					supportedTokens: map[string]bool{"TON": true},
					priceResponse: &PriceData{
						Symbol:    "TON",
						Price:     math.LegacyNewDec(340000).QuoInt64(100), // 3400.00
						Timestamp: time.Now().Unix(),
						Source:    "okx",
					},
				},
			},
			expectedPrice:   math.LegacyNewDec(325000).QuoInt64(100), // 中位数: 3250.00 (3200+3300)/2
			expectedError:   false,
			expectedSources: 4,
		},
	}

	// 运行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建模拟HTTP客户端
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// 返回一个空响应，实际内容由适配器的ParseResponse提供
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       http.NoBody,
					}, nil
				},
			}

			// 创建价格服务
			service := &PriceServiceImpl{
				adapters: tc.adapters,
				logger:   log.NewNopLogger(),
				client:   mockClient,
			}

			// 调用GetTokenPrice方法
			result, err := service.GetTokenPrice(context.Background(), tc.symbol)

			// 验证结果
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.symbol, result.Symbol)
				assert.Equal(t, tc.expectedPrice.String(), result.Price.String())
				assert.Equal(t, tc.expectedSources, len(result.SourcePrices))

				// 打印详细信息以便调试
				t.Logf("Symbol: %s, Aggregated Price: %s", result.Symbol, result.Price.String())
				for source, price := range result.SourcePrices {
					t.Logf("  Source: %s, Price: %s", source, price.Price.String())
				}
			}
		})
	}
}

// 测试价格聚合函数
func TestAggregatePrices(t *testing.T) {
	testCases := []struct {
		name          string
		sourcePrices  map[string]*PriceData
		expectedPrice math.LegacyDec
	}{
		{
			name: "odd number of prices",
			sourcePrices: map[string]*PriceData{
				"source1": {Price: math.LegacyNewDec(100)},
				"source2": {Price: math.LegacyNewDec(200)},
				"source3": {Price: math.LegacyNewDec(300)},
			},
			expectedPrice: math.LegacyNewDec(200), // 中位数
		},
		{
			name: "even number of prices",
			sourcePrices: map[string]*PriceData{
				"source1": {Price: math.LegacyNewDec(100)},
				"source2": {Price: math.LegacyNewDec(200)},
				"source3": {Price: math.LegacyNewDec(300)},
				"source4": {Price: math.LegacyNewDec(400)},
			},
			expectedPrice: math.LegacyNewDec(250), // (200+300)/2
		},
		{
			name: "single price",
			sourcePrices: map[string]*PriceData{
				"source1": {Price: math.LegacyNewDec(100)},
			},
			expectedPrice: math.LegacyNewDec(100),
		},
	}

	service := &PriceServiceImpl{logger: log.NewNopLogger()}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.aggregatePrices(tc.sourcePrices)
			assert.Equal(t, tc.expectedPrice.String(), result.String())
		})
	}
}

	// 排序价格
	for i := 0; i < len(prices); i++ {
		for j := i + 1; j < len(prices); j++ {
			if prices[i].GT(prices[j]) {
				prices[i], prices[j] = prices[j], prices[i]
			}
		}
	}

	// 取中位数
	if len(prices)%2 == 0 {
		// 偶数个元素，取中间两个的平均值
		middle1 := prices[len(prices)/2-1]
		middle2 := prices[len(prices)/2]
		return middle1.Add(middle2).QuoInt64(2)
	} else {
		// 奇数个元素，取中间值
		return prices[len(prices)/2]
	}
}
