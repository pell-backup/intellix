package server

import (
	"context"
	"fmt"
	"intellix/dvs/vrf/types"
	taskgateway "intellix/gateway/submitter"
)

var GatewayClient *taskgateway.Client

func (s Server) DVSResponsHandler(ctx context.Context, in *types.GenerateRandomNumberRequest) (*types.DVSResultResponse, error) {
	s.logger.Debug("DVSResponsHandler",
		"TaskIndex", in.TaskMetadata.TaskIndex,
		"TaskMetadata", fmt.Sprintf("%+v", in.TaskMetadata),
	)
	//pkgCtx := sdktypes.UnwrapContext(ctx)
	//
	//validatedData, err := pelldvs.GetDvsRequestValidatedData(pkgCtx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Convert []uint32 to bytes
	//groupNumbersBytes := make([]byte, len(in.TaskMetadata.GroupNumbers))
	//for i, num := range in.TaskMetadata.GroupNumbers {
	//	groupNumbersBytes[i] = byte(num)
	//}
	//
	//squared, _ := math.NewIntFromString(string(validatedData.Data))
	//// Construct task parameters
	//task := csquaringmanager.IIncredibleSquaringServiceManagerTask{
	//	NumberToBeSquared:        in.TaskMetadata.Squared.BigInt(),
	//	TaskCreatedBlock:         in.TaskMetadata.Height,
	//	GroupNumbers:             groupNumbersBytes,
	//	GroupThresholdPercentage: in.TaskMetadata.GroupThresholdPercentage,
	//}
	//
	//// Construct TaskResponse parameters
	//taskResponse := csquaringmanager.IIncredibleSquaringServiceManagerTaskResponse{
	//	ReferenceTaskIndex: in.TaskMetadata.TaskIndex,
	//	NumberSquared:      squared.BigInt(),
	//}
	//
	//// Construct NonSignerStakesAndSignature parameters
	//nonSignerPubkeysG1 := make([]csquaringmanager.BN254G1Point, len(validatedData.NonSignersPubkeysG1))
	//for i, pubkey := range validatedData.NonSignersPubkeysG1 {
	//	nonSignerPubkeysG1[i] = csquaringmanager.BN254G1Point{
	//		X: new(big.Int).SetBytes(pubkey[:32]),
	//		Y: new(big.Int).SetBytes(pubkey[32:]),
	//	}
	//}
	//
	//quorumApksG1 := []csquaringmanager.BN254G1Point{}
	//for _, apk := range validatedData.QuorumApksG1 {
	//	tapk := bls.NewZeroG1Point()
	//	_ = tapk.Unmarshal(apk)
	//	quorumApksG1 = append(quorumApksG1, csquaringmanager.BN254G1Point{
	//		X: tapk.X.BigInt(big.NewInt(0)),
	//		Y: tapk.Y.BigInt(big.NewInt(0)),
	//	})
	//}
	//
	//signersAggSigG1 := csquaringmanager.BN254G1Point{
	//	X: new(big.Int).SetBytes(validatedData.SignersAggSigG1[:32]),
	//	Y: new(big.Int).SetBytes(validatedData.SignersAggSigG1[32:]),
	//}
	//
	//nonSignerStakeIndices := make([][]uint32, len(validatedData.NonSignerStakeIndices))
	//for i, indices := range validatedData.NonSignerStakeIndices {
	//	nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	//}
	//
	//signersApkG2 := csquaringmanager.BN254G2Point{
	//	X: [2]*big.Int{
	//		new(big.Int).SetBytes(validatedData.SignersApkG2[:32]),
	//		new(big.Int).SetBytes(validatedData.SignersApkG2[32:64]),
	//	},
	//	Y: [2]*big.Int{
	//		new(big.Int).SetBytes(validatedData.SignersApkG2[64:96]),
	//		new(big.Int).SetBytes(validatedData.SignersApkG2[96:]),
	//	},
	//}
	//
	//nonSignerStakesAndSignature := csquaringmanager.IBLSSignatureVerifierNonSignerStakesAndSignature{
	//	NonSignerPubkeys:            nonSignerPubkeysG1,
	//	GroupApks:                   quorumApksG1,
	//	ApkG2:                       signersApkG2,
	//	Sigma:                       signersAggSigG1,
	//	NonSignerGroupBitmapIndices: validatedData.NonSignerQuorumBitmapIndices,
	//	GroupApkIndices:             validatedData.QuorumApkIndices,
	//	TotalStakeIndices:           validatedData.TotalStakeIndices,
	//	NonSignerStakeIndices:       nonSignerStakeIndices,
	//}
	//
	//s.logger.Debug("RespondToPriceTask",
	//	"task", task, "taskResponse", taskResponse,
	//	"nonSignerStakesAndSignature", nonSignerStakesAndSignature,
	//)
	//err = GatewayClient.RespondToPriceTask(uint64(pkgCtx.ChainID()), task, taskResponse, nonSignerStakesAndSignature)
	//if err != nil {
	//	s.logger.Error("Failed to respond to task", "error", err)
	//	return nil, err
	//}
	//
	//s.logger.Info("ProcessResponseNumberSquared Done")
	//
	//return &types.ResponseNumberSquaredOut{}, nil

	return nil, nil
}
