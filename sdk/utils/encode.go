package utils

import (
	"fmt"

	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
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
	bytes, err := arguments.Pack(&contractdataoracle.IDataOracleTaskResponse{
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

func AbiDecodeResponseTaskParam(data []byte) (*contractdataoracle.IDataOracleTaskResponse, error) {
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
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}

	// decode
	values, err := arguments.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack data: %w", err)
	}

	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected number of values: got %d, want 1", len(values))
	}

	r, ok := values[0].(struct {
		ReferenceTaskIndex uint32 `json:"referenceTaskIndex"`
		Data               []byte `json:"data"`
	})
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &contractdataoracle.IDataOracleTaskResponse{}, values[0])
	}

	return &contractdataoracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: r.ReferenceTaskIndex,
		Data:               r.Data,
	}, nil
}
