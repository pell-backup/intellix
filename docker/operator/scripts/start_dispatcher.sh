#!/bin/bash

set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/price-oracle-dvs/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/price-oracle-dvs/deployments/localhost"
  export DEBUG_ENABLED=${DEBUG_ENABLED:-false}
  export DEBUG_PORT=${DEBUG_PORT:-2345}

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
  DATA_ORACLE_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/DataOracle-Proxy.json" | jq -r .address)
  cat <<EOF > $PELLDVS_HOME/config/dispatcher.config.json
{
  "dvs_address": "tcp://$OPERATOR_RPC_SERVER",
  "chains": [
    {
      "chain_id": 1337,
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
