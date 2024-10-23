package taskgateway

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type TaskGatewayCfg struct {
	EthEndpoint         string
	CosmosNetworkUrl    string
	ContractAddress     string
	ContractFromAddress string
}

func (t TaskGatewayCfg) Validate() error {
	if t.EthEndpoint == "" {
		return fmt.Errorf("eth endpoint cannot be empty")
	}
	if t.CosmosNetworkUrl == "" {
		return fmt.Errorf("cosmos network url cannot be empty")
	}
	if t.ContractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}
	if t.ContractFromAddress == "" {
		return fmt.Errorf("contract from address cannot be empty")
	}
	return nil
}

func convertAddressToString(addrStr string) (*common.Address, error) {
	mixedCaseAddress, err := common.NewMixedcaseAddressFromString(addrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert address: %w", err)
	}
	addr := mixedCaseAddress.Address()
	return &addr, nil
}
