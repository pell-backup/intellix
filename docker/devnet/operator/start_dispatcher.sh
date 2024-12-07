#!/bin/bash

set -x
set -e

source "$(dirname "$0")/utils.sh"

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export HARDHAT_DVS_PATH="deployments/$NETWORK"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  ## TODO: remove this after the integration with the operator
  export OPERATOR_RPC_SERVER=${OPERATOR_RPC_SERVER:-operator:26657}

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


function setup_dispatcher_config {
  mkdir -p $PELLDVS_HOME/config
  DATA_ORACLE_ADDRESS=$(fetch_dvs_address "$HARDHAT_DVS_PATH/DataOracle-Proxy.json")
  cat <<EOF > $PELLDVS_HOME/config/dispatcher.config.json
{
  "dvs_address": "tcp://$OPERATOR_RPC_SERVER",
  "chains": [
    {
      "chain_id": $CHAIN_ID,
      "eth_url": "$ETH_WS_URL",
      "contract_address": "$DATA_ORACLE_ADDRESS"
    }
  ]
}
EOF
}

function start_dispatcher {
  if [ "$DEBUG_ENABLED" = "true" ]; then
    dlv exec /usr/bin/intellixd \
      --listen=:$DEBUG_PORT --headless=true --api-version=2 --accept-multiclient\
      --log --log-output=debugger \
      -- start-task-dispatcher
  else
    intellixd start-task-dispatcher
  fi
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Check if DVS is ready"
dvs_healthcheck

logt "Setup dispatcher config"
setup_dispatcher_config

touch /root/dispatcher_initialized

logt "Starting dispatcher..."
start_dispatcher
