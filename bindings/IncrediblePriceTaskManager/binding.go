// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIncrediblePriceTaskManager

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// IBLSSignatureCheckerNonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerNonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// IBLSSignatureCheckerQuorumStakeTotals is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerQuorumStakeTotals struct {
	SignedStakeForQuorum []*big.Int
	TotalStakeForQuorum  []*big.Int
}

// IIncrediblePriceTaskManagerTask is an auto generated low-level Go binding around an user-defined struct.
type IIncrediblePriceTaskManagerTask struct {
	PriceId                   [32]byte
	RequestData               []byte
	CallbackAddress           common.Address
	CallbackFunctionId        [4]byte
	TaskCreatedBlock          uint32
	QuorumNumbers             []byte
	QuorumThresholdPercentage uint32
}

// IIncrediblePriceTaskManagerTaskResponse is an auto generated low-level Go binding around an user-defined struct.
type IIncrediblePriceTaskManagerTaskResponse struct {
	ReferenceTaskIndex uint32
	Price              *big.Int
}

// IIncrediblePriceTaskManagerTaskResponseMetadata is an auto generated low-level Go binding around an user-defined struct.
type IIncrediblePriceTaskManagerTaskResponseMetadata struct {
	TaskResponsedBlock uint32
	HashOfNonSigners   [32]byte
}

