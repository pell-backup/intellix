package common

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/vrf"
	"intellix/dvs/vrf/types"
)

type VRFData struct {
	TaskMetadata *types.TaskMetadata
}

func ComputeVRF(sk keypair.PrivateKey, taskMetadata *types.TaskMetadata) ([]byte, []byte, error) {
	data, err := json.Marshal(&VRFData{
		TaskMetadata: taskMetadata,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("ComputeVRF failed to marshal vrfData: %s", err)
	}

	return vrf.Vrf(sk, data)
}

func VerifyVRF(pk keypair.PublicKey, task *types.TaskMetadata, value, proof []byte) (bool, error) {
	data, err := json.Marshal(&VRFData{
		TaskMetadata: task,
	})
	if err != nil {
		return false, fmt.Errorf("verifyVrf failed to marshal vrfData: %s", err)
	}

	return vrf.Verify(pk, data, value, proof)
}
