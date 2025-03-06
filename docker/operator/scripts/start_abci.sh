#!/bin/bash

set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/price-oracle-dvs/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/price-oracle-dvs/deployments/localhost"

  export INTELLIX_HOME=${INTELLIX_HOME:-/root/.intellix}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}

  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}
  export COSMOS_NODE_NAME=${COSMOS_NODE_NAME:-node0}
  export COSMOS_CHAIN_ID=${COSMOS_CHAIN_ID:-intellix}
  export COSMOS_KEYRING_BACKEND=${COSMOS_KEYRING_BACKEND:-test}

  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
}

function init_genesis {
  intellixd init $COSMOS_NODE_NAME --chain-id $COSMOS_CHAIN_ID
  sed -i 's/stake/uitlx/g' ~/.intellix/config/genesis.json

  export DEFAULT_KEY=${DEFAULT_KEY:-mykey}

  echo "Generating genesis.json"
  intellixd keys add $DEFAULT_KEY --keyring-backend test
  intellixd keys add $OPERATOR_KEY_NAME --keyring-backend test
  ACCOUNT_ADDRESS=$(intellixd keys show $DEFAULT_KEY -a --keyring-backend test)
  OPERATOR_ADDRESS=$(intellixd keys show $OPERATOR_KEY_NAME -a --keyring-backend test)
  
  intellixd genesis add-genesis-account $ACCOUNT_ADDRESS 20000000000uitlx
  intellixd genesis add-genesis-account $OPERATOR_ADDRESS 10000000000uitlx
  intellixd genesis gentx $DEFAULT_KEY 1000000uitlx \
    --chain-id $COSMOS_CHAIN_ID \
    --moniker $COSMOS_NODE_NAME \
    --keyring-backend test \
    --from $ACCOUNT_ADDRESS
  intellixd genesis collect-gentxs
}

function init_config {
  # max_body_bytes
  dasel put -f $INTELLIX_HOME/config/config.toml -v 1000000 'rpc.max_body_bytes' >> $INTELLIX_HOME/config/config.toml
  dasel put -f $INTELLIX_HOME/config/config.toml -v 'tcp://0.0.0.0:26657' 'rpc.laddr' >> $INTELLIX_HOME/config/config.toml
  # api max body bytes
  dasel put -f $INTELLIX_HOME/config/app.toml -v 1000000 'api.rpc-max-body-bytes' >> $INTELLIX_HOME/config/app.toml
}

function start_abci {
  intellixd start \
    --minimum-gas-prices=1uitlx \
    --api.enable=true \
    --api.address="tcp://0.0.0.0:1317" \
    --grpc.enable=true \
    --grpc.address="0.0.0.0:9090"
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

if [ ! -f /root/genesis_initialized ]; then

  logt "Init Genesis"
  init_genesis

  touch /root/genesis_initialized

fi

logt "Init Config"
init_config

logt "Starting ABCI in background..."
start_abci