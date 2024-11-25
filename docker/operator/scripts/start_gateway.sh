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

  export GATEWAY_PORT=${GATEWAY_PORT:-8949}
  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}
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

function setup_gateway_key {
  export GATEWAY_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
  if ! pelldvs keys show gateway --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure gateway $GATEWAY_KEY --home $PELLDVS_HOME >/dev/null
  fi

  export GATEWAY_ADDRESS=$(pelldvs keys show gateway --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
}


function setup_gateway_config {
  setup_gateway_key
  DATA_ORACLE_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/DataOracle-Proxy.json" | jq -r .address)

  # TODO: remove sender_address from config
  cat <<EOF > $PELLDVS_HOME/config/gateway.config.json
{
  "server_addr": "0.0.0.0:$GATEWAY_PORT",
  "sender_address": "$GATEWAY_ADDRESS",
  "eth_endpoint": "$ETH_WS_URL",
  "contract_address": "$DATA_ORACLE_ADDRESS",
  "private_key_store_path": "$PELLDVS_HOME/keys/gateway.ecdsa.key.json"
}
EOF
}

function start_gateway {
  ## TODO: add home dir flag
  PELLDVS_HOME=$PELLDVS_HOME intellixd start-task-gateway
}

function start_gateway_debug {
  ## TODO: add home dir flag
  go install github.com/go-delve/delve/cmd/dlv@latest
  dlv exec /usr/bin/intellixd \
    --listen=:2345 --headless=true --api-version=2 --accept-multiclient\
    -- start-task-gateway
  # PELLDVS_HOME=$PELLDVS_HOME intellixd start-task-gateway
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Check if DVS is ready"
dvs_healthcheck

logt "Setup gateway config"
setup_gateway_config

touch /root/gateway_initialized

logt "Starting gateway..."
start_gateway
