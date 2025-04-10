package common

import (
	"testing"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/stretchr/testify/assert"

	"intellix/dvs/vrf/types"
)

func TestComputeVRF(t *testing.T) {
	// Generate a test key pair
	sk, _, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	assert.NoError(t, err)

	// Create test task metadata
	taskMetadata := &types.TaskMetadata{
		// Add necessary fields for testing
	}

	// Test ComputeVRF
	value, proof, err := ComputeVRF(sk, taskMetadata)
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.NotNil(t, proof)
	assert.NotEmpty(t, value)
	assert.NotEmpty(t, proof)
}

func TestVerifyVRF(t *testing.T) {
	// Generate a test key pair
	sk, pk, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	assert.NoError(t, err)

	// Create test task metadata
	taskMetadata := &types.TaskMetadata{
		// Add necessary fields for testing
	}

	// First compute VRF
	value, proof, err := ComputeVRF(sk, taskMetadata)
	assert.NoError(t, err)

	// Test VerifyVRF with correct values
	valid, err := VerifyVRF(pk, taskMetadata, value, proof)
	assert.NoError(t, err)
	assert.True(t, valid)

	// Test VerifyVRF with incorrect value
	invalidValue := make([]byte, len(value))
	copy(invalidValue, value)
	invalidValue[0] ^= 0xFF // Flip some bits
	valid, err = VerifyVRF(pk, taskMetadata, invalidValue, proof)
	assert.NoError(t, err)
	assert.False(t, valid)

	// Test VerifyVRF with incorrect proof
	invalidProof := make([]byte, len(proof))
	copy(invalidProof, proof)
	invalidProof[0] ^= 0xFF // Flip some bits
	valid, err = VerifyVRF(pk, taskMetadata, value, invalidProof)
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestVerifyVRF_ErrorCases(t *testing.T) {
	// Generate a test key pair
	_, pk, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	assert.NoError(t, err)

	// Test with nil task metadata
	valid, err := VerifyVRF(pk, nil, []byte{}, []byte{})
	assert.False(t, valid)
}
