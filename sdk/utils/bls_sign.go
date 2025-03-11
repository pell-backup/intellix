package utils

import (
	"crypto/sha256"
	"fmt"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/0xPellNetwork/pelldvs/crypto/bls"
)

func SignWithBLS(blsKeyPair *bls.KeyPair, msgBytes []byte) []byte {
	if blsKeyPair == nil {
		return nil
	}

	hash := sha256.Sum256(msgBytes)
	var message [32]byte
	copy(message[:], hash[:32])
	signature := blsKeyPair.SignMessage(message)

	return signature.G1Point.Serialize()
}

func VerifyBLSSignature(pubKey *avsitypes.OperatorPubkeys, msgBytes []byte, signature []byte) error {
	if pubKey == nil {
		return fmt.Errorf("public key is nil")
	}

	if len(signature) == 0 {
		return fmt.Errorf("signature is empty")
	}

	// ensure signature is valid
	var message [32]byte
	copy(message[:], msgBytes[:32])

	g2PubKey := bls.NewZeroG2Point().Deserialize(pubKey.G2Pubkey)
	signatures := &bls.Signature{
		G1Point: bls.NewZeroG1Point().Deserialize(signature),
	}

	valid, err := signatures.Verify(g2PubKey, message)
	if err != nil {
		return fmt.Errorf("failed to verify BLS signature: %w", err)
	}

	if !valid {
		return fmt.Errorf("invalid BLS signature")
	}

	return nil
}
