set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

source "$(dirname "$0")/utils.sh"

function load_defaults {
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
  export AGGREGATOR_RPC_URL=${AGGREGATOR_RPC_URL:-dvs:26653}
  export OPERATOR_NODE_NAME=${OPERATOR_NODE_NAME:-operator}
  export OPERATOR_PUBLIC_RPC_URL=${OPERATOR_PUBLIC_RPC_URL}
  export GROUPS_TO_ADD=${GROUPS_TO_ADD:-0}

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}
}

## TODO: move operator config to seperated location
function init_pelldvs_config {
  pelldvs init --home $PELLDVS_HOME
  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }

  ## update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(fetch_pell_address "registry_router_factory")
  PELL_DELEGATION_MNAGER=$(fetch_pell_address "delegation_manager_proxy")
  PELL_DVS_DIRECTORY=$(fetch_pell_address "dvs_directory_proxy")

  update-config rpc_url "$ETH_RPC_URL"
  update-config pell_registry_router_factory_address "$REGISTRY_ROUTER_FACTORY_ADDRESS"
  update-config pell_delegation_manager_address "$PELL_DELEGATION_MNAGER"
  update-config pell_dvs_directory_address "$PELL_DVS_DIRECTORY"
  update-config pell_registry_router_address "$REGISTRY_ROUTER_ADDRESS"
  update-config aggregator_rpc_url "$AGGREGATOR_RPC_URL"

  ## FIXME: don't use absolute path for key
  update-config operator_bls_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.bls.key.json"
  update-config operator_ecdsa_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.ecdsa.key.json"
}

function register_operator_to_dvs {
  pelldvs client operator register-operator-to-dvs \
    --home $PELLDVS_HOME \
    --from $OPERATOR_KEY_NAME \
    --groups ${GROUPS_TO_ADD} \
    --socket ${OPERATOR_PUBLIC_RPC_URL}
  show_dvs_operator_info $OPERATOR_ADDRESS
}


load_defaults

logt "init pelldvs config"
init_pelldvs_config

logt "setup operator key"
source "$(dirname "$0")/setup_operator_key.sh"

if [ "$1" == "register" ]; then 
  register_operator_to_dvs
else
  show_dvs_operator_info "$OPERATOR_ADDRESS"
fi
