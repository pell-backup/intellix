package logger

import (
	sdklog "cosmossdk.io/log"
	cmtlog "github.com/cometbft/cometbft/libs/log"
)

type CometBFTLogAdapter struct {
	sdkLogger sdklog.Logger
}

func NewCometBFTLogAdapter(sdkLogger sdklog.Logger) cmtlog.Logger {
	return &CometBFTLogAdapter{sdkLogger: sdkLogger}
}

func (a *CometBFTLogAdapter) Debug(msg string, keyVals ...interface{}) {
	a.sdkLogger.Debug(msg, keyVals...)
}

func (a *CometBFTLogAdapter) Info(msg string, keyVals ...interface{}) {
	a.sdkLogger.Info(msg, keyVals...)
}

func (a *CometBFTLogAdapter) Error(msg string, keyVals ...interface{}) {
	a.sdkLogger.Error(msg, keyVals...)
}

func (a *CometBFTLogAdapter) With(keyVals ...interface{}) cmtlog.Logger {
	return &CometBFTLogAdapter{sdkLogger: a.sdkLogger.With(keyVals...)}
}
