package utils

import (
	contractDataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

func AbiEncodeResponseTaskParam(taskIndex uint32, data []byte) ([]byte, error) {
	taskResponseType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "referenceTaskIndex",
			Type: "uint32",
		},
		{
			Name: "data",
			Type: "bytes",
		},
	})
	if err != nil {
		return nil, err
	}
	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}
	bytes, err := arguments.Pack(&contractDataOracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: taskIndex,
		Data:               data,
	})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func DigestKeccak256(data []byte) []byte {
	var taskResponseDigest [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

	return taskResponseDigest[:]
}
