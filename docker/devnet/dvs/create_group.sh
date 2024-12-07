
set -e
set -x

function load_defaults {
  export GITHUB_TOKEN=${GITHUB_TOKEN}
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
}

fetch_contract_address() {
  curl -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3.raw" \
      https://api.github.com/repos/0xPellNetwork/contracts/contents/$1 | jq -r '.address'
}

function update_pelldvs_config {
  pelldvs init --home "$PELLDVS_HOME"

  # TODO: should get address from contract
  REGISTRY_ROUTER_ADDRESS=$(cat $PELLDVS_HOME/RegistryRouterAddress.json | jq -r .address)
  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" ~/.pelldvs/config/config.toml
  }
  update-config rpc_url "$ETH_RPC_URL"
  update-config registry_router_address "$REGISTRY_ROUTER_ADDRESS"
}

function create_group {
  STBTC_STRATEGY_ADDRESS=$(fetch_contract_address "deployments/$NETWORK/stBTC-Strategy-Proxy.json")
  BTCB_STRATEGY_ADDRESS=$(fetch_contract_address "deployments/$NETWORK/BTCB-Strategy-Proxy.json")
  cat > $PELLDVS_HOME/group-0-config.json <<EOF
{
  "minimum_stake": 0,
  "pool_params": [
    {
      "chain_id": $CHAIN_ID,
      "multiplier": 1,
      "pool": "$STBTC_STRATEGY_ADDRESS"
    },
    {
      "chain_id": $CHAIN_ID,
      "multiplier": 1,
      "pool": "$BTCB_STRATEGY_ADDRESS"
    }
  ],
  "operator_set_params": {
    "kick_bi_ps_of_operator_stake": 10,
    "kick_bi_ps_of_total_stake": 10,
    "max_operator_count": 1000
  }
}
EOF

  pelldvs client dvs create-group \
    --home $PELLDVS_HOME \
    --from admin \
    --config $PELLDVS_HOME/group-0-config.json
}

function show_group {
  REGISTRY_ROUTER_ADDRESS=$(cat $PELLDVS_HOME/RegistryRouterAddress.json | jq -r .address)
  GROUP_COUNT=$(cast call "$REGISTRY_ROUTER_ADDRESS" "groupCount()" --rpc-url "$ETH_RPC_URL")
  echo "Group Count From Registry Router in Pell EVM: $GROUP_COUNT"
}

load_defaults

if [ "$1" == "create" ]; then
  update_pelldvs_config
  create_group
else
  show_group
fi
