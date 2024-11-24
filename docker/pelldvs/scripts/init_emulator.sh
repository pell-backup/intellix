#!/bin/bash

set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function hardhat_healthcheck {
  set +e
  while true; do
    ssh hardhat "test -f /root/contracts_deployed_completed"
    if [ $? -eq 0 ]; then
      echo "Contracts deployed, proceeding to the next step..."
      break
    fi
    echo "Contracts not deployed, retrying in 1 second..."
    sleep 1
  done
  set -e
}


function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/price-oracle-dvs/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/price-oracle-dvs/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  export ADMIN_KEY_FILE="$PELLDVS_HOME/keys/admin.ecdsa.key.json"
}


function setup_admin_key {
  ## create admin key
  # echo  -ne '\n\n' | pelldvs keys create admin --key-type=ecdsa --insecure > /tmp/admin.key
  # ADMIN_KEY=$(cat /tmp/admin.key | sed -n 's/.*\/\/[[:space:]]*\([0-9a-f]\{64\}\)[[:space:]]*\/\/.*/\1/p')

  ## For development purposes, we use a predefined admin key from Hardhat's first account
  ## This key is used to deploy contracts in the contract template repo
  ADMIN_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
  echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure admin $ADMIN_KEY --home $PELLDVS_HOME >/dev/null

  export ADMIN_ADDRESS=$(pelldvs keys show admin --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | jq -r .address)
}

function update_pelldvs_config {
  pelldvs init --home $PELLDVS_HOME
  ## Update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellRegistryRouterFactory.json" | jq -r .address)
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
  REGISTRY_ROUTER_ADDRESS_FILE="/root/RegistryRouterAddress.json"
  pelldvs client dvs create-registry-router \
    --home $PELLDVS_HOME \
    --from admin \
    --initialOwner $ADMIN_ADDRESS \
    --dvsChainApprover $ADMIN_ADDRESS \
    --churnApprover $ADMIN_ADDRESS \
    --ejector $ADMIN_ADDRESS \
    --pauser $ADMIN_ADDRESS \
    --unpauser $ADMIN_ADDRESS \
    --initialPausedStatus false \
    --saveToFile $REGISTRY_ROUTER_ADDRESS_FILE \
    --forceSave true

  ## Export registry router address
  export PELL_REGISTRY_ROUTER=$(cat $REGISTRY_ROUTER_ADDRESS_FILE | jq -r .address)
}

function init_pell_emulator {
  ## Initialize pell_emulator and config will be written in /root/.pelldvs/pell_emulator/contract.address.json
  pelldvs pell_emulator init \
    --home $PELLDVS_HOME

  ## Get contracts addresses from Hardhat
  PELL_DELEGATION_MNAGER=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellDelegationManager-Proxy.json" | jq -r .address)
  PELL_DVS_DIRECTORY=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellDVSDirectory-Proxy.json" | jq -r .address)
  PELL_STRATEGY_MANAGER=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/PellStrategyManager-Proxy.json" | jq -r .address)
  PELL_REGISTRY_INTERACTOR=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/RegistryInteractor-Proxy.json" | jq -r .address)
  STAKING_STRATEGY_MANAGER=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/StrategyManager-Proxy.json" | jq -r .address)
  STAKING_DELEGATION_MANAGER=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/DelegationManager-Proxy.json" | jq -r .address)
  SERVICE_OMNI_OPERATOR_SHARES_MANAGER=$(ssh hardhat "cat $HARDHAT_CONTRACTS_PATH/OmniOperatorSharesManager-Proxy.json" | jq -r .address)

  DVS_BLS_APK_REGISTRY=$(ssh hardhat "cat $HARDHAT_DVS_PATH/BLSApkRegistry-Proxy.json" | jq -r .address)
  DVS_REGISTRY_COORDINATOR=$(ssh hardhat "cat $HARDHAT_DVS_PATH/RegistryCoordinator-Proxy.json" | jq -r .address)
  DVS_INDEX_REGISTRY=$(ssh hardhat "cat $HARDHAT_DVS_PATH/IndexRegistry-Proxy.json" | jq -r .address)
  DVS_STAKE_REGISTRY=$(ssh hardhat "cat $HARDHAT_DVS_PATH/StakeRegistry-Proxy.json" | jq -r .address)

  ## Update emulator contracts addresses
  update-emulator-address() {
    JQ_EXPR="$1"
    jq "$JQ_EXPR" /root/.pelldvs/pell_emulator/contract.address.json >/tmp/tmp.json &&
      mv /tmp/tmp.json /root/.pelldvs/pell_emulator/contract.address.json
  }

  update-emulator-address '.PellDelegationManager = "'$PELL_DELEGATION_MNAGER'"'
  update-emulator-address '.PellDVSDirectory = "'$PELL_DVS_DIRECTORY'"'
  update-emulator-address '.PellStrategyManager = "'$PELL_STRATEGY_MANAGER'"'
  update-emulator-address '.PellRegistryInteractor = "'PELL_REGISTRY_INTERACTOR'"'
  update-emulator-address '.PellRegistryRouter = "'$PELL_REGISTRY_ROUTER'"'
  update-emulator-address '.StakingStrategyManager = "'$STAKING_STRATEGY_MANAGER'"'
  update-emulator-address '.StakingDelegationManager = "'$STAKING_DELEGATION_MANAGER'"'
  update-emulator-address '.ServiceOmniOperatorSharesManager = "'$SERVICE_OMNI_OPERATOR_SHARES_MANAGER'"'

  update-emulator-address '.DVSBLSApkRegistry = "'$DVS_BLS_APK_REGISTRY'"'
  update-emulator-address '.DVSRegistryCoordinator = "'$DVS_REGISTRY_COORDINATOR'"'
  update-emulator-address '.DVSIndexRegistry = "'$DVS_INDEX_REGISTRY'"'
  update-emulator-address '.DVSStakeRegistry = "'$DVS_STAKE_REGISTRY'"'
  cat /root/.pelldvs/pell_emulator/contract.address.json | jq

}

function start_pell_emulator {
  ## start pell_emulator
  pelldvs pell_emulator start \
    --home "$PELLDVS_HOME" \
    --rpc-url "$ETH_RPC_URL" \
    --ws-url "$ETH_WS_URL" \
    --deployer-key-file "$ADMIN_KEY_FILE"
}

# start sshd
/usr/sbin/sshd &

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Wait for Hardhat to be ready"
hardhat_healthcheck

logt "Update PellDVS Config"
update_pelldvs_config

if [ ! -f /root/emulator_initialized ]; then
  logt "Setup Admin Key"
  setup_admin_key

  logt "Create Registry Router"
  create_registry_router

  logt "Initialize Pell Emulator"
  init_pell_emulator

  touch /root/emulator_initialized
else
  logt "Pell Emulator already initialized, skipping..."
fi

logt "Start Pell Emulator"
start_pell_emulator
