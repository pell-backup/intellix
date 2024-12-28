package pelldvs

import (
	"fmt"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"
	sdktypes "intellix/sdk/types"
)

func GetDvsRequestValidatedData(ctx sdktypes.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	reqData := ctx.DvsPostResponseData()

	validatedDataMsg, err := dvsservermanager.DecodeMsg(reqData)
	if err != nil {
		return nil, err
	}
	validatedData, ok := validatedDataMsg.(*dvstypes.RequestPostRequestValidatedData)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &dvstypes.RequestPostRequestValidatedData{}, validatedData)
	}
	return validatedData, nil
}
