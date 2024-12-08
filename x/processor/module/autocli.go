package processor

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "intellix/api/intellix/processor"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "ShowProcessor",
					Use:            "show-processor [id]",
					Short:          "Shows a processor by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},

				{
					RpcMethod:      "ListProcessor",
					Use:            "list-processor",
					Short:          "Lists all processors",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateProcessor",
					Use:            "create-processor [config] [wasm-code]",
					Short:          "Creates a new processor",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "config"}, {ProtoField: "wasmCode"}},
				},
				{
					RpcMethod: "VoteRequestProcessor",
				},
				{
					RpcMethod: "VoteResponseProcessor",
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
