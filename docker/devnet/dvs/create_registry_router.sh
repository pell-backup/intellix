set -x
set -e

function load_defaults {
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
}

fetch_pell_address() {
  KEY=$1
  curl https://raw.githubusercontent.com/0xPellNetwork/network-config/refs/heads/main/testnet/system_contract.json | jq -r ".$KEY"
}

function update_pelldvs_config {
  pelldvs init --home "$PELLDVS_HOME"

  ## Update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(fetch_pell_address "registry_router_factory")

  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" ~/.pelldvs/config/config.toml
  }
  update-config rpc_url "$ETH_RPC_URL"
  update-config registry_router_factory_address "$REGISTRY_ROUTER_FACTORY_ADDRESS"
}

function create_registry_router {
  ## Create registry router
  export ADMIN_ADDRESS=$(pelldvs keys show admin --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
  REGISTRY_ROUTER_ADDRESS_FILE="$PELLDVS_HOME/RegistryRouterAddress.json"
  pelldvs client dvs create-registry-router \
    --home $PELLDVS_HOME \
    --from admin \
    --initial-owner $ADMIN_ADDRESS \
    --dvs-chain-approver $ADMIN_ADDRESS \
    --churn-approver $ADMIN_ADDRESS \
    --ejector $ADMIN_ADDRESS \
    --pauser $ADMIN_ADDRESS \
    --unpauser $ADMIN_ADDRESS \
    --initial-paused-status false \
    --save-to-file $REGISTRY_ROUTER_ADDRESS_FILE \
    --force-save true
}

load_defaults
update_pelldvs_config
create_registry_router
