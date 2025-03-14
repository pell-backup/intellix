package server

import (
	"context"
	"net/http"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

// PriceData 表示从数据源获取的价格数据
type PriceData struct {
	Symbol    string         // 代币符号
	Price     math.LegacyDec // 价格
	Timestamp int64          // 时间戳
	Source    string         // 数据源
	Raw       interface{}    // 原始响应数据
}

// TokenPrice 表示聚合后的代币价格
type TokenPrice struct {
	Symbol       string                // 代币符号
	Price        math.LegacyDec        // 聚合后的价格
	Timestamp    int64                 // 时间戳
	SourcePrices map[string]*PriceData // 各数据源价格
}

// DataSourceAdapter 定义数据源适配器接口
type DataSourceAdapter interface {
	// 获取数据源名称
	Name() string

	// 检查数据源是否支持指定代币
	SupportsToken(symbol string) bool

	// 构建API请求
	BuildRequest(ctx context.Context, symbol string) (*http.Request, error)

	// 解析API响应
	ParseResponse(resp *http.Response, symbol string) (*PriceData, error)
}

// HTTPClient 定义HTTP客户端接口
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// PriceServiceImpl 实现价格服务
type PriceServiceImpl struct {
	adapters []DataSourceAdapter
	logger   log.Logger
	client   HTTPClient
}
