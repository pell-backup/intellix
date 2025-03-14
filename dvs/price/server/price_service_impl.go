package server

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
)

// GetTokenPrice 获取单个代币的价格
func (s *PriceServiceImpl) GetTokenPrice(ctx context.Context, symbol string) (*TokenPrice, error) {
	sourcePrices := make(map[string]*PriceData)

	// 从所有支持该代币的数据源获取价格
	for _, adapter := range s.adapters {
		if !adapter.SupportsToken(symbol) {
			continue
		}

		// 构建请求
		req, err := adapter.BuildRequest(ctx, symbol)
		if err != nil {
			s.logger.Error("Error building request", "source", adapter.Name(), "error", err)
			continue
		}

		// 执行请求
		resp, err := s.client.Do(req)
		if err != nil {
			s.logger.Error("Error executing request", "source", adapter.Name(), "error", err)
			continue
		}

		// 解析响应
		priceData, err := adapter.ParseResponse(resp, symbol)
		if err != nil {
			s.logger.Error("Error parsing response", "source", adapter.Name(), "error", err)
			continue
		}

		// 存储结果
		sourcePrices[adapter.Name()] = priceData
	}

	// 检查是否有足够的数据源返回结果
	if len(sourcePrices) == 0 {
		return nil, fmt.Errorf("no data sources returned valid prices for %s", symbol)
	}

	// 聚合价格
	aggregatedPrice := s.aggregatePrices(sourcePrices)

	return &TokenPrice{
		Symbol:       symbol,
		Price:        aggregatedPrice,
		Timestamp:    time.Now().Unix(),
		SourcePrices: sourcePrices,
	}, nil
}

// GetTokenPrices 获取多个代币的价格
func (s *PriceServiceImpl) GetTokenPrices(ctx context.Context, symbols []string) (map[string]*TokenPrice, error) {
	result := make(map[string]*TokenPrice)

	for _, symbol := range symbols {
		price, err := s.GetTokenPrice(ctx, symbol)
		if err != nil {
			s.logger.Error("Error getting price for token", "symbol", symbol, "error", err)
			continue
		}
		result[symbol] = price
	}

	return result, nil
}

// aggregatePrices 聚合多个数据源的价格
func (s *PriceServiceImpl) aggregatePrices(sourcePrices map[string]*PriceData) math.LegacyDec {
	// 提取所有价格
	var prices []math.LegacyDec
	for _, data := range sourcePrices {
		prices = append(prices, data.Price)
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
