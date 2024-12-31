package pelldvs

import (
	"fmt"
	dvstypes "intellix/pkg/pelldvs/types"
	sdktypes "intellix/sdk/types"
)

func GetDvsRequestValidatedData(ctx sdktypes.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	validatedData := ctx.ValidateResponse()
	if validatedData == nil {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	return validatedData, nil
}