// ContractIncrediblePriceTaskManagerMetaData contains all meta data concerning the ContractIncrediblePriceTaskManager contract.
var ContractIncrediblePriceTaskManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_taskResponseWindowBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"TASK_CHALLENGE_WINDOW_BLOCK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TASK_RESPONSE_WINDOW_BLOCK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"aggregator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allTaskHashes\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allTaskResponses\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkSignatures\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.QuorumStakeTotals\",\"components\":[{\"name\":\"signedStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"totalStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"generator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPriceId\",\"inputs\":[{\"name\":\"baseSymbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"quoteSymbol\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTaskResponseWindowBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_aggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_generator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"latestPriceInfos\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"price\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"publishTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestTaskNum\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"omniOperatorSharesManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIOmniOperatorSharesManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raiseAndResolveChallenge\",\"inputs\":[{\"name\":\"actualPrice\",\"type\":\"int192\",\"internalType\":\"int192\"},{\"name\":\"task\",\"type\":\"tuple\",\"internalType\":\"structIIncrediblePriceTaskManager.Task\",\"components\":[{\"name\":\"priceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"requestData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callbackAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackFunctionId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"taskCreatedBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"taskResponse\",\"type\":\"tuple\",\"internalType\":\"structIIncrediblePriceTaskManager.TaskResponse\",\"components\":[{\"name\":\"referenceTaskIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"price\",\"type\":\"int192\",\"internalType\":\"int192\"}]},{\"name\":\"taskResponseMetadata\",\"type\":\"tuple\",\"internalType\":\"structIIncrediblePriceTaskManager.TaskResponseMetadata\",\"components\":[{\"name\":\"taskResponsedBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashOfNonSigners\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"pubkeysOfNonSigningOperators\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestPrice\",\"inputs\":[{\"name\":\"baseSymbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"quoteSymbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"callbackAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackFunctionId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStaleStakesForbidden\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staleStakesForbidden\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"taskNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"taskSuccesfullyChallenged\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trySignatureAndApkVerification\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"apk\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"pairingSuccessful\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"siganatureIsValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateAggregator\",\"inputs\":[{\"name\":\"_aggregator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateGenerator\",\"inputs\":[{\"name\":\"_generator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePrice\",\"inputs\":[{\"name\":\"task\",\"type\":\"tuple\",\"internalType\":\"structIIncrediblePriceTaskManager.Task\",\"components\":[{\"name\":\"priceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"requestData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callbackAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackFunctionId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"taskCreatedBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"taskResponse\",\"type\":\"tuple\",\"internalType\":\"structIIncrediblePriceTaskManager.TaskResponse\",\"components\":[{\"name\":\"referenceTaskIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"price\",\"type\":\"int192\",\"internalType\":\"int192\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewTaskCreated\",\"inputs\":[{\"name\":\"taskIndex\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"task\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIIncrediblePriceTaskManager.Task\",\"components\":[{\"name\":\"priceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"requestData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callbackAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callbackFunctionId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"taskCreatedBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentage\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestPrice\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"priceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"baseSymbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"quoteSymbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"times\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleStakesForbiddenUpdate\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TaskChallengedSuccessfully\",\"inputs\":[{\"name\":\"taskIndex\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TaskChallengedUnsuccessfully\",\"inputs\":[{\"name\":\"taskIndex\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"challenger\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TaskCompleted\",\"inputs\":[{\"name\":\"taskIndex\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TaskResponded\",\"inputs\":[{\"name\":\"taskResponse\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIIncrediblePriceTaskManager.TaskResponse\",\"components\":[{\"name\":\"referenceTaskIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"price\",\"type\":\"int192\",\"internalType\":\"int192\"}]},{\"name\":\"taskResponseMetadata\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIIncrediblePriceTaskManager.TaskResponseMetadata\",\"components\":[{\"name\":\"taskResponsedBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"hashOfNonSigners\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateAggregator\",\"inputs\":[{\"name\":\"previousAggregator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAggregator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateGenerator\",\"inputs\":[{\"name\":\"previousGenerator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newGenerator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdatePrice\",\"inputs\":[{\"name\":\"priceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"price\",\"type\":\"int192\",\"indexed\":false,\"internalType\":\"int192\"},{\"name\":\"publishTime\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162004c8e38038062004c8e8339810160408190526200003591620001ea565b81806001600160a01b03166080816001600160a01b031681525050806001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200008f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000b5919062000231565b6001600160a01b031660a0816001600160a01b031681525050806001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200010d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000133919062000231565b6001600160a01b031660c0816001600160a01b03168152505060a0516001600160a01b03166366cef7106040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200018d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001b3919062000231565b6001600160a01b031660e0525063ffffffff16610100525062000258565b6001600160a01b0381168114620001e757600080fd5b50565b60008060408385031215620001fe57600080fd5b82516200020b81620001d1565b602084015190925063ffffffff811681146200022657600080fd5b809150509250929050565b6000602082840312156200024457600080fd5b81516200025181620001d1565b9392505050565b60805160a05160c05160e051610100516149a4620002ea600039600081816101fe015281816104c60152610ef301526000818161032501526127a90152600081816102fe0152818161182e015261298b01526000818161034c01528181612b610152612d230152600081816103730152818161099801528181610b83015281816124d4015261284601526149a46000f3fe608060405234801561001057600080fd5b506004361061018f5760003560e01c80636efb4636116100e45780639fe4ee47116100925780639fe4ee4714610458578063a5f3c6991461046b578063b98d09081461047e578063bdaf86581461048b578063c0c53b8b1461049e578063f2fde38b146104b1578063f5c9899d146104c4578063f63c5bab146104ea57600080fd5b80636efb463614610395578063715018a6146103b657806372d18e8d146103be5780637afa1eed146103cc5780638b00ce7c146103df5780638da5cb5b146103ef5780639a22de761461040057600080fd5b8063416c7e5e11610141578063416c7e5e146102a05780634d3141ce146102b35780635decc3f5146102c65780635df45946146102f957806366cef7101461032057806368304835146103475780636d14a9871461036e57600080fd5b806305880a65146101945780630ed21268146101ba578063171f1d5b146101cf5780631ad43189146101f9578063245a7bfc146102355780632cb223d5146102605780632d89f6fc14610280575b600080fd5b6101a76101a236600461398a565b6104f2565b6040519081526020015b60405180910390f35b6101cd6101c8366004613a5c565b6107bf565b005b6101e26101dd366004613bca565b610830565b6040805192151583529015156020830152016101b1565b6102207f000000000000000000000000000000000000000000000000000000000000000081565b60405163ffffffff90911681526020016101b1565b609c54610248906001600160a01b031681565b6040516001600160a01b0390911681526020016101b1565b6101a761026e366004613c1b565b609a6020526000908152604090205481565b6101a761028e366004613c1b565b60996020526000908152604090205481565b6101cd6102ae366004613c36565b610996565b6101a76102c1366004613c58565b610ad5565b6102e96102d4366004613c1b565b609b6020526000908152604090205460ff1681565b60405190151581526020016101b1565b6102487f000000000000000000000000000000000000000000000000000000000000000081565b6102487f000000000000000000000000000000000000000000000000000000000000000081565b6102487f000000000000000000000000000000000000000000000000000000000000000081565b6102487f000000000000000000000000000000000000000000000000000000000000000081565b6103a86103a3366004613f65565b610b0e565b6040516101b1929190614024565b6101cd610c7d565b60985463ffffffff16610220565b609d54610248906001600160a01b031681565b6098546102209063ffffffff1681565b6033546001600160a01b0316610248565b61043661040e36600461406d565b609760205260009081526040902054601781900b90600160c01b90046001600160401b031682565b6040805160179390930b83526001600160401b039091166020830152016101b1565b6101cd610466366004613a5c565b610c91565b6101cd6104793660046140aa565b610d02565b6065546102e99060ff1681565b6101cd610499366004614130565b611398565b6101cd6104ac3660046141ba565b611966565b6101cd6104bf366004613a5c565b611aab565b7f0000000000000000000000000000000000000000000000000000000000000000610220565b610220606481565b60006104fc6137f3565b61050581611b21565b50610566604051806040016040528060048152602001636261736560e01b8152508c8c8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869493925050611b3c9050565b6105c76040518060400160405280600581526020016471756f746560d81b8152508a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869493925050611b3c9050565b60408051808201909152600581526474696d657360d81b60208201526305f5e100906105f590839083611b55565b60006106038d8d8d8d610ad5565b905080336001600160a01b03167fcffe5e47d40a1edafcd79aaafab384bbdcdf0c48c2b6c32fe79dcdb89ef3023d8f8f8f8f8860405161064795949392919061422e565b60405180910390a36040805160e08101825260008082526060602083018190529282018190528282018190526080820181905260a082019290925260c081019190915281815283516020808301919091526001600160a01b038b166040808401919091526001600160e01b03198b16606084015263ffffffff43811660808501528a1660c08401528051601f890183900483028101830190915287815290889088908190840183828082843760009201919091525050505060a08201526040516107159082906020016142b8565b60408051601f1981840301815282825280516020918201206098805463ffffffff9081166000908152609990945293909220555416907f36e52b7ae44ca8f480b92f4c6107dcd429968909437e0bb327dca642b85d0a30906107789084906142b8565b60405180910390a26098546107949063ffffffff166001614369565b6098805463ffffffff191663ffffffff92909216919091179055509c9b505050505050505050505050565b6107c7611b69565b609d54604080516001600160a01b03928316815291831660208301527f774efe45a4bbe359a5d19d71e74cbf8a7215d5b4ae22694324fa7f5d136cd678910160405180910390a1609d80546001600160a01b0319166001600160a01b0392909216919091179055565b60008060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001878760000151886020015188600001516000600281106108785761087861438d565b60200201518951600160200201518a6020015160006002811061089d5761089d61438d565b60200201518b602001516001600281106108b9576108b961438d565b602090810291909101518c518d8301516040516109169a99989796959401988952602089019790975260408801959095526060870193909352608086019190915260a085015260c084015260e08301526101008201526101200190565b6040516020818303038152906040528051906020012060001c61093991906143a3565b905061098861095261094b8884611bc3565b8690611c48565b61095a611cd1565b61097e61096f85610969611d91565b90611bc3565b6109788c611db2565b90611c48565b886201d4c0611e35565b909890975095505050505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109f4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a1891906143c5565b6001600160a01b0316336001600160a01b031614610ac95760405162461bcd60e51b815260206004820152605c60248201527f424c535369676e6174757265436865636b65722e6f6e6c79436f6f7264696e6160448201527f746f724f776e65723a2063616c6c6572206973206e6f7420746865206f776e6560648201527f72206f6620746865207265676973747279436f6f7264696e61746f7200000000608482015260a4015b60405180910390fd5b610ad281612059565b50565b600084848484604051602001610aee94939291906143e2565b604051602081830303815290604052805190602001209050949350505050565b60408051808201909152606080825260208201526000610b30868686866120a0565b6000610bf987878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060408051639aa1653d60e01b815290516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169350639aa1653d925060048083019260209291908290030181865afa158015610bd0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bf49190614404565b612288565b9050600080610c09878785612309565b91509150610c168261265d565b9150600080610c298b8b8b8b87896126ec565b91509150610c388c828a612eb8565b6000898460200151604051602001610c51929190614427565b60408051808303601f190181529190528051602090910120929d929c50919a5050505050505050505050565b610c85611b69565b610c8f6000612f9d565b565b610c99611b69565b609c54604080516001600160a01b03928316815291831660208301527f35a4ee0f60435b1bc206614dffe6bfd98db7e7450ab2ff2da4611d15ce8ffd69910160405180910390a1609c80546001600160a01b0319166001600160a01b0392909216919091179055565b609c546001600160a01b03163314610d5c5760405162461bcd60e51b815260206004820152601d60248201527f41676772656761746f72206d757374206265207468652063616c6c65720000006044820152606401610ac0565b6000610d6e60a0850160808601613c1b565b9050366000610d8060a087018761446f565b90925090506000610d9760e0880160c08901613c1b565b905060996000610daa6020890189613c1b565b63ffffffff1663ffffffff1681526020019081526020016000205487604051602001610dd691906144fa565b6040516020818303038152906040528051906020012014610e5f5760405162461bcd60e51b815260206004820152603d60248201527f737570706c696564207461736b20646f6573206e6f74206d617463682074686560448201527f206f6e65207265636f7264656420696e2074686520636f6e74726163740000006064820152608401610ac0565b6000609a81610e7160208a018a613c1b565b63ffffffff1663ffffffff1681526020019081526020016000205414610eee5760405162461bcd60e51b815260206004820152602c60248201527f41676772656761746f722068617320616c726561647920726573706f6e64656460448201526b20746f20746865207461736b60a01b6064820152608401610ac0565b610f187f000000000000000000000000000000000000000000000000000000000000000085614369565b63ffffffff164363ffffffff161115610f895760405162461bcd60e51b815260206004820152602d60248201527f41676772656761746f722068617320726573706f6e64656420746f207468652060448201526c7461736b20746f6f206c61746560981b6064820152608401610ac0565b600086604051602001610f9c91906145ef565b604051602081830303815290604052805190602001209050600080610fc48387878a8c610b0e565b9150915060005b858110156110c3578460ff1683602001518281518110610fed57610fed61438d565b6020026020010151610fff91906145fd565b6001600160601b03166064846000015183815181106110205761102061438d565b60200260200101516001600160601b031661103b9190614620565b10156110b1576040805162461bcd60e51b81526020600482015260248101919091527f5369676e61746f7269657320646f206e6f74206f776e206174206c656173742060448201527f7468726573686f6c642070657263656e74616765206f6620612071756f72756d6064820152608401610ac0565b806110bb81614637565b915050610fcb565b5060408051808201825263ffffffff431681526020808201849052915190916110f0918c91849101614650565b60405160208183030381529060405280519060200120609a60008c600001602081019061111d9190613c1b565b63ffffffff1663ffffffff1681526020019081526020016000208190555060405180604001604052808b6020016020810190611159919061467c565b60170b815263ffffffff42166020918201528c3560009081526097825260409020825192909101516001600160401b0316600160c01b026001600160c01b0390921691909117905562061a805a10156111f45760405162461bcd60e51b815260206004820181905260248201527f4d7573742070726f7669646520636f6e73756d657220656e6f756768206761736044820152606401610ac0565b600061120660608d0160408e01613a5c565b6001600160a01b031661121f60808e0160608f01614697565b8d600001358d6020016020810190611237919061467c565b604080516001600160e01b031990941660208501526024840192909252901b6044820152605c0160408051601f1981840301815290829052611278916146b2565b6000604051808303816000865af19150503d80600081146112b5576040519150601f19603f3d011682016040523d82523d6000602084013e6112ba565b606091505b50509050806112fd5760405162461bcd60e51b815260206004820152600f60248201526e18d85b1b189858dac819985a5b1959608a1b6044820152606401610ac0565b7f68b0dddd02e8a85939b35ac692e68162f8a1dbdd71737d8b20620abbedff5c3d8b8360405161132e929190614650565b60405180910390a18b357f3502e3be13f3f611c5bd255b5ee64196626f6ebaac138cf0f4db61166976e55361136960408e0160208f0161467c565b6040805160179290920b825263ffffffff421660208301520160405180910390a2505050505050505050505050565b60006113a76020850185613c1b565b63ffffffff81166000908152609a60205260409020549091506114165760405162461bcd60e51b815260206004820152602160248201527f5461736b206861736e2774206265656e20726573706f6e64656420746f2079656044820152601d60fa1b6064820152608401610ac0565b83836040516020016114299291906146ce565b60408051601f19818403018152918152815160209283012063ffffffff84166000908152609a909352912054146114c85760405162461bcd60e51b815260206004820152603d60248201527f5461736b20726573706f6e736520646f6573206e6f74206d617463682074686560448201527f206f6e65207265636f7264656420696e2074686520636f6e74726163740000006064820152608401610ac0565b63ffffffff81166000908152609b602052604090205460ff16156115605760405162461bcd60e51b815260206004820152604360248201527f54686520726573706f6e736520746f2074686973207461736b2068617320616c60448201527f7265616479206265656e206368616c6c656e676564207375636365737366756c606482015262363c9760e91b608482015260a401610ac0565b606461156f6020850185613c1b565b6115799190614369565b63ffffffff164363ffffffff1611156115f45760405162461bcd60e51b815260206004820152603760248201527f546865206368616c6c656e676520706572696f6420666f72207468697320746160448201527639b5903430b99030b63932b0b23c9032bc3834b932b21760491b6064820152608401610ac0565b6000611606604086016020870161467c565b601788810b91900b149050600181900361165457604051339063ffffffff8416907ffd3e26beeb5967fc5a57a0446914eabc45b4aa474c67a51b4b5160cac60ddb0590600090a3505061195f565b600083516001600160401b0381111561166f5761166f613a79565b604051908082528060200260200182016040528015611698578160200160208202803683370190505b50905060005b845181101561170a576116db8582815181106116bc576116bc61438d565b6020026020010151805160009081526020918201519091526040902090565b8282815181106116ed576116ed61438d565b60209081029190910101528061170281614637565b91505061169e565b50600061171d60a0890160808a01613c1b565b8260405160200161172f929190614427565b604051602081830303815290604052805190602001209050856020013581146117d95760405162461bcd60e51b815260206004820152605060248201527f546865207075626b657973206f66206e6f6e2d7369676e696e67206f7065726160448201527f746f727320737570706c69656420627920746865206368616c6c656e6765722060648201526f30b932903737ba1031b7b93932b1ba1760811b608482015260a401610ac0565b600085516001600160401b038111156117f4576117f4613a79565b60405190808252806020026020018201604052801561181d578160200160208202803683370190505b50905060005b8651811015611910577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e8bb9ae685838151811061186d5761186d61438d565b60200260200101516040518263ffffffff1660e01b815260040161189391815260200190565b602060405180830381865afa1580156118b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118d491906143c5565b8282815181106118e6576118e661438d565b6001600160a01b03909216602092830291909101909101528061190881614637565b915050611823565b5063ffffffff85166000818152609b6020526040808220805460ff19166001179055513392917fc20d1bb0f1623680306b83d4ff4bb99a2beb9d86d97832f3ca40fd13a29df1ec91a350505050505b5050505050565b600054610100900460ff16158080156119865750600054600160ff909116105b806119a05750303b1580156119a0575060005460ff166001145b611a035760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610ac0565b6000805460ff191660011790558015611a26576000805461ff0019166101001790555b611a2f84612f9d565b609c80546001600160a01b038086166001600160a01b031992831617909255609d8054928516929091169190911790558015611aa5576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050565b611ab3611b69565b6001600160a01b038116611b185760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610ac0565b610ad281612f9d565b611b296137f3565b611b3582610100612fef565b5090919050565b611b468383613047565b611b508382613047565b505050565b611b5f8383613047565b611b50838261305e565b6033546001600160a01b03163314610c8f5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610ac0565b611bcb61380d565b611bd3613827565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080611c0257fe5b5080611c405760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b6044820152606401610ac0565b505092915050565b611c5061380d565b611c58613845565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080611c9357fe5b5080611c405760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b6044820152606401610ac0565b611cd9613863565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b611d9961380d565b5060408051808201909152600181526002602082015290565b611dba61380d565b60008080611dd660008051602061490f833981519152866143a3565b90505b611de2816130c1565b909350915060008051602061490f8339815191528283098303611e1b576040805180820190915290815260208101919091529392505050565b60008051602061490f833981519152600182089050611dd9565b604080518082018252868152602080820186905282518084019093528683528201849052600091829190611e67613888565b60005b600281101561202c576000611e80826006614620565b9050848260028110611e9457611e9461438d565b60200201515183611ea6836000614701565b600c8110611eb657611eb661438d565b6020020152848260028110611ecd57611ecd61438d565b60200201516020015183826001611ee49190614701565b600c8110611ef457611ef461438d565b6020020152838260028110611f0b57611f0b61438d565b6020020151515183611f1e836002614701565b600c8110611f2e57611f2e61438d565b6020020152838260028110611f4557611f4561438d565b6020020151516001602002015183611f5e836003614701565b600c8110611f6e57611f6e61438d565b6020020152838260028110611f8557611f8561438d565b602002015160200151600060028110611fa057611fa061438d565b602002015183611fb1836004614701565b600c8110611fc157611fc161438d565b6020020152838260028110611fd857611fd861438d565b602002015160200151600160028110611ff357611ff361438d565b602002015183612004836005614701565b600c81106120145761201461438d565b6020020152508061202481614637565b915050611e6a565b506120356138a7565b60006020826101808560088cfa9151919c9115159b50909950505050505050505050565b6065805460ff19168215159081179091556040519081527f40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc9060200160405180910390a150565b60008390036120ff5760405162461bcd60e51b8152602060048201526037602482015260008051602061494f8339815191526044820152761c995cce88195b5c1d1e481c5d5bdc9d5b481a5b9c1d5d604a1b6064820152608401610ac0565b60408101515183148015612117575060a08101515183145b8015612127575060c08101515183145b8015612137575060e08101515183145b6121a15760405162461bcd60e51b8152602060048201526041602482015260008051602061494f83398151915260448201527f7265733a20696e7075742071756f72756d206c656e677468206d69736d6174636064820152600d60fb1b608482015260a401610ac0565b805151602082015151146122195760405162461bcd60e51b81526020600482015260446024820181905260008051602061494f833981519152908201527f7265733a20696e707574206e6f6e7369676e6572206c656e677468206d69736d6064820152630c2e8c6d60e31b608482015260a401610ac0565b4363ffffffff168263ffffffff1610611aa55760405162461bcd60e51b815260206004820152603c602482015260008051602061494f83398151915260448201527f7265733a20696e76616c6964207265666572656e636520626c6f636b000000006064820152608401610ac0565b60008061229484613143565b9050808360ff166001901b116123005760405162461bcd60e51b815260206004820152603f602482015260008051602061492f83398151915260448201527f69746d61703a206269746d61702065786365656473206d61782076616c7565006064820152608401610ac0565b90505b92915050565b61231161380d565b6040805180820190915260608082526020820152604051806040016040528060008152602001600081525091508360200151516001600160401b0381111561235b5761235b613a79565b604051908082528060200260200182016040528015612384578160200160208202803683370190505b5081526020840151516001600160401b038111156123a4576123a4613a79565b6040519080825280602002602001820160405280156123cd578160200160208202803683370190505b50602082015260005b846020015151811015612654576123fc856020015182815181106116bc576116bc61438d565b826020015182815181106124125761241261438d565b602090810291909101015280156124d2576020820151612433600183614714565b815181106124435761244361438d565b602002602001015160001c826020015182815181106124645761246461438d565b602002602001015160001c116124d2576040805162461bcd60e51b815260206004820152602481019190915260008051602061494f83398151915260448201527f7265733a206e6f6e5369676e65725075626b657973206e6f7420736f727465646064820152608401610ac0565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166304ec6351836020015183815181106125175761251761438d565b602002602001015188886000015185815181106125365761253661438d565b60200260200101516040518463ffffffff1660e01b81526004016125739392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015612590573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125b49190614727565b6001600160c01b0316826000015182815181106125d3576125d361438d565b60200260200101818152505061264061263961260d86856000015185815181106125ff576125ff61438d565b6020026020010151166132af565b876020015184815181106126235761262361438d565b60200260200101516132da90919063ffffffff16565b8490611c48565b92508061264c81614637565b9150506123d6565b50935093915050565b61266561380d565b815115801561267657506020820151155b15612694575050604080518082019091526000808252602082015290565b60405180604001604052808360000151815260200160008051602061490f83398151915284602001516126c791906143a3565b6126df9060008051602061490f833981519152614714565b905292915050565b919050565b604080518082019091526060808252602082015261270861380d565b866001600160401b0381111561272057612720613a79565b604051908082528060200260200182016040528015612749578160200160208202803683370190505b506020830152866001600160401b0381111561276757612767613a79565b604051908082528060200260200182016040528015612790578160200160208202803683370190505b50825260655460ff166000816127a7576000612829565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c448feb86040518163ffffffff1660e01b8152600401602060405180830381865afa158015612805573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128299190614750565b905060005b89811015612ea7578215612989578863ffffffff16827f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663249a0c428e8e868181106128855761288561438d565b60405160e085901b6001600160e01b031916815292013560f81c600483015250602401602060405180830381865afa1580156128c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128e99190614750565b6128f39190614701565b116129895760405162461bcd60e51b8152602060048201526066602482015260008051602061494f83398151915260448201527f7265733a205374616b6552656769737472792075706461746573206d7573742060648201527f62652077697468696e207769746864726177616c44656c6179426c6f636b732060848201526577696e646f7760d01b60a482015260c401610ac0565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166368bccaac8c8c848181106129ca576129ca61438d565b9050013560f81c60f81b60f81c8b8b60a0015185815181106129ee576129ee61438d565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015612a4a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a6e9190614769565b6001600160401b031916612a91896040015183815181106116bc576116bc61438d565b67ffffffffffffffff191614612b2d5760405162461bcd60e51b8152602060048201526061602482015260008051602061494f83398151915260448201527f7265733a2071756f72756d41706b206861736820696e2073746f72616765206460648201527f6f6573206e6f74206d617463682070726f76696465642071756f72756d2061706084820152606b60f81b60a482015260c401610ac0565b612b5d88604001518281518110612b4657612b4661438d565b602002602001015187611c4890919063ffffffff16565b95507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c8294c568c8c84818110612ba057612ba061438d565b9050013560f81c60f81b60f81c8b8b60c001518581518110612bc457612bc461438d565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015612c20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c449190614794565b85602001518281518110612c5a57612c5a61438d565b6001600160601b03909216602092830291909101820152850151805182908110612c8657612c8661438d565b602002602001015185600001518281518110612ca457612ca461438d565b60200260200101906001600160601b031690816001600160601b0316815250506000805b896020015151811015612e9257612d1c89600001518281518110612cee57612cee61438d565b60200260200101518e8e86818110612d0857612d0861438d565b600192013560f81c9290921c811614919050565b15612e80577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663f2be94ae8e8e86818110612d6257612d6261438d565b9050013560f81c60f81b60f81c8d8c602001518581518110612d8657612d8661438d565b60200260200101518e60e001518881518110612da457612da461438d565b60200260200101518781518110612dbd57612dbd61438d565b60209081029190910101516040516001600160e01b031960e087901b16815260ff909416600485015263ffffffff92831660248501526044840191909152166064820152608401602060405180830381865afa158015612e21573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e459190614794565b8751805185908110612e5957612e5961438d565b60200260200101818151612e6d91906147bd565b6001600160601b03169052506001909101905b80612e8a81614637565b915050612cc8565b50508080612e9f90614637565b91505061282e565b508492505050965096945050505050565b600080612ecf858585606001518660800151610830565b9150915081612f405760405162461bcd60e51b8152602060048201526043602482015260008051602061494f83398151915260448201527f7265733a2070616972696e6720707265636f6d70696c652063616c6c206661696064820152621b195960ea1b608482015260a401610ac0565b8061195f5760405162461bcd60e51b8152602060048201526039602482015260008051602061494f8339815191526044820152781c995cce881cda59db985d1d5c99481a5cc81a5b9d985b1a59603a1b6064820152608401610ac0565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b612ff76137f3565b6130026020836143a3565b1561302a576130126020836143a3565b61301d906020614714565b6130279083614701565b91505b506020828101829052604080518085526000815290920101905290565b61305482600383516133b1565b611b5082826134b8565b67ffffffffffffffff1981121561307d5761307982826134d2565b5050565b6001600160401b03811315613096576130798282613514565b600081126130aa57613079826000836133b1565b6130798260016130bc846000196147dd565b6133b1565b6000808060008051602061490f833981519152600360008051602061490f8339815191528660008051602061490f833981519152888909090890506000613137827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f5260008051602061490f833981519152613537565b91959194509092505050565b6000610100825111156131ba5760405162461bcd60e51b81526020600482015260446024820181905260008051602061492f833981519152908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610ac0565b81516000036131cb57506000919050565b600080836000815181106131e1576131e161438d565b0160200151600160f89190911c81901b92505b84518110156132a65784818151811061320f5761320f61438d565b0160200151600160f89190911c1b91508282116132925760405162461bcd60e51b8152602060048201526047602482015260008051602061492f83398151915260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610ac0565b9181179161329f81614637565b90506131f4565b50909392505050565b6000805b8215612303576132c4600184614714565b90921691806132d2816147fd565b9150506132b3565b6132e261380d565b6102008261ffff161061332a5760405162461bcd60e51b815260206004820152601060248201526f7363616c61722d746f6f2d6c6172676560801b6044820152606401610ac0565b8161ffff1660010361333d575081612303565b6040805180820190915260008082526020820181905284906001905b8161ffff168661ffff16106133a657600161ffff871660ff83161c81169003613389576133868484611c48565b93505b6133938384611c48565b92506201fffe600192831b169101613359565b509195945050505050565b6017816001600160401b0316116133d557611aa58360e0600585901b1683176135e0565b60ff816001600160401b031611613411576133fb836018611fe0600586901b16176135e0565b50611aa5836001600160401b03831660016135f8565b61ffff816001600160401b03161161344e57613438836019611fe0600586901b16176135e0565b50611aa5836001600160401b03831660026135f8565b63ffffffff816001600160401b03161161348d5761347783601a611fe0600586901b16176135e0565b50611aa5836001600160401b03831660046135f8565b6134a283601b611fe0600586901b16176135e0565b50611aa5836001600160401b03831660086135f8565b6134c06137f3565b61230083846000015151848551613619565b6134dd8260c36135e0565b50613079826134ee836000196147dd565b60405160200161350091815260200190565b6040516020818303038152906040526136f6565b61351f8260c26135e0565b50613079828260405160200161350091815260200190565b6000806135426138a7565b61354a6138c5565b602080825281810181905260408201819052606082018890526080820187905260a082018690528260c08360056107d05a03fa9250828061358757fe5b50826135d55760405162461bcd60e51b815260206004820152601a60248201527f424e3235342e6578704d6f643a2063616c6c206661696c7572650000000000006044820152606401610ac0565b505195945050505050565b6135e86137f3565b6123008384600001515184613703565b6136006137f3565b613611848560000151518585613751565b949350505050565b6136216137f3565b825182111561362f57600080fd5b602085015161363e8386614701565b111561367157613671856136618760200151878661365c9190614701565b6137c5565b61366c906002614620565b6137dc565b6000808651805187602083010193508088870111156136905787860182525b505050602084015b602084106136d057805182526136af602083614701565b91506136bc602082614701565b90506136c9602085614714565b9350613698565b51815160001960208690036101000a019081169019919091161790525083949350505050565b61305482600283516133b1565b61370b6137f3565b8360200151831061372b5761372b848560200151600261366c9190614620565b8351805160208583010184815350808503613747576001810182525b5093949350505050565b6137596137f3565b60208501516137688584614701565b111561377c5761377c856136618685614701565b6000600161378c84610100614902565b6137969190614714565b90508551838682010185831982511617815250805184870111156137ba5783860181525b509495945050505050565b6000818311156137d6575081612303565b50919050565b81516137e88383612fef565b50611aa583826134b8565b604051806040016040528060608152602001600081525090565b604051806040016040528060008152602001600081525090565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b60405180604001604052806138766138e3565b81526020016138836138e3565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60008083601f84011261391357600080fd5b5081356001600160401b0381111561392a57600080fd5b60208301915083602082850101111561394257600080fd5b9250929050565b6001600160a01b0381168114610ad257600080fd5b80356001600160e01b0319811681146126e757600080fd5b803563ffffffff811681146126e757600080fd5b600080600080600080600080600060c08a8c0312156139a857600080fd5b89356001600160401b03808211156139bf57600080fd5b6139cb8d838e01613901565b909b50995060208c01359150808211156139e457600080fd5b6139f08d838e01613901565b909950975060408c01359150613a0582613949565b819650613a1460608d0161395e565b9550613a2260808d01613976565b945060a08c0135915080821115613a3857600080fd5b50613a458c828d01613901565b915080935050809150509295985092959850929598565b600060208284031215613a6e57600080fd5b813561230081613949565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715613ab157613ab1613a79565b60405290565b60405161010081016001600160401b0381118282101715613ab157613ab1613a79565b604051601f8201601f191681016001600160401b0381118282101715613b0257613b02613a79565b604052919050565b600060408284031215613b1c57600080fd5b613b24613a8f565b9050813581526020820135602082015292915050565b600082601f830112613b4b57600080fd5b613b53613a8f565b806040840185811115613b6557600080fd5b845b81811015613b7f578035845260209384019301613b67565b509095945050505050565b600060808284031215613b9c57600080fd5b613ba4613a8f565b9050613bb08383613b3a565b8152613bbf8360408401613b3a565b602082015292915050565b6000806000806101208587031215613be157600080fd5b84359350613bf28660208701613b0a565b9250613c018660608701613b8a565b9150613c108660e08701613b0a565b905092959194509250565b600060208284031215613c2d57600080fd5b61230082613976565b600060208284031215613c4857600080fd5b8135801515811461230057600080fd5b60008060008060408587031215613c6e57600080fd5b84356001600160401b0380821115613c8557600080fd5b613c9188838901613901565b90965094506020870135915080821115613caa57600080fd5b50613cb787828801613901565b95989497509550505050565b60006001600160401b03821115613cdc57613cdc613a79565b5060051b60200190565b600082601f830112613cf757600080fd5b81356020613d0c613d0783613cc3565b613ada565b82815260059290921b84018101918181019086841115613d2b57600080fd5b8286015b84811015613d4d57613d4081613976565b8352918301918301613d2f565b509695505050505050565b600082601f830112613d6957600080fd5b81356020613d79613d0783613cc3565b82815260069290921b84018101918181019086841115613d9857600080fd5b8286015b84811015613d4d57613dae8882613b0a565b835291830191604001613d9c565b600082601f830112613dcd57600080fd5b81356020613ddd613d0783613cc3565b82815260059290921b84018101918181019086841115613dfc57600080fd5b8286015b84811015613d4d5780356001600160401b03811115613e1f5760008081fd5b613e2d8986838b0101613ce6565b845250918301918301613e00565b60006101808284031215613e4e57600080fd5b613e56613ab7565b905081356001600160401b0380821115613e6f57600080fd5b613e7b85838601613ce6565b83526020840135915080821115613e9157600080fd5b613e9d85838601613d58565b60208401526040840135915080821115613eb657600080fd5b613ec285838601613d58565b6040840152613ed48560608601613b8a565b6060840152613ee68560e08601613b0a565b6080840152610120840135915080821115613f0057600080fd5b613f0c85838601613ce6565b60a0840152610140840135915080821115613f2657600080fd5b613f3285838601613ce6565b60c0840152610160840135915080821115613f4c57600080fd5b50613f5984828501613dbc565b60e08301525092915050565b600080600080600060808688031215613f7d57600080fd5b8535945060208601356001600160401b0380821115613f9b57600080fd5b613fa789838a01613901565b9096509450849150613fbb60408901613976565b93506060880135915080821115613fd157600080fd5b50613fde88828901613e3b565b9150509295509295909350565b600081518084526020808501945080840160005b838110156137ba5781516001600160601b031687529582019590820190600101613fff565b604081526000835160408084015261403f6080840182613feb565b90506020850151603f1984830301606085015261405c8282613feb565b925050508260208301529392505050565b60006020828403121561407f57600080fd5b5035919050565b600060e082840312156137d657600080fd5b6000604082840312156137d657600080fd5b6000806000608084860312156140bf57600080fd5b83356001600160401b03808211156140d657600080fd5b6140e287838801614086565b94506140f18760208801614098565b9350606086013591508082111561410757600080fd5b5061411486828701613e3b565b9150509250925092565b8035601781900b81146126e757600080fd5b600080600080600060e0868803121561414857600080fd5b6141518661411e565b945060208601356001600160401b038082111561416d57600080fd5b61417989838a01614086565b95506141888960408a01614098565b94506141978960808a01614098565b935060c08801359150808211156141ad57600080fd5b50613fde88828901613d58565b6000806000606084860312156141cf57600080fd5b83356141da81613949565b925060208401356141ea81613949565b915060408401356141fa81613949565b809150509250925092565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b606081526000614242606083018789614205565b8281036020840152614255818688614205565b9150508260408301529695505050505050565b60005b8381101561428357818101518382015260200161426b565b50506000910152565b600081518084526142a4816020860160208601614268565b601f01601f19169290920160200192915050565b60208152815160208201526000602083015160e060408401526142df61010084018261428c565b905060018060a01b03604085015116606084015263ffffffff60e01b606085015116608084015263ffffffff60808501511660a084015260a0840151601f198483030160c0850152614331828261428c565b91505060c084015161434b60e085018263ffffffff169052565b509392505050565b634e487b7160e01b600052601160045260246000fd5b63ffffffff81811683821601908082111561438657614386614353565b5092915050565b634e487b7160e01b600052603260045260246000fd5b6000826143c057634e487b7160e01b600052601260045260246000fd5b500690565b6000602082840312156143d757600080fd5b815161230081613949565b8385823760008482016000815283858237600093019283525090949350505050565b60006020828403121561441657600080fd5b815160ff8116811461230057600080fd5b63ffffffff60e01b8360e01b1681526000600482018351602080860160005b8381101561446257815185529382019390820190600101614446565b5092979650505050505050565b6000808335601e1984360301811261448657600080fd5b8301803591506001600160401b038211156144a057600080fd5b60200191503681900382131561394257600080fd5b6000808335601e198436030181126144cc57600080fd5b83016020810192503590506001600160401b038111156144eb57600080fd5b80360382131561394257600080fd5b6020815281356020820152600061451460208401846144b5565b60e0604085015261452a61010085018284614205565b915050604084013561453b81613949565b6001600160a01b03166060848101919091526001600160e01b03199061456290860161395e565b16608084015261457460808501613976565b63ffffffff811660a08501525061458e60a08501856144b5565b848303601f190160c08601526145a5838284614205565b925050506145b560c08501613976565b63ffffffff811660e085015261434b565b63ffffffff6145d482613976565b1682526145e36020820161411e565b60170b60208301525050565b6040810161230382846145c6565b6001600160601b03818116838216028082169190828114611c4057611c40614353565b808202811582820484141761230357612303614353565b60006001820161464957614649614353565b5060010190565b6080810161465e82856145c6565b63ffffffff8351166040830152602083015160608301529392505050565b60006020828403121561468e57600080fd5b6123008261411e565b6000602082840312156146a957600080fd5b6123008261395e565b600082516146c4818460208701614268565b9190910192915050565b608081016146dc82856145c6565b63ffffffff6146ea84613976565b166040830152602083013560608301529392505050565b8082018082111561230357612303614353565b8181038181111561230357612303614353565b60006020828403121561473957600080fd5b81516001600160c01b038116811461230057600080fd5b60006020828403121561476257600080fd5b5051919050565b60006020828403121561477b57600080fd5b815167ffffffffffffffff198116811461230057600080fd5b6000602082840312156147a657600080fd5b81516001600160601b038116811461230057600080fd5b6001600160601b0382811682821603908082111561438657614386614353565b818103600083128015838313168383128216171561438657614386614353565b600061ffff80831681810361481457614814614353565b6001019392505050565b600181815b8085111561485957816000190482111561483f5761483f614353565b8085161561484c57918102915b93841c9390800290614823565b509250929050565b60008261487057506001612303565b8161487d57506000612303565b8160018114614893576002811461489d576148b9565b6001915050612303565b60ff8411156148ae576148ae614353565b50506001821b612303565b5060208310610133831016604e8410600b84101617156148dc575081810a612303565b6148e6838361481e565b80600019048211156148fa576148fa614353565b029392505050565b6000612300838361486156fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd474269746d61705574696c732e6f72646572656442797465734172726179546f42424c535369676e6174757265436865636b65722e636865636b5369676e617475a264697066735822122024f34602d6bcbd4a7381d3dcec2fb2794c3d855bf25b4b012683ce211992c7af64736f6c63430008140033",
}

// ContractIncrediblePriceTaskManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIncrediblePriceTaskManagerMetaData.ABI instead.
var ContractIncrediblePriceTaskManagerABI = ContractIncrediblePriceTaskManagerMetaData.ABI

// ContractIncrediblePriceTaskManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractIncrediblePriceTaskManagerMetaData.Bin instead.
var ContractIncrediblePriceTaskManagerBin = ContractIncrediblePriceTaskManagerMetaData.Bin

// DeployContractIncrediblePriceTaskManager deploys a new Ethereum contract, binding an instance of ContractIncrediblePriceTaskManager to it.
func DeployContractIncrediblePriceTaskManager(auth *bind.TransactOpts, backend bind.ContractBackend, _registryCoordinator common.Address, _taskResponseWindowBlock uint32) (common.Address, *types.Transaction, *ContractIncrediblePriceTaskManager, error) {
	parsed, err := ContractIncrediblePriceTaskManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractIncrediblePriceTaskManagerBin), backend, _registryCoordinator, _taskResponseWindowBlock)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractIncrediblePriceTaskManager{ContractIncrediblePriceTaskManagerCaller: ContractIncrediblePriceTaskManagerCaller{contract: contract}, ContractIncrediblePriceTaskManagerTransactor: ContractIncrediblePriceTaskManagerTransactor{contract: contract}, ContractIncrediblePriceTaskManagerFilterer: ContractIncrediblePriceTaskManagerFilterer{contract: contract}}, nil
}

