package pelldvs

import (
	"fmt"
	dvstypes "intellix/sdk/pelldvs/types"
	sdktypes "intellix/sdk/types"
)

func GetDvsRequestValidatedData(ctx sdktypes.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	validatedData := ctx.ValidatedResponse()
	if validatedData == nil {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	return validatedData, nil
}
