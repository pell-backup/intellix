set -x

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/price-oracle-dvs/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/price-oracle-dvs/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  export COSMOS_NODE_URI=${COSMOS_NODE_URI:-http://abci:26657}
}

function abci_healthcheck {
  set +e
  while true; do
    curl -s $COSMOS_NODE_URI >/dev/null
    if [ $? -eq 0 ]; then
      echo "Cosmos ABCI is ready, proceeding to the next step..."
      break
    fi
    echo "Cosmos ABCI is not ready, retrying in 2 second..."
    sleep 2
  done
  set -e
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
  sleep 15
  set -e
}

function assert_gt {
  # use awk to compare numbers as big numbers are not supported by bash
  if [ $(awk "BEGIN {print ($1 > $2)}") -ne 1 ]; then
    echo "[FAIL] Expected $1 to be greater than $2"
    exit 1
  fi
  echo "[PASS] Expected $1 to be greater than $2"
}

load_defaults
operator_healthcheck
abci_healthcheck

ADMIN_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/PriceOraclePayInNativeConsumer.json" | jq -r .address)

export TIMEOUT_FOR_TASK_PROCESS=${TIMEOUT_FOR_TASK_PROCESS:-20}

### ----------------
## create a new task
cast send "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "requestPrice(string)" "ETH" --private-key "$ADMIN_KEY" --rpc-url "$ETH_RPC_URL"

# wait for the task to be processed
echo "wait ${TIMEOUT_FOR_TASK_PROCESS} seconds for the task to be processed"
sleep ${TIMEOUT_FOR_TASK_PROCESS}
RESULT=$(cast call "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "price()" --private-key "$ADMIN_KEY" --rpc-url "$ETH_RPC_URL" | cast to-dec)
assert_gt "$RESULT" "0"

# cast call "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "allTaskResponses(uint32)" $((TASK_NUMBER - 1))
# RETRIEVER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorStateRetriever.json" | jq -r .address)
# cast call "$RETRIEVER_ADDRESS" "GetGROUPsDVSStateAtBlock(uint32)" $TASK_ID --private-key "$ADMIN_KEY"

# update operator socket address, in init_operator.sh we have set the operator socket address to http://$(hostname -i):26657
# here we update it to http://operator:26657
logt "update operator sokcet address and sleep 5 seconds to wait for the dispatcher has processed it "
ssh operator "bash /root/scripts/update_operator_socket.sh http://operator:26657"
logt ""

### ----------------
## create a new task
cast send "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "requestPrice(string)" "BTC" --private-key "$ADMIN_KEY" --rpc-url "$ETH_RPC_URL"

# wait for the task to be processed
echo "wait ${TIMEOUT_FOR_TASK_PROCESS} seconds for the task to be processed"
sleep ${TIMEOUT_FOR_TASK_PROCESS}
RESULT=$(cast call "$PRICE_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "price()" --private-key "$ADMIN_KEY" --rpc-url "$ETH_RPC_URL" | cast to-dec)
assert_gt "$RESULT" "0"

### ----------------
## create a new vrf task
VRF_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/VRFOraclePayInNativeConsumer.json" | jq -r .address)

# cast send "$VRF_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" "requestRandomWords(uint256)" 2 --private-key "$ADMIN_KEY" --rpc-url "$ETH_RPC_URL"

TX_HASH=$(
  cast send "$VRF_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" \
    "requestRandomWords(uint256)" 2 \
    --private-key "$ADMIN_KEY" \
    --rpc-url "$ETH_RPC_URL" \
    --json \
  | jq -r .transactionHash
)

echo "transactionHash: $TX_HASH"

sleep ${TIMEOUT_FOR_TASK_PROCESS}
cast tx "$TX_HASH" --rpc-url "$ETH_RPC_URL"

LOGS_JSON=$(
  cast receipt "$TX_HASH" --rpc-url "$ETH_RPC_URL" --json
)

DATA=$(
  echo "$LOGS_JSON" \
    | jq -r '.logs[]
      | select(.topics[0] == "0x7106aad702536308648985e660b725d368593a51be2b5134b4c1b2f884492c86")
      | .data'
)

REQUEST_ID="${DATA:0:66}"

echo "Got requestId: $REQUEST_ID"

# Run cast and capture the output
output=$(cast call \
  "$VRF_ORACLE_PAY_IN_NATIVE_CONSUMER_ADDRESS" \
  "getRequestStatus(bytes32)(bool,uint256[])" \
  "$REQUEST_ID" \
  --rpc-url "$ETH_RPC_URL")

echo "Output: $output"

if [[ "$output" == $'false\n[]' ]]; then
    echo "Error: Output was 'false\n[]', exiting..."
    exit 1
else
    echo "Everything is okay."
fi
