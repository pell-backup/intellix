package server

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

// 使用更新的适配器系统获取价格
func fetchRawPricesWithAdapters(ctx context.Context, logger log.Logger, baseSymbol, quoteSymbol string, cmcApiKey string) (map[string]math.LegacyDec, error) {
	// 创建价格服务
	service := NewPriceService(logger, cmcApiKey)

	// 获取价格
	result, err := service.GetTokenPrice(ctx, baseSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price for %s: %w", baseSymbol, err)
	}

	// 转换为 map[string]math.LegacyDec 格式
	prices := make(map[string]math.LegacyDec)
	for source, priceData := range result.SourcePrices {
		prices[source] = priceData.Price
	}

	return prices, nil
}

// 替换原有的 fetchRawPrices 函数，使用新的适配器系统
func fetchRawPricesEnhanced(ctx context.Context, logger log.Logger, baseSymbol, quoteSymbol string, tickConverter PriceTickConverterByDataSource, cmcApiKey string) (map[string]math.LegacyDec, error) {
	// 首先尝试使用新的适配器系统获取价格
	prices, err := fetchRawPricesWithAdapters(ctx, logger, baseSymbol, quoteSymbol, cmcApiKey)
	if err == nil && len(prices) > 0 {
		logger.Info("Successfully fetched prices using new adapters", "baseSymbol", baseSymbol, "sources", len(prices))
		return prices, nil
	}

	// 如果新系统失败，回退到原有系统
	logger.Info("Falling back to legacy price fetching system", "baseSymbol", baseSymbol, "error", err)
	return fetchRawPrices(ctx, logger, baseSymbol, quoteSymbol, tickConverter)
}
