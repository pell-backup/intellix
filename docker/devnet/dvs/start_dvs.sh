#!/bin/bash

set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

source "$(dirname "$0")/utils.sh"

function load_defaults {
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export HARDHAT_DVS_PATH="deployments/$NETWORK"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export SERVICE_CHAIN_RPC_URL=${SERVICE_CHAIN_RPC_URL:-https://bsc-testnet.blockpi.network/v1/rpc/public}
  export SERVICE_CHAIN_WS_URL=${SERVICE_CHAIN_WS_URL}

  export AGGREGATOR_RPC_LADDR=${AGGREGATOR_RPC_LADDR:-0.0.0.0:26653}
}

function init_aggregator {
  mkdir -p $PELLDVS_HOME/config
  REGISTRY_ROUTER_ADDRESS=$(cat $PELLDVS_HOME/RegistryRouterAddress.json | jq -r .address)

  cat <<EOF > $PELLDVS_HOME/config/aggregator.json
{
    "aggregator_rpc_server": "$AGGREGATOR_RPC_LADDR",
    "operator_response_timeout": "10s",
    "pell_registry_router_address": "$REGISTRY_ROUTER_ADDRESS",
    "chain_config_path": "$PELLDVS_HOME/config/chain.detail.json",
    "indexer_batch_size": 1000,
    "indexer_listen_interval": 5000000000
}
EOF

  DVS_OPERATOR_KEY_MANAGER=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorKeyManager-Proxy.json")
  DVS_CENTRAL_SCHEDULER=$(fetch_dvs_address "$HARDHAT_DVS_PATH/CentralScheduler-Proxy.json")
  DVS_OPERATOR_INFO_PROVIDER=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorInfoProvider.json")
  DVS_OPERATOR_INDEX_MANAGER=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorIndexManager-Proxy.json")
  cat <<EOF > $PELLDVS_HOME/config/chain.detail.json
{
  "$CHAIN_ID": {
    "rpc_url": "$SERVICE_CHAIN_RPC_URL",
    "operator_info_provider_address": "$DVS_OPERATOR_INFO_PROVIDER",
    "operator_key_manager_address": "$DVS_OPERATOR_KEY_MANAGER",
    "central_scheduler_address": "$DVS_CENTRAL_SCHEDULER",
    "operator_index_manager_address": "$DVS_OPERATOR_INDEX_MANAGER"
  }
}
EOF
}

function start_aggregator {
  pelldvs start-aggregator --home "$PELLDVS_HOME"
}

logt "Load Default Values for ENV Vars if not set."
load_defaults

if [ ! -f /root/aggregator_initialized ]; then
  logt "Init aggregator"
  init_aggregator
  touch /root/aggregator_initialized
fi

logt "Starting aggregator..."
start_aggregator