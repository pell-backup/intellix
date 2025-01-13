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
      "gas_limit": 100000
    }
  }
}
EOF
}

function start_gateway {
  if [ "$DEBUG_ENABLED" = "true" ]; then
    dlv exec /usr/bin/intellixd \
      --listen=:$DEBUG_PORT --headless=true --api-version=2 --accept-multiclient\
      -- start-task-gateway
  else
    intellixd start-task-gateway
  fi
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

contract_healthcheck

logt "Setup gateway config"
setup_gateway_config

touch /root/gateway_initialized

logt "Starting gateway..."
start_gateway
