#!/bin/bash

set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/price-oracle-dvs/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/price-oracle-dvs/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}

  export AGGREGATOR_RPC_LADDR=${AGGREGATOR_RPC_LADDR:-0.0.0.0:26653}
}

function init_aggregator {
  mkdir -p $PELLDVS_HOME/config
  cat <<EOF > $PELLDVS_HOME/config/aggregator.json
{
    "aggregator_rpc_server": "$AGGREGATOR_RPC_LADDR",
    "operator_response_timeout": "10s",
    "chain_config_path": "$PELLDVS_HOME/config/chain.detail.json"
}
EOF

  DVS_BLS_APK_REGISTRY=$(ssh hardhat "cat $HARDHAT_DVS_PATH/BLSApkRegistry-Proxy.json" | jq -r .address)
  DVS_REGISTRY_COORDINATOR=$(ssh hardhat "cat $HARDHAT_DVS_PATH/RegistryCoordinator-Proxy.json" | jq -r .address)
  DVS_OPERATOR_STATE_RETRIEVER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorStateRetriever.json" | jq -r .address)
  cat <<EOF > $PELLDVS_HOME/config/chain.detail.json
{
  "1337": {
    "rpc_url": "$ETH_WS_URL",
    "operator_state_retriever_address": "$DVS_OPERATOR_STATE_RETRIEVER",
    "bls_apk_registry_address": "$DVS_BLS_APK_REGISTRY",
    "registry_coordinator_address": "$DVS_REGISTRY_COORDINATOR"
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