package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"intellix/dvs/price/server"

	"github.com/0xPellNetwork/pelldvs-libs/log"
)

func main() {
	// 解析命令行参数
	symbol := flag.String("symbol", "BTC", "Token symbol to fetch price for")
	cmcApiKey := flag.String("cmc-api-key", "", "CoinMarketCap API key")
	timeout := flag.Int("timeout", 30, "Timeout in seconds")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	allSymbols := flag.String("all", "", "Comma-separated list of symbols to fetch prices for")
	flag.Parse()

	// 设置日志
	var logger log.Logger
	if *verbose {
		logger = log.NewNopLogger() // 使用 NopLogger，实际项目中可以替换为适当的 logger
	} else {
		logger = log.NewNopLogger()
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	// 创建价格服务
	service := server.NewPriceService(logger, *cmcApiKey)

	// 获取价格
	if *allSymbols != "" {
		// 获取多个代币价格
		symbols := strings.Split(*allSymbols, ",")
		fmt.Printf("Fetching prices for %d tokens: %s\n", len(symbols), *allSymbols)

		results, err := service.GetTokenPrices(ctx, symbols)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// 打印结果
		fmt.Println("\n=== Results ===")
		for symbol, price := range results {
			fmt.Printf("Symbol: %s, Price: %s\n", symbol, price.Price.String())
			fmt.Printf("  Sources (%d):\n", len(price.SourcePrices))
			for source, sourcePrice := range price.SourcePrices {
				fmt.Printf("    %s: %s\n", source, sourcePrice.Price.String())
			}
			fmt.Println()
		}
	} else {
		// 获取单个代币价格
		fmt.Printf("Fetching price for %s...\n", *symbol)

		result, err := service.GetTokenPrice(ctx, *symbol)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// 打印结果
		fmt.Println("\n=== Result ===")
		fmt.Printf("Symbol: %s\n", result.Symbol)
		fmt.Printf("Aggregated Price: %s\n", result.Price.String())
		fmt.Printf("Timestamp: %s\n", time.Unix(result.Timestamp, 0).Format(time.RFC3339))
		fmt.Printf("Sources (%d):\n", len(result.SourcePrices))

		for source, price := range result.SourcePrices {
			fmt.Printf("  %s: %s\n", source, price.Price.String())
		}
	}
}
