#!/bin/bash

set -x

source "$(dirname "$0")/utils.sh"

function load_defaults {
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export ETH_RPC_URL=${SERVICE_CHAIN_RPC_URL:-https://bsc-testnet.blockpi.network/v1/rpc/public}
  export HARDHAT_DVS_PATH="deployments/$NETWORK"
}

load_defaults
PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS=$(fetch_dvs_address "$HARDHAT_DVS_PATH/PriceOraclePayInNativeConsumer.json")

RESULT=$(cast call "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "price()" | cast to-dec)
cast send "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "requestPrice(string)" "ETH" --private-key "$ADMIN_KEY"