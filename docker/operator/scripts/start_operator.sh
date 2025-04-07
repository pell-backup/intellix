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
  export INTELLIX_HOME=${INTELLIX_HOME:-/root/.intellix}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  export GATEWAY_ADDR=${GATEWAY_ADDR:-gateway:8949}
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
  export DEBUG_ENABLED=${DEBUG_ENABLED:-false}
  export DEBUG_PORT=${DEBUG_PORT:-2345}

  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}
  export COSMOS_KEYRING_BACKEND=${COSMOS_KEYRING_BACKEND:-test}
  export COSMOS_CHAIN_ID=${COSMOS_CHAIN_ID:-intellix}
  export COSMOS_NODE_URI=${COSMOS_NODE_URI:-http://abci:26657}
  export OPERATOR_RPC_SERVER=${OPERATOR_RPC_SERVER:-operator:26657}

  export COINMARKETCAP_API_KEY=${COINMARKETCAP_API_KEY:-""}
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

function gateway_healthcheck {
  set +e
  while true; do
    curl -s $GATEWAY_ADDR >/dev/null
    if [ $? -eq 52 ]; then
      echo "Gateway is ready, proceeding to the next step..."
      break
    fi
    echo "Gateway not ready, retrying in 2 seconds..."
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

function gen_cosmos_key {
  # TODO: generate new key and use admin to faceut
  mkdir -p "$PELLDVS_HOME/keyring-test/"
  scp abci:/root/.intellix/keyring-test/* "$PELLDVS_HOME/keyring-test/"
}

function setup_operator_config {
  ## migrate to dvs logic after fix
  export OPERATOR_ADDRESS=$(pelldvs keys show $OPERATOR_KEY_NAME --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
  ## TODO: use operator key on config.toml and gateway should be on app.toml
  cat <<EOF > $PELLDVS_HOME/config/operator.config.json
{
  "operator_address": "$OPERATOR_ADDRESS",
  "gateway_addr": "$GATEWAY_ADDR",
  "cosmos_node_uri": "$COSMOS_NODE_URI",
  "cosmos_chain_id": "$COSMOS_CHAIN_ID",
  "price_tick_converter_config": {
    "binance": {
      "USD": "USDT"
    },
    "coinbase": {
      "USD": "USD"
    },
    "okx": {
      "USD": "USDT"
    },
    "gate": {
      "USD": "USDT"
    },
    "coinmarketcap": {
      "USD": "USD"
    }
  },
  "api_key": {
    "coinmarketcap": "$COINMARKETCAP_API_KEY"
  },
  "symbol_list": {
    "crypto": {
      "TON": "USD",
      "BTC": "USD",
      "SAVM": "USD",
      "DOGE": "USD",
      "ORDI": "USD",
      "PEPE": "USD",
      "CAT": "USD",
      "GOAT": "USD",
      "NEIROCTO": "USD",
      "XLM": "USD",
      "XRP": "USD",
      "ACT": "USD",
      "DOGS": "USD",
      "MEME": "USD",
      "BNB": "USD",
      "ETH": "USD",
      "MUBI": "USD",
      "BB": "USD",
      "CATS": "USD",
      "DAI": "USD",
      "LADYS": "USD",
      "SUNDOG": "USD",
      "BAN": "USD",
      "PUSS": "USD",
      "MOODENG": "USD",
      "USDC": "USD",
      "NEIRO": "USD",
      "1000SATS": "USD",
      "NEIROETH": "USD",
      "AUCTION": "USD",
      "WIF": "USD",
      "USDT": "USD",
      "BOME": "USD"
    },
    "a_shares": {
      "1_000001": "USD",
      "0_399001": "USD"
    },
    "hk_shares": {
      "100_HSI": "USD"
    },
    "us_stocks": {
      "MSTR": "USD",
      "BB": "USD",
      "COIN": "USD",
      "NVDA": "USD",
      "TSLA": "USD",
      "PLTR": "USD",
      "AAPL": "USD",
      "AMZN": "USD",
      "GOOGL": "USD"
    }
  }
}
EOF

  logt "Operator config created:"
  cat $PELLDVS_HOME/config/operator.config.json
}

function start_operator {
  if [ "$DEBUG_ENABLED" = "true" ]; then
    dlv exec /usr/bin/intellixd \
      --listen=:$DEBUG_PORT --headless=true --api-version=2 --accept-multiclient\
      -- start-operator
  else
    intellixd start-operator
  fi
}

function upload_wasm_script {
  ssh abci "intellixd tx processor create-processor 'Intellix' ./scripts/processor_data/mock_processor.wasm --from $OPERATOR_KEY_NAME --chain-id $COSMOS_CHAIN_ID --keyring-backend test --gas auto --fees 2000000uixn -y"
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Check if DVS is ready"
dvs_healthcheck

logt "Check if Gateway is ready"
#gateway_healthcheck

if [ ! -f /root/operator_initialized ]; then
  logt "Init operator"
  source "$(dirname "$0")/init_operator.sh"
  gen_cosmos_key
  upload_wasm_script

  touch /root/operator_initialized
fi

logt "Setup operator config"
setup_operator_config

logt "Setup dispatcher config"
setup_dispatcher_config

touch /root/dispatcher_initialized

logt "Starting operator..."
start_operator
