package types

import "github.com/0xPellNetwork/pelldvs/aggregator"

func NewValidatedResponse(validatedData *aggregator.ValidatedResponse) *RequestPostRequestValidatedData {
	var nonSignersPubkeysG1 []*G1Point
	for _, pubkey := range validatedData.NonSignersPubkeysG1 {
		x := pubkey.X.Bytes()
		y := pubkey.Y.Bytes()
		nonSignersPubkeysG1 = append(nonSignersPubkeysG1, &G1Point{
			X: x[:],
			Y: y[:],
		})
	}

	var quorumApksG1 []*G1Point
	for _, pubkey := range validatedData.QuorumApksG1 {
		x := pubkey.X.Bytes()
		y := pubkey.Y.Bytes()
		quorumApksG1 = append(quorumApksG1, &G1Point{
			X: x[:],
			Y: y[:],
		})
	}

	var signersApkG2 *G2Point
	if validatedData.SignersApkG2 != nil {
		xReal := validatedData.SignersApkG2.X.A0.Bytes()
		xImag := validatedData.SignersApkG2.X.A1.Bytes()
		yReal := validatedData.SignersApkG2.Y.A0.Bytes()
		yImag := validatedData.SignersApkG2.Y.A1.Bytes()
		signersApkG2 = &G2Point{
			XReal: xReal[:],
			XImag: xImag[:],
			YReal: yReal[:],
			YImag: yImag[:],
		}
	}

	var signersAggSigG1 *G1Point
	if validatedData.SignersAggSigG1 != nil {
		x := validatedData.SignersAggSigG1.X.Bytes()
		y := validatedData.SignersAggSigG1.Y.Bytes()
		signersAggSigG1 = &G1Point{
			X: x[:],
			Y: y[:],
		}
	}

	var nonSignerStakeIndices []*UInt32List
	for _, stakeIndices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, &UInt32List{
			Values: stakeIndices,
		})
	}
	var errMsg string
	if validatedData.Err != nil {
		errMsg = validatedData.Err.Error()
	}

	resp := &RequestPostRequestValidatedData{
		Data:                         validatedData.Data,
		Error:                        errMsg,
		Hash:                         validatedData.Hash,
		NonSignersPubkeysG1:          nonSignersPubkeysG1,
		QuorumApksG1:                 quorumApksG1,
		SignersApkG2:                 signersApkG2,
		SignersAggSigG1:              signersAggSigG1,
		NonSignerQuorumBitmapIndices: validatedData.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             validatedData.QuorumApkIndices,
		TotalStakeIndices:            validatedData.TotalStakeIndices,
		NonSignerStakeIndices:        nonSignerStakeIndices,
	}

	return resp
}
