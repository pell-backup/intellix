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

  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}
  export COSMOS_NODE_NAME=${COSMOS_NODE_NAME:-node0}
  export COSMOS_CHAIN_ID=${COSMOS_CHAIN_ID:-intellix}
  export COSMOS_KEYRING_BACKEND=${COSMOS_KEYRING_BACKEND:-test}
}


function init_genesis {
  intellixd init $COSMOS_NODE_NAME --chain-id $COSMOS_CHAIN_ID

  export DEFAULT_KEY=${DEFAULT_KEY:-mykey}

  echo "Generating genesis.json"
  intellixd keys add $DEFAULT_KEY --keyring-backend test
  ACCOUNT_ADDRESS=$(intellixd keys show $DEFAULT_KEY -a --keyring-backend test)
  
  intellixd genesis add-genesis-account $ACCOUNT_ADDRESS 10000000000stake
  intellixd genesis gentx $DEFAULT_KEY 10000000000stake \
    --chain-id $COSMOS_CHAIN_ID \
    --moniker $COSMOS_NODE_NAME \
    --keyring-backend test \
    --from $ACCOUNT_ADDRESS
  intellixd genesis collect-gentxs
}

function init_validator {
  echo "Initializing validator"
  export DEFAULT_VALIDATOR_KEY=${DEFAULT_VALIDATOR_KEY:-validator}
  intellixd keys add $DEFAULT_VALIDATOR_KEY --keyring-backend test

  export DEFAULT_VALIDATOR_ADDR=$(intellixd keys show $DEFAULT_VALIDATOR_KEY -a --keyring-backend test)

  cat <<EOF > $HOME/$DEFAULT_VALIDATOR_KEY.json
{
        "pubkey": $(intellixd tendermint show-validator),
        "amount": "1000000stake",
        "moniker": "$COSMOS_NODE_NAME",
        "identity": "",
        "commission-rate": "0.1",
        "commission-max-rate": "0.2",
        "commission-max-change-rate": "0.01",
        "min-self-delegation": "1"
}
EOF

  intellixd tx staking create-validator $HOME/$DEFAULT_VALIDATOR_KEY.json \
    --from $DEFAULT_KEY \
    --chain-id $COSMOS_CHAIN_ID \
    --keyring-backend test -y

  intellixd tx bank send $ACCOUNT_ADDRESS $DEFAULT_VALIDATOR_ADDR 10000000stake \
    --chain-id $COSMOS_CHAIN_ID \
    --keyring-backend test -y
}

function start_abci {
  # intellixd --home "$PELLDVS_HOME"
  PELLDVS_HOME=$PELLDVS_HOME intellixd start \
    --minimum-gas-prices=0.001stake \
    --api.enable=true \
    --api.address="tcp://0.0.0.0:1317" \
    --grpc.enable=true \
    --grpc.address="0.0.0.0:9090"
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Init Genesis"
init_genesis

logt "Starting ABCI..."
start_abci

logt "Init Validator"
init_validator