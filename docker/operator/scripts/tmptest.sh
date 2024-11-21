
set -x

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/price-oracle-dvs/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/price-oracle-dvs/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
}

function operator_healthcheck {
  set +e
  while true; do
    ssh operator "test -f /root/operator_initialized"
    if [ $? -eq 0 ]; then
      echo "Operator initialized, proceeding to the next step..."
      break
    fi
    echo "Operator not initialized, retrying in 2 second..."
    sleep 2
  done
  ## Wait for operator to be ready
  sleep 3
  set -e
}

function assert_eq {
  if [ "$1" != "$2" ]; then
    echo "[FAIL] Expected $1 to be equal to $2"
    exit 1
  fi
  echo "[PASS] Expected $1 to be equal to $2"
}

load_defaults
operator_healthcheck

# TODO: replace with actual test
URL=http://operator:8090/hello
response=$(curl -sS -X GET "$URL")
echo "response: $response"

assert_eq "$response" "world"

ADMIN_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/PriceOraclePayInNativeConsumer.json" | jq -r .address)

## create a new task
cast send "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "requestPrice(string)" "BTC" --private-key "$ADMIN_KEY" --rpc-url "$ETH_RPC_URL"

## wait for the task to be processed
export TIMEOUT_FOR_TASK_PROCESS=${TIMEOUT_FOR_TASK_PROCESS:-8}
export TIMEOUT_FOR_TASK_PROCESS=$TIMEOUT_FOR_TASK_PROCESS
echo "wait ${TIMEOUT_FOR_TASK_PROCESS} seconds for the task to be processed"
sleep ${TIMEOUT_FOR_TASK_PROCESS}
RESULT=$(cast call "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "price()" --private-key "$ADMIN_KEY" | cast to-dec
assert_eq "$RESULT" "0"
