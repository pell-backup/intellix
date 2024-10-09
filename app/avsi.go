package app

import (
	"context"

	avsi "github.com/0xPellNetwork/pelldvs/application"
)

func (app *App) ProcessRequest(ctx context.Context, req *avsi.RequestProcessRequest) (*avsi.ResponseProcessRequest, error) {
	return &avsi.ResponseProcessRequest{}, nil
}

func (app *App) PostRequest(ctx context.Context, req *avsi.RequestPostRequest) (*avsi.ResponsePostRequest, error) {
	return &avsi.ResponsePostRequest{}, nil
}
