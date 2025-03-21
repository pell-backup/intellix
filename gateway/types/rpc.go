package types

type RespondToTaskResponse struct {
	Error string `json:"error"`
}

// RPCTaskRaw represents a serializable version of TaskRaw
type RPCTaskRaw struct {
	TaskType                  int64  `json:"task_type"`
	TaskIndex                 uint32 `json:"task_index"`
	RequestID                 []byte `json:"request_id"`
	FeeToken                  string `json:"fee_token"`
	Payment                   string `json:"payment"`
	RequestData               []byte `json:"request_data"`
	CallbackAddress           string `json:"callback_address"`
	CallbackFunctionID        []byte `json:"callback_function_id"`
	TaskCreatedBlock          uint32 `json:"task_created_block"`
	QuorumNumbers             []byte `json:"quorum_numbers"`
	QuorumThresholdPercentage uint32 `json:"quorum_threshold_percentage"`
	AdvanceDecode             bool   `json:"advance_decode"`
}

type RPCTaskResponse struct {
	ReferenceTaskIndex uint32 `json:"reference_task_index"`
	Data               []byte `json:"data"`
}

type RPCValidatedData struct {
	Data                         []byte     `json:"data,omitempty"`
	Error                        string     `json:"error,omitempty"`
	Hash                         []byte     `json:"hash,omitempty"`
	NonSignersPubkeysG1          [][]byte   `json:"non_signers_pubkeys_g_1,omitempty"`
	QuorumApksG1                 [][]byte   `json:"quorum_apks_g_1,omitempty"`
	SignersApkG2                 []byte     `json:"signers_apk_g_2,omitempty"`
	SignersAggSigG1              []byte     `json:"signers_agg_sig_g_1,omitempty"`
	NonSignerQuorumBitmapIndices []uint32   `json:"non_signer_quorum_bitmap_indices,omitempty"`
	QuorumApkIndices             []uint32   `json:"quorum_apk_indices,omitempty"`
	TotalStakeIndices            []uint32   `json:"total_stake_indices,omitempty"`
	NonSignerStakeIndices        [][]uint32 `json:"non_signer_stake_indices,omitempty"`
}

type RPCVoteFinalizedRequestIn struct {
	ChainID        int64             `json:"chain_id"`
	TaskRaw        *RPCTaskRaw       `json:"task_raw"`
	ValidatedData  *RPCValidatedData `json:"validated_data"`
	RespToTaskData []byte            `json:"resp_to_task_data"` // decoded data
}