// ContractIncrediblePriceTaskManager is an auto generated Go binding around an Ethereum contract.
type ContractIncrediblePriceTaskManager struct {
	ContractIncrediblePriceTaskManagerCaller     // Read-only binding to the contract
	ContractIncrediblePriceTaskManagerTransactor // Write-only binding to the contract
	ContractIncrediblePriceTaskManagerFilterer   // Log filterer for contract events
}

// ContractIncrediblePriceTaskManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIncrediblePriceTaskManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIncrediblePriceTaskManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIncrediblePriceTaskManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIncrediblePriceTaskManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIncrediblePriceTaskManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIncrediblePriceTaskManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIncrediblePriceTaskManagerSession struct {
	Contract     *ContractIncrediblePriceTaskManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ContractIncrediblePriceTaskManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIncrediblePriceTaskManagerCallerSession struct {
	Contract *ContractIncrediblePriceTaskManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// ContractIncrediblePriceTaskManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIncrediblePriceTaskManagerTransactorSession struct {
	Contract     *ContractIncrediblePriceTaskManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// ContractIncrediblePriceTaskManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIncrediblePriceTaskManagerRaw struct {
	Contract *ContractIncrediblePriceTaskManager // Generic contract binding to access the raw methods on
}

// ContractIncrediblePriceTaskManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIncrediblePriceTaskManagerCallerRaw struct {
	Contract *ContractIncrediblePriceTaskManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIncrediblePriceTaskManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIncrediblePriceTaskManagerTransactorRaw struct {
	Contract *ContractIncrediblePriceTaskManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIncrediblePriceTaskManager creates a new instance of ContractIncrediblePriceTaskManager, bound to a specific deployed contract.
func NewContractIncrediblePriceTaskManager(address common.Address, backend bind.ContractBackend) (*ContractIncrediblePriceTaskManager, error) {
	contract, err := bindContractIncrediblePriceTaskManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManager{ContractIncrediblePriceTaskManagerCaller: ContractIncrediblePriceTaskManagerCaller{contract: contract}, ContractIncrediblePriceTaskManagerTransactor: ContractIncrediblePriceTaskManagerTransactor{contract: contract}, ContractIncrediblePriceTaskManagerFilterer: ContractIncrediblePriceTaskManagerFilterer{contract: contract}}, nil
}

// NewContractIncrediblePriceTaskManagerCaller creates a new read-only instance of ContractIncrediblePriceTaskManager, bound to a specific deployed contract.
func NewContractIncrediblePriceTaskManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractIncrediblePriceTaskManagerCaller, error) {
	contract, err := bindContractIncrediblePriceTaskManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerCaller{contract: contract}, nil
}

// NewContractIncrediblePriceTaskManagerTransactor creates a new write-only instance of ContractIncrediblePriceTaskManager, bound to a specific deployed contract.
func NewContractIncrediblePriceTaskManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIncrediblePriceTaskManagerTransactor, error) {
	contract, err := bindContractIncrediblePriceTaskManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerTransactor{contract: contract}, nil
}

// NewContractIncrediblePriceTaskManagerFilterer creates a new log filterer instance of ContractIncrediblePriceTaskManager, bound to a specific deployed contract.
func NewContractIncrediblePriceTaskManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIncrediblePriceTaskManagerFilterer, error) {
	contract, err := bindContractIncrediblePriceTaskManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerFilterer{contract: contract}, nil
}

// bindContractIncrediblePriceTaskManager binds a generic wrapper to an already deployed contract.
func bindContractIncrediblePriceTaskManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIncrediblePriceTaskManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIncrediblePriceTaskManager.Contract.ContractIncrediblePriceTaskManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.ContractIncrediblePriceTaskManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.ContractIncrediblePriceTaskManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIncrediblePriceTaskManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.contract.Transact(opts, method, params...)
}

// TASKCHALLENGEWINDOWBLOCK is a free data retrieval call binding the contract method 0xf63c5bab.
//
// Solidity: function TASK_CHALLENGE_WINDOW_BLOCK() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) TASKCHALLENGEWINDOWBLOCK(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "TASK_CHALLENGE_WINDOW_BLOCK")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TASKCHALLENGEWINDOWBLOCK is a free data retrieval call binding the contract method 0xf63c5bab.
//
// Solidity: function TASK_CHALLENGE_WINDOW_BLOCK() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) TASKCHALLENGEWINDOWBLOCK() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TASKCHALLENGEWINDOWBLOCK(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// TASKCHALLENGEWINDOWBLOCK is a free data retrieval call binding the contract method 0xf63c5bab.
//
// Solidity: function TASK_CHALLENGE_WINDOW_BLOCK() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) TASKCHALLENGEWINDOWBLOCK() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TASKCHALLENGEWINDOWBLOCK(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// TASKRESPONSEWINDOWBLOCK is a free data retrieval call binding the contract method 0x1ad43189.
//
// Solidity: function TASK_RESPONSE_WINDOW_BLOCK() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) TASKRESPONSEWINDOWBLOCK(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "TASK_RESPONSE_WINDOW_BLOCK")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TASKRESPONSEWINDOWBLOCK is a free data retrieval call binding the contract method 0x1ad43189.
//
// Solidity: function TASK_RESPONSE_WINDOW_BLOCK() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) TASKRESPONSEWINDOWBLOCK() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TASKRESPONSEWINDOWBLOCK(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// TASKRESPONSEWINDOWBLOCK is a free data retrieval call binding the contract method 0x1ad43189.
//
// Solidity: function TASK_RESPONSE_WINDOW_BLOCK() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) TASKRESPONSEWINDOWBLOCK() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TASKRESPONSEWINDOWBLOCK(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) Aggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "aggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) Aggregator() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Aggregator(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// Aggregator is a free data retrieval call binding the contract method 0x245a7bfc.
//
// Solidity: function aggregator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) Aggregator() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Aggregator(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// AllTaskHashes is a free data retrieval call binding the contract method 0x2d89f6fc.
//
// Solidity: function allTaskHashes(uint32 ) view returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) AllTaskHashes(opts *bind.CallOpts, arg0 uint32) ([32]byte, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "allTaskHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AllTaskHashes is a free data retrieval call binding the contract method 0x2d89f6fc.
//
// Solidity: function allTaskHashes(uint32 ) view returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) AllTaskHashes(arg0 uint32) ([32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.AllTaskHashes(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// AllTaskHashes is a free data retrieval call binding the contract method 0x2d89f6fc.
//
// Solidity: function allTaskHashes(uint32 ) view returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) AllTaskHashes(arg0 uint32) ([32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.AllTaskHashes(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// AllTaskResponses is a free data retrieval call binding the contract method 0x2cb223d5.
//
// Solidity: function allTaskResponses(uint32 ) view returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) AllTaskResponses(opts *bind.CallOpts, arg0 uint32) ([32]byte, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "allTaskResponses", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AllTaskResponses is a free data retrieval call binding the contract method 0x2cb223d5.
//
// Solidity: function allTaskResponses(uint32 ) view returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) AllTaskResponses(arg0 uint32) ([32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.AllTaskResponses(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// AllTaskResponses is a free data retrieval call binding the contract method 0x2cb223d5.
//
// Solidity: function allTaskResponses(uint32 ) view returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) AllTaskResponses(arg0 uint32) ([32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.AllTaskResponses(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) BlsApkRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "blsApkRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.BlsApkRegistry(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.BlsApkRegistry(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) CheckSignatures(opts *bind.CallOpts, msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "checkSignatures", msgHash, quorumNumbers, referenceBlockNumber, params)

	if err != nil {
		return *new(IBLSSignatureCheckerQuorumStakeTotals), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSSignatureCheckerQuorumStakeTotals)).(*IBLSSignatureCheckerQuorumStakeTotals)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) CheckSignatures(msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.CheckSignatures(&_ContractIncrediblePriceTaskManager.CallOpts, msgHash, quorumNumbers, referenceBlockNumber, params)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) CheckSignatures(msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.CheckSignatures(&_ContractIncrediblePriceTaskManager.CallOpts, msgHash, quorumNumbers, referenceBlockNumber, params)
}

// Generator is a free data retrieval call binding the contract method 0x7afa1eed.
//
// Solidity: function generator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) Generator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "generator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Generator is a free data retrieval call binding the contract method 0x7afa1eed.
//
// Solidity: function generator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) Generator() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Generator(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// Generator is a free data retrieval call binding the contract method 0x7afa1eed.
//
// Solidity: function generator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) Generator() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Generator(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// GetPriceId is a free data retrieval call binding the contract method 0x4d3141ce.
//
// Solidity: function getPriceId(string baseSymbol, string quoteSymbol) pure returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) GetPriceId(opts *bind.CallOpts, baseSymbol string, quoteSymbol string) ([32]byte, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "getPriceId", baseSymbol, quoteSymbol)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPriceId is a free data retrieval call binding the contract method 0x4d3141ce.
//
// Solidity: function getPriceId(string baseSymbol, string quoteSymbol) pure returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) GetPriceId(baseSymbol string, quoteSymbol string) ([32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.GetPriceId(&_ContractIncrediblePriceTaskManager.CallOpts, baseSymbol, quoteSymbol)
}

// GetPriceId is a free data retrieval call binding the contract method 0x4d3141ce.
//
// Solidity: function getPriceId(string baseSymbol, string quoteSymbol) pure returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) GetPriceId(baseSymbol string, quoteSymbol string) ([32]byte, error) {
	return _ContractIncrediblePriceTaskManager.Contract.GetPriceId(&_ContractIncrediblePriceTaskManager.CallOpts, baseSymbol, quoteSymbol)
}

