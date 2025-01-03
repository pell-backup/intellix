#!/bin/bash

set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export INTELLIX_HOME=${INTELLIX_HOME:-/root/.intellix}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}

  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}
  export COSMOS_NODE_NAME=${COSMOS_NODE_NAME:-node0}
  export COSMOS_CHAIN_ID=${COSMOS_CHAIN_ID:-intellix}
  export COSMOS_KEYRING_BACKEND=${COSMOS_KEYRING_BACKEND:-test}

  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
  export VALIDATOR_KEY_NAME=${VALIDATOR_KEY_NAME:-validator}
  export VALIDATOR_KEY_MNEMONIC=${VALIDATOR_KEY_MNEMONIC} 
  export OPERATOR_KEY_MNEMONIC=${OPERATOR_KEY_MNEMONIC}
}


function init_genesis {
  intellixd init $COSMOS_NODE_NAME --chain-id $COSMOS_CHAIN_ID


  echo "Generating genesis.json"
  if intellixd keys show $VALIDATOR_KEY_NAME --keyring-backend test; then
    echo "Validator key already exists"
  else
    echo $VALIDATOR_KEY_MNEMONIC | intellixd keys add $VALIDATOR_KEY_NAME --recover --keyring-backend=test
  fi

  if intellixd keys show $OPERATOR_KEY_NAME --keyring-backend test; then
    echo "Operator key already exists"
  else
    echo $OPERATOR_KEY_MNEMONIC | intellixd keys add $OPERATOR_KEY_NAME --recover --keyring-backend=test
  fi

  ACCOUNT_ADDRESS=$(intellixd keys show $VALIDATOR_KEY_NAME -a --keyring-backend test)
  OPERATOR_ADDRESS=$(intellixd keys show $OPERATOR_KEY_NAME -a --keyring-backend test)
  
  intellixd genesis add-genesis-account $ACCOUNT_ADDRESS 20000000000stake
  intellixd genesis add-genesis-account $OPERATOR_ADDRESS 10000000000stake
  intellixd genesis gentx $VALIDATOR_KEY_NAME 1000000stake \
    --chain-id $COSMOS_CHAIN_ID \
    --moniker $COSMOS_NODE_NAME \
    --keyring-backend test \
    --from $ACCOUNT_ADDRESS
  intellixd genesis collect-gentxs
}

function init_config {
  dasel put -f $INTELLIX_HOME/config/config.toml -v 'tcp://0.0.0.0:26657' 'rpc.laddr' > $INTELLIX_HOME/config/config.toml
}

function start_abci {
  intellixd start \
    --minimum-gas-prices=0.01stake \
    --api.enable=true \
    --api.address="tcp://0.0.0.0:1317" \
    --grpc.enable=true \
    --grpc.address="0.0.0.0:9090"
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

if [ ! -f ${INTELLIX_HOME}/config/genesis.json ]; then

  logt "Init Genesis"
  init_genesis

  touch /root/genesis_initialized

fi

logt "Init Config"
init_config

logt "Starting ABCI in background..."
start_abci