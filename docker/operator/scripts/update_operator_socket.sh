#!/usr/bin/env bash
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}

  source "$(dirname "$0")/utils.sh"
}

function update_socket() {
  local socket_address=$1
  logt "recieved socket address: $socket_address"
  if [ -z "$socket_address" ]; then
    socket_address="http://operator:26657"
  fi

  REGISTRY_ROUTER_ADDRESS=$(ssh emulator "cat /root/RegistryRouterAddress.json" | jq -r .address)
  pelldvs client operator update-socket \
    --home $PELLDVS_HOME \
    --from $OPERATOR_KEY_NAME \
    --rpc-url $ETH_RPC_URL \
    --registry-router "$REGISTRY_ROUTER_ADDRESS" \
    $socket_address
}

load_defaults

logt "Updating operator socket address"
update_socket $1