// GetTaskResponseWindowBlock is a free data retrieval call binding the contract method 0xf5c9899d.
//
// Solidity: function getTaskResponseWindowBlock() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) GetTaskResponseWindowBlock(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "getTaskResponseWindowBlock")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetTaskResponseWindowBlock is a free data retrieval call binding the contract method 0xf5c9899d.
//
// Solidity: function getTaskResponseWindowBlock() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) GetTaskResponseWindowBlock() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.GetTaskResponseWindowBlock(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// GetTaskResponseWindowBlock is a free data retrieval call binding the contract method 0xf5c9899d.
//
// Solidity: function getTaskResponseWindowBlock() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) GetTaskResponseWindowBlock() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.GetTaskResponseWindowBlock(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// LatestPriceInfos is a free data retrieval call binding the contract method 0x9a22de76.
//
// Solidity: function latestPriceInfos(bytes32 ) view returns(int192 price, uint64 publishTime)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) LatestPriceInfos(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Price       *big.Int
	PublishTime uint64
}, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "latestPriceInfos", arg0)

	outstruct := new(struct {
		Price       *big.Int
		PublishTime uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Price = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.PublishTime = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// LatestPriceInfos is a free data retrieval call binding the contract method 0x9a22de76.
//
// Solidity: function latestPriceInfos(bytes32 ) view returns(int192 price, uint64 publishTime)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) LatestPriceInfos(arg0 [32]byte) (struct {
	Price       *big.Int
	PublishTime uint64
}, error) {
	return _ContractIncrediblePriceTaskManager.Contract.LatestPriceInfos(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// LatestPriceInfos is a free data retrieval call binding the contract method 0x9a22de76.
//
// Solidity: function latestPriceInfos(bytes32 ) view returns(int192 price, uint64 publishTime)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) LatestPriceInfos(arg0 [32]byte) (struct {
	Price       *big.Int
	PublishTime uint64
}, error) {
	return _ContractIncrediblePriceTaskManager.Contract.LatestPriceInfos(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// LatestTaskNum is a free data retrieval call binding the contract method 0x8b00ce7c.
//
// Solidity: function latestTaskNum() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) LatestTaskNum(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "latestTaskNum")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LatestTaskNum is a free data retrieval call binding the contract method 0x8b00ce7c.
//
// Solidity: function latestTaskNum() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) LatestTaskNum() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.LatestTaskNum(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// LatestTaskNum is a free data retrieval call binding the contract method 0x8b00ce7c.
//
// Solidity: function latestTaskNum() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) LatestTaskNum() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.LatestTaskNum(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// OmniOperatorSharesManager is a free data retrieval call binding the contract method 0x66cef710.
//
// Solidity: function omniOperatorSharesManager() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) OmniOperatorSharesManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "omniOperatorSharesManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OmniOperatorSharesManager is a free data retrieval call binding the contract method 0x66cef710.
//
// Solidity: function omniOperatorSharesManager() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) OmniOperatorSharesManager() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.OmniOperatorSharesManager(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// OmniOperatorSharesManager is a free data retrieval call binding the contract method 0x66cef710.
//
// Solidity: function omniOperatorSharesManager() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) OmniOperatorSharesManager() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.OmniOperatorSharesManager(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) Owner() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Owner(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) Owner() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Owner(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RegistryCoordinator(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RegistryCoordinator(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) StakeRegistry() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.StakeRegistry(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) StakeRegistry() (common.Address, error) {
	return _ContractIncrediblePriceTaskManager.Contract.StakeRegistry(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) StaleStakesForbidden(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "staleStakesForbidden")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) StaleStakesForbidden() (bool, error) {
	return _ContractIncrediblePriceTaskManager.Contract.StaleStakesForbidden(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) StaleStakesForbidden() (bool, error) {
	return _ContractIncrediblePriceTaskManager.Contract.StaleStakesForbidden(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) TaskNumber(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "taskNumber")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) TaskNumber() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TaskNumber(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) TaskNumber() (uint32, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TaskNumber(&_ContractIncrediblePriceTaskManager.CallOpts)
}

// TaskSuccesfullyChallenged is a free data retrieval call binding the contract method 0x5decc3f5.
//
// Solidity: function taskSuccesfullyChallenged(uint32 ) view returns(bool)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) TaskSuccesfullyChallenged(opts *bind.CallOpts, arg0 uint32) (bool, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "taskSuccesfullyChallenged", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TaskSuccesfullyChallenged is a free data retrieval call binding the contract method 0x5decc3f5.
//
// Solidity: function taskSuccesfullyChallenged(uint32 ) view returns(bool)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) TaskSuccesfullyChallenged(arg0 uint32) (bool, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TaskSuccesfullyChallenged(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// TaskSuccesfullyChallenged is a free data retrieval call binding the contract method 0x5decc3f5.
//
// Solidity: function taskSuccesfullyChallenged(uint32 ) view returns(bool)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) TaskSuccesfullyChallenged(arg0 uint32) (bool, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TaskSuccesfullyChallenged(&_ContractIncrediblePriceTaskManager.CallOpts, arg0)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCaller) TrySignatureAndApkVerification(opts *bind.CallOpts, msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	var out []interface{}
	err := _ContractIncrediblePriceTaskManager.contract.Call(opts, &out, "trySignatureAndApkVerification", msgHash, apk, apkG2, sigma)

	outstruct := new(struct {
		PairingSuccessful bool
		SiganatureIsValid bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PairingSuccessful = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.SiganatureIsValid = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TrySignatureAndApkVerification(&_ContractIncrediblePriceTaskManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerCallerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TrySignatureAndApkVerification(&_ContractIncrediblePriceTaskManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _aggregator, address _generator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, _aggregator common.Address, _generator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "initialize", initialOwner, _aggregator, _generator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _aggregator, address _generator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) Initialize(initialOwner common.Address, _aggregator common.Address, _generator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Initialize(&_ContractIncrediblePriceTaskManager.TransactOpts, initialOwner, _aggregator, _generator)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address _aggregator, address _generator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) Initialize(initialOwner common.Address, _aggregator common.Address, _generator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.Initialize(&_ContractIncrediblePriceTaskManager.TransactOpts, initialOwner, _aggregator, _generator)
}

// RaiseAndResolveChallenge is a paid mutator transaction binding the contract method 0xbdaf8658.
//
// Solidity: function raiseAndResolveChallenge(int192 actualPrice, (bytes32,bytes,address,bytes4,uint32,bytes,uint32) task, (uint32,int192) taskResponse, (uint32,bytes32) taskResponseMetadata, (uint256,uint256)[] pubkeysOfNonSigningOperators) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) RaiseAndResolveChallenge(opts *bind.TransactOpts, actualPrice *big.Int, task IIncrediblePriceTaskManagerTask, taskResponse IIncrediblePriceTaskManagerTaskResponse, taskResponseMetadata IIncrediblePriceTaskManagerTaskResponseMetadata, pubkeysOfNonSigningOperators []BN254G1Point) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "raiseAndResolveChallenge", actualPrice, task, taskResponse, taskResponseMetadata, pubkeysOfNonSigningOperators)
}

// RaiseAndResolveChallenge is a paid mutator transaction binding the contract method 0xbdaf8658.
//
// Solidity: function raiseAndResolveChallenge(int192 actualPrice, (bytes32,bytes,address,bytes4,uint32,bytes,uint32) task, (uint32,int192) taskResponse, (uint32,bytes32) taskResponseMetadata, (uint256,uint256)[] pubkeysOfNonSigningOperators) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) RaiseAndResolveChallenge(actualPrice *big.Int, task IIncrediblePriceTaskManagerTask, taskResponse IIncrediblePriceTaskManagerTaskResponse, taskResponseMetadata IIncrediblePriceTaskManagerTaskResponseMetadata, pubkeysOfNonSigningOperators []BN254G1Point) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RaiseAndResolveChallenge(&_ContractIncrediblePriceTaskManager.TransactOpts, actualPrice, task, taskResponse, taskResponseMetadata, pubkeysOfNonSigningOperators)
}

// RaiseAndResolveChallenge is a paid mutator transaction binding the contract method 0xbdaf8658.
//
// Solidity: function raiseAndResolveChallenge(int192 actualPrice, (bytes32,bytes,address,bytes4,uint32,bytes,uint32) task, (uint32,int192) taskResponse, (uint32,bytes32) taskResponseMetadata, (uint256,uint256)[] pubkeysOfNonSigningOperators) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) RaiseAndResolveChallenge(actualPrice *big.Int, task IIncrediblePriceTaskManagerTask, taskResponse IIncrediblePriceTaskManagerTaskResponse, taskResponseMetadata IIncrediblePriceTaskManagerTaskResponseMetadata, pubkeysOfNonSigningOperators []BN254G1Point) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RaiseAndResolveChallenge(&_ContractIncrediblePriceTaskManager.TransactOpts, actualPrice, task, taskResponse, taskResponseMetadata, pubkeysOfNonSigningOperators)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RenounceOwnership(&_ContractIncrediblePriceTaskManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RenounceOwnership(&_ContractIncrediblePriceTaskManager.TransactOpts)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x05880a65.
//
// Solidity: function requestPrice(string baseSymbol, string quoteSymbol, address callbackAddress, bytes4 callbackFunctionId, uint32 quorumThresholdPercentage, bytes quorumNumbers) returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) RequestPrice(opts *bind.TransactOpts, baseSymbol string, quoteSymbol string, callbackAddress common.Address, callbackFunctionId [4]byte, quorumThresholdPercentage uint32, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "requestPrice", baseSymbol, quoteSymbol, callbackAddress, callbackFunctionId, quorumThresholdPercentage, quorumNumbers)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x05880a65.
//
// Solidity: function requestPrice(string baseSymbol, string quoteSymbol, address callbackAddress, bytes4 callbackFunctionId, uint32 quorumThresholdPercentage, bytes quorumNumbers) returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) RequestPrice(baseSymbol string, quoteSymbol string, callbackAddress common.Address, callbackFunctionId [4]byte, quorumThresholdPercentage uint32, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RequestPrice(&_ContractIncrediblePriceTaskManager.TransactOpts, baseSymbol, quoteSymbol, callbackAddress, callbackFunctionId, quorumThresholdPercentage, quorumNumbers)
}

// RequestPrice is a paid mutator transaction binding the contract method 0x05880a65.
//
// Solidity: function requestPrice(string baseSymbol, string quoteSymbol, address callbackAddress, bytes4 callbackFunctionId, uint32 quorumThresholdPercentage, bytes quorumNumbers) returns(bytes32)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) RequestPrice(baseSymbol string, quoteSymbol string, callbackAddress common.Address, callbackFunctionId [4]byte, quorumThresholdPercentage uint32, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.RequestPrice(&_ContractIncrediblePriceTaskManager.TransactOpts, baseSymbol, quoteSymbol, callbackAddress, callbackFunctionId, quorumThresholdPercentage, quorumNumbers)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) SetStaleStakesForbidden(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "setStaleStakesForbidden", value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.SetStaleStakesForbidden(&_ContractIncrediblePriceTaskManager.TransactOpts, value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.SetStaleStakesForbidden(&_ContractIncrediblePriceTaskManager.TransactOpts, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TransferOwnership(&_ContractIncrediblePriceTaskManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.TransferOwnership(&_ContractIncrediblePriceTaskManager.TransactOpts, newOwner)
}

// UpdateAggregator is a paid mutator transaction binding the contract method 0x9fe4ee47.
//
// Solidity: function updateAggregator(address _aggregator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) UpdateAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "updateAggregator", _aggregator)
}

// UpdateAggregator is a paid mutator transaction binding the contract method 0x9fe4ee47.
//
// Solidity: function updateAggregator(address _aggregator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) UpdateAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.UpdateAggregator(&_ContractIncrediblePriceTaskManager.TransactOpts, _aggregator)
}

// UpdateAggregator is a paid mutator transaction binding the contract method 0x9fe4ee47.
//
// Solidity: function updateAggregator(address _aggregator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) UpdateAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.UpdateAggregator(&_ContractIncrediblePriceTaskManager.TransactOpts, _aggregator)
}

// UpdateGenerator is a paid mutator transaction binding the contract method 0x0ed21268.
//
// Solidity: function updateGenerator(address _generator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) UpdateGenerator(opts *bind.TransactOpts, _generator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "updateGenerator", _generator)
}

// UpdateGenerator is a paid mutator transaction binding the contract method 0x0ed21268.
//
// Solidity: function updateGenerator(address _generator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) UpdateGenerator(_generator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.UpdateGenerator(&_ContractIncrediblePriceTaskManager.TransactOpts, _generator)
}

