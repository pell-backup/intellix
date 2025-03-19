package dispatcher

import (
	"encoding/json"
	"fmt"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/vrf"
)

type VRFData struct {
	Task *contractdataoracle.ContractDataOracleNewTaskCreated
}

func computeVrf(sk keypair.PrivateKey, task *contractdataoracle.ContractDataOracleNewTaskCreated) ([]byte, []byte, error) {
	data, err := json.Marshal(&VRFData{
		Task: task,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("computeVrf failed to marshal vrfData: %s", err)
	}

	return vrf.Vrf(sk, data)
}
