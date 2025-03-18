package server

import (
	"context"
	"fmt"
	"intellix/dvs/vrf/types"
	"math/big"

	"cosmossdk.io/math"
)

func (s Server) GenerateRandomNumber(ctx context.Context, request *types.GenerateRandomNumberRequest) (*types.GenerateRandomNumberResponse, error) {
	s.logger.Info("GenerateRandomNumber", "in", fmt.Sprintf("%+v", request))

	randomNumber := big.NewInt(1)

	s.logger.Info("Generated random number", "result", randomNumber)
	return &types.GenerateRandomNumberResponse{
		TaskIndex:    request.Task.TaskIndex,
		RandomNumber: math.NewIntFromBigInt(randomNumber),
	}, nil
}