// UpdateGenerator is a paid mutator transaction binding the contract method 0x0ed21268.
//
// Solidity: function updateGenerator(address _generator) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) UpdateGenerator(_generator common.Address) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.UpdateGenerator(&_ContractIncrediblePriceTaskManager.TransactOpts, _generator)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0xa5f3c699.
//
// Solidity: function updatePrice((bytes32,bytes,address,bytes4,uint32,bytes,uint32) task, (uint32,int192) taskResponse, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactor) UpdatePrice(opts *bind.TransactOpts, task IIncrediblePriceTaskManagerTask, taskResponse IIncrediblePriceTaskManagerTaskResponse, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.contract.Transact(opts, "updatePrice", task, taskResponse, nonSignerStakesAndSignature)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0xa5f3c699.
//
// Solidity: function updatePrice((bytes32,bytes,address,bytes4,uint32,bytes,uint32) task, (uint32,int192) taskResponse, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerSession) UpdatePrice(task IIncrediblePriceTaskManagerTask, taskResponse IIncrediblePriceTaskManagerTaskResponse, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.UpdatePrice(&_ContractIncrediblePriceTaskManager.TransactOpts, task, taskResponse, nonSignerStakesAndSignature)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0xa5f3c699.
//
// Solidity: function updatePrice((bytes32,bytes,address,bytes4,uint32,bytes,uint32) task, (uint32,int192) taskResponse, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerTransactorSession) UpdatePrice(task IIncrediblePriceTaskManagerTask, taskResponse IIncrediblePriceTaskManagerTaskResponse, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIncrediblePriceTaskManager.Contract.UpdatePrice(&_ContractIncrediblePriceTaskManager.TransactOpts, task, taskResponse, nonSignerStakesAndSignature)
}

// ContractIncrediblePriceTaskManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerInitializedIterator struct {
	Event *ContractIncrediblePriceTaskManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerInitialized represents a Initialized event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractIncrediblePriceTaskManagerInitializedIterator, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerInitializedIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerInitialized)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseInitialized(log types.Log) (*ContractIncrediblePriceTaskManagerInitialized, error) {
	event := new(ContractIncrediblePriceTaskManagerInitialized)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerNewTaskCreatedIterator is returned from FilterNewTaskCreated and is used to iterate over the raw logs and unpacked data for NewTaskCreated events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerNewTaskCreatedIterator struct {
	Event *ContractIncrediblePriceTaskManagerNewTaskCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerNewTaskCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerNewTaskCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerNewTaskCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerNewTaskCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerNewTaskCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerNewTaskCreated represents a NewTaskCreated event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerNewTaskCreated struct {
	TaskIndex uint32
	Task      IIncrediblePriceTaskManagerTask
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNewTaskCreated is a free log retrieval operation binding the contract event 0x36e52b7ae44ca8f480b92f4c6107dcd429968909437e0bb327dca642b85d0a30.
//
// Solidity: event NewTaskCreated(uint32 indexed taskIndex, (bytes32,bytes,address,bytes4,uint32,bytes,uint32) task)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterNewTaskCreated(opts *bind.FilterOpts, taskIndex []uint32) (*ContractIncrediblePriceTaskManagerNewTaskCreatedIterator, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "NewTaskCreated", taskIndexRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerNewTaskCreatedIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "NewTaskCreated", logs: logs, sub: sub}, nil
}

// WatchNewTaskCreated is a free log subscription operation binding the contract event 0x36e52b7ae44ca8f480b92f4c6107dcd429968909437e0bb327dca642b85d0a30.
//
// Solidity: event NewTaskCreated(uint32 indexed taskIndex, (bytes32,bytes,address,bytes4,uint32,bytes,uint32) task)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchNewTaskCreated(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerNewTaskCreated, taskIndex []uint32) (event.Subscription, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "NewTaskCreated", taskIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerNewTaskCreated)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "NewTaskCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewTaskCreated is a log parse operation binding the contract event 0x36e52b7ae44ca8f480b92f4c6107dcd429968909437e0bb327dca642b85d0a30.
//
// Solidity: event NewTaskCreated(uint32 indexed taskIndex, (bytes32,bytes,address,bytes4,uint32,bytes,uint32) task)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseNewTaskCreated(log types.Log) (*ContractIncrediblePriceTaskManagerNewTaskCreated, error) {
	event := new(ContractIncrediblePriceTaskManagerNewTaskCreated)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "NewTaskCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerOwnershipTransferredIterator struct {
	Event *ContractIncrediblePriceTaskManagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractIncrediblePriceTaskManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerOwnershipTransferredIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerOwnershipTransferred)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractIncrediblePriceTaskManagerOwnershipTransferred, error) {
	event := new(ContractIncrediblePriceTaskManagerOwnershipTransferred)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerRequestPriceIterator is returned from FilterRequestPrice and is used to iterate over the raw logs and unpacked data for RequestPrice events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerRequestPriceIterator struct {
	Event *ContractIncrediblePriceTaskManagerRequestPrice // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerRequestPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerRequestPrice)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerRequestPrice)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerRequestPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerRequestPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerRequestPrice represents a RequestPrice event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerRequestPrice struct {
	Sender      common.Address
	PriceId     [32]byte
	BaseSymbol  string
	QuoteSymbol string
	Times       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRequestPrice is a free log retrieval operation binding the contract event 0xcffe5e47d40a1edafcd79aaafab384bbdcdf0c48c2b6c32fe79dcdb89ef3023d.
//
// Solidity: event RequestPrice(address indexed sender, bytes32 indexed priceId, string baseSymbol, string quoteSymbol, int256 times)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterRequestPrice(opts *bind.FilterOpts, sender []common.Address, priceId [][32]byte) (*ContractIncrediblePriceTaskManagerRequestPriceIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var priceIdRule []interface{}
	for _, priceIdItem := range priceId {
		priceIdRule = append(priceIdRule, priceIdItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "RequestPrice", senderRule, priceIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerRequestPriceIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "RequestPrice", logs: logs, sub: sub}, nil
}

// WatchRequestPrice is a free log subscription operation binding the contract event 0xcffe5e47d40a1edafcd79aaafab384bbdcdf0c48c2b6c32fe79dcdb89ef3023d.
//
// Solidity: event RequestPrice(address indexed sender, bytes32 indexed priceId, string baseSymbol, string quoteSymbol, int256 times)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchRequestPrice(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerRequestPrice, sender []common.Address, priceId [][32]byte) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var priceIdRule []interface{}
	for _, priceIdItem := range priceId {
		priceIdRule = append(priceIdRule, priceIdItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "RequestPrice", senderRule, priceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerRequestPrice)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "RequestPrice", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestPrice is a log parse operation binding the contract event 0xcffe5e47d40a1edafcd79aaafab384bbdcdf0c48c2b6c32fe79dcdb89ef3023d.
//
// Solidity: event RequestPrice(address indexed sender, bytes32 indexed priceId, string baseSymbol, string quoteSymbol, int256 times)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseRequestPrice(log types.Log) (*ContractIncrediblePriceTaskManagerRequestPrice, error) {
	event := new(ContractIncrediblePriceTaskManagerRequestPrice)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "RequestPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator is returned from FilterStaleStakesForbiddenUpdate and is used to iterate over the raw logs and unpacked data for StaleStakesForbiddenUpdate events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator struct {
	Event *ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate represents a StaleStakesForbiddenUpdate event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate struct {
	Value bool
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterStaleStakesForbiddenUpdate is a free log retrieval operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterStaleStakesForbiddenUpdate(opts *bind.FilterOpts) (*ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdateIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "StaleStakesForbiddenUpdate", logs: logs, sub: sub}, nil
}

// WatchStaleStakesForbiddenUpdate is a free log subscription operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchStaleStakesForbiddenUpdate(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate) (event.Subscription, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStaleStakesForbiddenUpdate is a log parse operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseStaleStakesForbiddenUpdate(log types.Log) (*ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate, error) {
	event := new(ContractIncrediblePriceTaskManagerStaleStakesForbiddenUpdate)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator is returned from FilterTaskChallengedSuccessfully and is used to iterate over the raw logs and unpacked data for TaskChallengedSuccessfully events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator struct {
	Event *ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully represents a TaskChallengedSuccessfully event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully struct {
	TaskIndex  uint32
	Challenger common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTaskChallengedSuccessfully is a free log retrieval operation binding the contract event 0xc20d1bb0f1623680306b83d4ff4bb99a2beb9d86d97832f3ca40fd13a29df1ec.
//
// Solidity: event TaskChallengedSuccessfully(uint32 indexed taskIndex, address indexed challenger)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterTaskChallengedSuccessfully(opts *bind.FilterOpts, taskIndex []uint32, challenger []common.Address) (*ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "TaskChallengedSuccessfully", taskIndexRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerTaskChallengedSuccessfullyIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "TaskChallengedSuccessfully", logs: logs, sub: sub}, nil
}

// WatchTaskChallengedSuccessfully is a free log subscription operation binding the contract event 0xc20d1bb0f1623680306b83d4ff4bb99a2beb9d86d97832f3ca40fd13a29df1ec.
//
// Solidity: event TaskChallengedSuccessfully(uint32 indexed taskIndex, address indexed challenger)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchTaskChallengedSuccessfully(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully, taskIndex []uint32, challenger []common.Address) (event.Subscription, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "TaskChallengedSuccessfully", taskIndexRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskChallengedSuccessfully", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskChallengedSuccessfully is a log parse operation binding the contract event 0xc20d1bb0f1623680306b83d4ff4bb99a2beb9d86d97832f3ca40fd13a29df1ec.
//
// Solidity: event TaskChallengedSuccessfully(uint32 indexed taskIndex, address indexed challenger)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseTaskChallengedSuccessfully(log types.Log) (*ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully, error) {
	event := new(ContractIncrediblePriceTaskManagerTaskChallengedSuccessfully)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskChallengedSuccessfully", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator is returned from FilterTaskChallengedUnsuccessfully and is used to iterate over the raw logs and unpacked data for TaskChallengedUnsuccessfully events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator struct {
	Event *ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully represents a TaskChallengedUnsuccessfully event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully struct {
	TaskIndex  uint32
	Challenger common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTaskChallengedUnsuccessfully is a free log retrieval operation binding the contract event 0xfd3e26beeb5967fc5a57a0446914eabc45b4aa474c67a51b4b5160cac60ddb05.
//
// Solidity: event TaskChallengedUnsuccessfully(uint32 indexed taskIndex, address indexed challenger)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterTaskChallengedUnsuccessfully(opts *bind.FilterOpts, taskIndex []uint32, challenger []common.Address) (*ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "TaskChallengedUnsuccessfully", taskIndexRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfullyIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "TaskChallengedUnsuccessfully", logs: logs, sub: sub}, nil
}

// WatchTaskChallengedUnsuccessfully is a free log subscription operation binding the contract event 0xfd3e26beeb5967fc5a57a0446914eabc45b4aa474c67a51b4b5160cac60ddb05.
//
// Solidity: event TaskChallengedUnsuccessfully(uint32 indexed taskIndex, address indexed challenger)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchTaskChallengedUnsuccessfully(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully, taskIndex []uint32, challenger []common.Address) (event.Subscription, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}
	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "TaskChallengedUnsuccessfully", taskIndexRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskChallengedUnsuccessfully", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskChallengedUnsuccessfully is a log parse operation binding the contract event 0xfd3e26beeb5967fc5a57a0446914eabc45b4aa474c67a51b4b5160cac60ddb05.
//
// Solidity: event TaskChallengedUnsuccessfully(uint32 indexed taskIndex, address indexed challenger)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseTaskChallengedUnsuccessfully(log types.Log) (*ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully, error) {
	event := new(ContractIncrediblePriceTaskManagerTaskChallengedUnsuccessfully)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskChallengedUnsuccessfully", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerTaskCompletedIterator is returned from FilterTaskCompleted and is used to iterate over the raw logs and unpacked data for TaskCompleted events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskCompletedIterator struct {
	Event *ContractIncrediblePriceTaskManagerTaskCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerTaskCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerTaskCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerTaskCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerTaskCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerTaskCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerTaskCompleted represents a TaskCompleted event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskCompleted struct {
	TaskIndex uint32
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTaskCompleted is a free log retrieval operation binding the contract event 0x9a144f228a931b9d0d1696fbcdaf310b24b5d2d21e799db623fc986a0f547430.
//
// Solidity: event TaskCompleted(uint32 indexed taskIndex)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterTaskCompleted(opts *bind.FilterOpts, taskIndex []uint32) (*ContractIncrediblePriceTaskManagerTaskCompletedIterator, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "TaskCompleted", taskIndexRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerTaskCompletedIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "TaskCompleted", logs: logs, sub: sub}, nil
}

// WatchTaskCompleted is a free log subscription operation binding the contract event 0x9a144f228a931b9d0d1696fbcdaf310b24b5d2d21e799db623fc986a0f547430.
//
// Solidity: event TaskCompleted(uint32 indexed taskIndex)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchTaskCompleted(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerTaskCompleted, taskIndex []uint32) (event.Subscription, error) {

	var taskIndexRule []interface{}
	for _, taskIndexItem := range taskIndex {
		taskIndexRule = append(taskIndexRule, taskIndexItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "TaskCompleted", taskIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerTaskCompleted)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskCompleted is a log parse operation binding the contract event 0x9a144f228a931b9d0d1696fbcdaf310b24b5d2d21e799db623fc986a0f547430.
//
// Solidity: event TaskCompleted(uint32 indexed taskIndex)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseTaskCompleted(log types.Log) (*ContractIncrediblePriceTaskManagerTaskCompleted, error) {
	event := new(ContractIncrediblePriceTaskManagerTaskCompleted)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerTaskRespondedIterator is returned from FilterTaskResponded and is used to iterate over the raw logs and unpacked data for TaskResponded events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskRespondedIterator struct {
	Event *ContractIncrediblePriceTaskManagerTaskResponded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerTaskRespondedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerTaskResponded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerTaskResponded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerTaskRespondedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerTaskRespondedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerTaskResponded represents a TaskResponded event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerTaskResponded struct {
	TaskResponse         IIncrediblePriceTaskManagerTaskResponse
	TaskResponseMetadata IIncrediblePriceTaskManagerTaskResponseMetadata
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterTaskResponded is a free log retrieval operation binding the contract event 0x68b0dddd02e8a85939b35ac692e68162f8a1dbdd71737d8b20620abbedff5c3d.
//
// Solidity: event TaskResponded((uint32,int192) taskResponse, (uint32,bytes32) taskResponseMetadata)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterTaskResponded(opts *bind.FilterOpts) (*ContractIncrediblePriceTaskManagerTaskRespondedIterator, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "TaskResponded")
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerTaskRespondedIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "TaskResponded", logs: logs, sub: sub}, nil
}

// WatchTaskResponded is a free log subscription operation binding the contract event 0x68b0dddd02e8a85939b35ac692e68162f8a1dbdd71737d8b20620abbedff5c3d.
//
// Solidity: event TaskResponded((uint32,int192) taskResponse, (uint32,bytes32) taskResponseMetadata)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchTaskResponded(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerTaskResponded) (event.Subscription, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "TaskResponded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerTaskResponded)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskResponded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskResponded is a log parse operation binding the contract event 0x68b0dddd02e8a85939b35ac692e68162f8a1dbdd71737d8b20620abbedff5c3d.
//
// Solidity: event TaskResponded((uint32,int192) taskResponse, (uint32,bytes32) taskResponseMetadata)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseTaskResponded(log types.Log) (*ContractIncrediblePriceTaskManagerTaskResponded, error) {
	event := new(ContractIncrediblePriceTaskManagerTaskResponded)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "TaskResponded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerUpdateAggregatorIterator is returned from FilterUpdateAggregator and is used to iterate over the raw logs and unpacked data for UpdateAggregator events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerUpdateAggregatorIterator struct {
	Event *ContractIncrediblePriceTaskManagerUpdateAggregator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerUpdateAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerUpdateAggregator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerUpdateAggregator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerUpdateAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerUpdateAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerUpdateAggregator represents a UpdateAggregator event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerUpdateAggregator struct {
	PreviousAggregator common.Address
	NewAggregator      common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpdateAggregator is a free log retrieval operation binding the contract event 0x35a4ee0f60435b1bc206614dffe6bfd98db7e7450ab2ff2da4611d15ce8ffd69.
//
// Solidity: event UpdateAggregator(address previousAggregator, address newAggregator)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterUpdateAggregator(opts *bind.FilterOpts) (*ContractIncrediblePriceTaskManagerUpdateAggregatorIterator, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "UpdateAggregator")
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerUpdateAggregatorIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "UpdateAggregator", logs: logs, sub: sub}, nil
}

// WatchUpdateAggregator is a free log subscription operation binding the contract event 0x35a4ee0f60435b1bc206614dffe6bfd98db7e7450ab2ff2da4611d15ce8ffd69.
//
// Solidity: event UpdateAggregator(address previousAggregator, address newAggregator)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchUpdateAggregator(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerUpdateAggregator) (event.Subscription, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "UpdateAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerUpdateAggregator)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "UpdateAggregator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateAggregator is a log parse operation binding the contract event 0x35a4ee0f60435b1bc206614dffe6bfd98db7e7450ab2ff2da4611d15ce8ffd69.
//
// Solidity: event UpdateAggregator(address previousAggregator, address newAggregator)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseUpdateAggregator(log types.Log) (*ContractIncrediblePriceTaskManagerUpdateAggregator, error) {
	event := new(ContractIncrediblePriceTaskManagerUpdateAggregator)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "UpdateAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerUpdateGeneratorIterator is returned from FilterUpdateGenerator and is used to iterate over the raw logs and unpacked data for UpdateGenerator events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerUpdateGeneratorIterator struct {
	Event *ContractIncrediblePriceTaskManagerUpdateGenerator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerUpdateGeneratorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerUpdateGenerator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerUpdateGenerator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerUpdateGeneratorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerUpdateGeneratorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerUpdateGenerator represents a UpdateGenerator event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerUpdateGenerator struct {
	PreviousGenerator common.Address
	NewGenerator      common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpdateGenerator is a free log retrieval operation binding the contract event 0x774efe45a4bbe359a5d19d71e74cbf8a7215d5b4ae22694324fa7f5d136cd678.
//
// Solidity: event UpdateGenerator(address previousGenerator, address newGenerator)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterUpdateGenerator(opts *bind.FilterOpts) (*ContractIncrediblePriceTaskManagerUpdateGeneratorIterator, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "UpdateGenerator")
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerUpdateGeneratorIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "UpdateGenerator", logs: logs, sub: sub}, nil
}

// WatchUpdateGenerator is a free log subscription operation binding the contract event 0x774efe45a4bbe359a5d19d71e74cbf8a7215d5b4ae22694324fa7f5d136cd678.
//
// Solidity: event UpdateGenerator(address previousGenerator, address newGenerator)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchUpdateGenerator(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerUpdateGenerator) (event.Subscription, error) {

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "UpdateGenerator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerUpdateGenerator)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "UpdateGenerator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateGenerator is a log parse operation binding the contract event 0x774efe45a4bbe359a5d19d71e74cbf8a7215d5b4ae22694324fa7f5d136cd678.
//
// Solidity: event UpdateGenerator(address previousGenerator, address newGenerator)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseUpdateGenerator(log types.Log) (*ContractIncrediblePriceTaskManagerUpdateGenerator, error) {
	event := new(ContractIncrediblePriceTaskManagerUpdateGenerator)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "UpdateGenerator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIncrediblePriceTaskManagerUpdatePriceIterator is returned from FilterUpdatePrice and is used to iterate over the raw logs and unpacked data for UpdatePrice events raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerUpdatePriceIterator struct {
	Event *ContractIncrediblePriceTaskManagerUpdatePrice // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncrediblePriceTaskManagerUpdatePriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncrediblePriceTaskManagerUpdatePrice)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncrediblePriceTaskManagerUpdatePrice)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncrediblePriceTaskManagerUpdatePriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncrediblePriceTaskManagerUpdatePriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncrediblePriceTaskManagerUpdatePrice represents a UpdatePrice event raised by the ContractIncrediblePriceTaskManager contract.
type ContractIncrediblePriceTaskManagerUpdatePrice struct {
	PriceId     [32]byte
	Price       *big.Int
	PublishTime uint32
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdatePrice is a free log retrieval operation binding the contract event 0x3502e3be13f3f611c5bd255b5ee64196626f6ebaac138cf0f4db61166976e553.
//
// Solidity: event UpdatePrice(bytes32 indexed priceId, int192 price, uint32 publishTime)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) FilterUpdatePrice(opts *bind.FilterOpts, priceId [][32]byte) (*ContractIncrediblePriceTaskManagerUpdatePriceIterator, error) {

	var priceIdRule []interface{}
	for _, priceIdItem := range priceId {
		priceIdRule = append(priceIdRule, priceIdItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.FilterLogs(opts, "UpdatePrice", priceIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncrediblePriceTaskManagerUpdatePriceIterator{contract: _ContractIncrediblePriceTaskManager.contract, event: "UpdatePrice", logs: logs, sub: sub}, nil
}

// WatchUpdatePrice is a free log subscription operation binding the contract event 0x3502e3be13f3f611c5bd255b5ee64196626f6ebaac138cf0f4db61166976e553.
//
// Solidity: event UpdatePrice(bytes32 indexed priceId, int192 price, uint32 publishTime)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) WatchUpdatePrice(opts *bind.WatchOpts, sink chan<- *ContractIncrediblePriceTaskManagerUpdatePrice, priceId [][32]byte) (event.Subscription, error) {

	var priceIdRule []interface{}
	for _, priceIdItem := range priceId {
		priceIdRule = append(priceIdRule, priceIdItem)
	}

	logs, sub, err := _ContractIncrediblePriceTaskManager.contract.WatchLogs(opts, "UpdatePrice", priceIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncrediblePriceTaskManagerUpdatePrice)
				if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "UpdatePrice", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatePrice is a log parse operation binding the contract event 0x3502e3be13f3f611c5bd255b5ee64196626f6ebaac138cf0f4db61166976e553.
//
// Solidity: event UpdatePrice(bytes32 indexed priceId, int192 price, uint32 publishTime)
func (_ContractIncrediblePriceTaskManager *ContractIncrediblePriceTaskManagerFilterer) ParseUpdatePrice(log types.Log) (*ContractIncrediblePriceTaskManagerUpdatePrice, error) {
	event := new(ContractIncrediblePriceTaskManagerUpdatePrice)
	if err := _ContractIncrediblePriceTaskManager.contract.UnpackLog(event, "UpdatePrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
