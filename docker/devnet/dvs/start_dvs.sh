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

  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}

  export AGGREGATOR_RPC_PORT=${AGGREGATOR_RPC_PORT:-26653}
  export AGGREGATOR_RPC_LADDR=${AGGREGATOR_RPC_LADDR:-0.0.0.0:$AGGREGATOR_RPC_PORT}
  export AGGREGATOR_INDEXER_START_HEIGHT=${AGGREGATOR_INDEXER_START_HEIGHT:-299527}
  export AGGREGATOR_INDEXER_BATCH_SIZE=${AGGREGATOR_INDEXER_BATCH_SIZE:-100}
}

function init_aggregator {
  if [ ! -d "$PELLDVS_HOME" ]; then
    pelldvs init --home "$PELLDVS_HOME"
  fi

  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }
  local PELL_DELEGATION_MNAGER=$(fetch_pell_address "delegation_manager_proxy")

  if [ -z "$PELL_DELEGATION_MNAGER" ]; then
    echo "Pell Delegation Manager not found"
    exit 1
  fi

  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    echo "Registry Router Address not found"
    exit 1
  fi

  if [ -z "$ETH_RPC_URL" ]; then
    echo "ETH_RPC_URL not found"
    exit 1
  fi

  update-config rpc_url "$ETH_RPC_URL"
  update-config pell_registry_router_address "$REGISTRY_ROUTER_ADDRESS"
  update-config pell_delegation_manager_address "$PELL_DELEGATION_MNAGER"

  cat <<EOF > $PELLDVS_HOME/config/aggregator.json
{
    "aggregator_rpc_server": "$AGGREGATOR_RPC_LADDR",
    "operator_response_timeout": "10s",
    "pell_registry_router_address": "$REGISTRY_ROUTER_ADDRESS",
    "chain_config_path": "$PELLDVS_HOME/config/chain.detail.json",
    "indexer_start_height": $AGGREGATOR_INDEXER_START_HEIGHT,
    "indexer_batch_size": $AGGREGATOR_INDEXER_BATCH_SIZE,
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