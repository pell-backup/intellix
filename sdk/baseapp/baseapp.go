package baseapp

import (
	"github.com/0xPellNetwork/pelldvs/libs/log"
)

type BaseApp struct {
	logger log.Logger
}

func NewBaseApp(
	logger log.Logger,
) *BaseApp {
	app := &BaseApp{
		logger: logger,
	}
	return app
}
