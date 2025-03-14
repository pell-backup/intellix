#!/bin/bash

set -x
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
  export DEBUG_ENABLED=${DEBUG_ENABLED:-false}
  export DEBUG_PORT=${DEBUG_PORT:-2345}

  export GATEWAY_PORT=${GATEWAY_PORT:-8949}
  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}

  export CHAIN_ID=${CHAIN_ID:-1337}
  ## TODO: remove this after the integration with the operator
  export OPERATOR_RPC_SERVER=${OPERATOR_RPC_SERVER:-operator:26657}

  export AGGREGATOR_INDEXER_START_HEIGHT=${AGGREGATOR_INDEXER_START_HEIGHT:-0}
  export AGGREGATOR_INDEXER_BATCH_SIZE=${AGGREGATOR_INDEXER_BATCH_SIZE:-1000}

  export SERVICE_CHAIN_ID=${SERVICE_CHAIN_ID:-1337}
  export SERVICE_CHAIN_RPC_URL=${SERVICE_CHAIN_RPC_URL:-http://eth:8545}
  export SERVICE_CHAIN_WS_URL=${SERVICE_CHAIN_WS_URL:-ws://eth:8545}
}

function setup_gateway_key {
  export GATEWAY_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
  if ! pelldvs keys show gateway --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure gateway $GATEWAY_KEY --home $PELLDVS_HOME >/dev/null
  fi

  export GATEWAY_ADDRESS=$(pelldvs keys show gateway --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
}

function contract_healthcheck {
  set +e
  while true; do
    DATA_ORACLE_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/DataOracle-Proxy.json" | jq -r .address)
    if [ -z "$DATA_ORACLE_ADDRESS" ]; then
      echo "DATA_ORACLE_ADDRESS is empty, waiting..."
      sleep 5
    else
      echo "contract_address is ready"
      break
    fi
  done
  set -e
}

function dvs_healthcheck {
  set +e
  while true; do
    curl -s $AGGREGATOR_RPC_SERVER >/dev/null
    if [ $? -eq 52 ]; then
      echo "DVS RPC port is ready, proceeding to the next step..."
      break
    fi
    echo "DVS RPC port not ready, retrying in 2 seconds..."
    sleep 2
  done
  ## Wait for aggregator to be ready
  sleep 3
  set -e
}

function init_pelldvs_config {
  pelldvs init --home $PELLDVS_HOME

  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }

  ## update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellRegistryRouterFactory.json" | jq -r .address)
  PELL_DELEGATION_MNAGER=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellDelegationManager-Proxy.json" | jq -r .address)
  PELL_DVS_DIRECTORY=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellDVSDirectory-Proxy.json" | jq -r .address)
  REGISTRY_ROUTER_ADDRESS=$(ssh emulator "cat /root/RegistryRouterAddress.json" | jq -r .address)

  update-config aggregator_rpc_url "$AGGREGATOR_RPC_URL"
	update-config interactor_config_path "$PELLDVS_HOME/config/interactor_config.json"

  ## FIXME: don't use absolute path for key
  update-config operator_bls_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.bls.key.json"
  update-config operator_ecdsa_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.ecdsa.key.json"

  DVS_OPERATOR_KEY_MANAGER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorKeyManager-Proxy.json" | jq -r .address)
  DVS_CENTRAL_SCHEDULER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/CentralScheduler-Proxy.json" | jq -r .address)
  DVS_OPERATOR_INFO_PROVIDER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorInfoProvider.json" | jq -r .address)
  DVS_OPERATOR_INDEX_MANAGER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorIndexManager-Proxy.json" | jq -r .address)

  cat <<EOF > $PELLDVS_HOME/config/interactor_config.json
{
    "rpc_url": "$ETH_RPC_URL",
    "chain_id": $CHAIN_ID,
    "contract_config": {
      "indexer_start_height": $AGGREGATOR_INDEXER_START_HEIGHT,
      "indexer_batch_size": $AGGREGATOR_INDEXER_BATCH_SIZE,
      "pell_registry_router_factory": "$REGISTRY_ROUTER_FACTORY_ADDRESS",
      "pell_dvs_directory": "$PELL_DVS_DIRECTORY",
      "pell_delegation_manager": "$PELL_DELEGATION_MNAGER",
      "pell_registry_router": "$REGISTRY_ROUTER_ADDRESS",
      "dvs_configs": {
        "$SERVICE_CHAIN_ID": {
          "chain_id": $SERVICE_CHAIN_ID,
          "rpc_url": "$SERVICE_CHAIN_RPC_URL",
          "ws_url": "$SERVICE_CHAIN_WS_URL",
          "operator_info_provider": "$DVS_OPERATOR_INFO_PROVIDER",
          "operator_key_manager": "$DVS_OPERATOR_KEY_MANAGER",
          "central_scheduler": "$DVS_CENTRAL_SCHEDULER",
          "operator_index_manager": "$DVS_OPERATOR_INDEX_MANAGER"
        }
      }
    }
}
EOF

cat $PELLDVS_HOME/config/interactor_config.json

tail -n 10 $PELLDVS_HOME/config/config.toml

}

function setup_gateway_config {
  setup_gateway_key
  DATA_ORACLE_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/DataOracle-Proxy.json" | jq -r .address)

  # TODO: remove sender_address from config
  cat <<EOF > $PELLDVS_HOME/config/gateway.config.json
{
  "server_addr": "0.0.0.0:$GATEWAY_PORT",
  "sender_address": "$GATEWAY_ADDRESS",
  "private_key_store_path": "$PELLDVS_HOME/keys/gateway.ecdsa.key.json",
  "chains": {
    "1337": {
      "eth_endpoint": "$ETH_WS_URL",
      "contract_address": "$DATA_ORACLE_ADDRESS",
      "chain_id": 1337,
      "gas_limit": 1000000,
      "chain_id": $SERVICE_CHAIN_ID,
      "rpc_url": "$SERVICE_CHAIN_RPC_URL",
      "ws_url": "$SERVICE_CHAIN_WS_URL",
      "contract_address": "$DATA_ORACLE_ADDRESS"
    }
  }
}
EOF

cat $PELLDVS_HOME/config/gateway.config.json

}

function start_gateway {
  if [ "$DEBUG_ENABLED" = "true" ]; then
    dlv exec /usr/bin/intellixd \
      --listen=:$DEBUG_PORT --headless=true --api-version=2 --accept-multiclient\
      -- start-task-gateway
  else
    intellixd start-task-gateway --home $PELLDVS_HOME
  fi
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

dvs_healthcheck

logt "setup pelldvs config"
init_pelldvs_config

logt "Setup gateway config"
setup_gateway_config

touch /root/gateway_initialized

logt "Starting gateway..."
start_gateway
