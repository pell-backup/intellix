
set -x

function load_defaults {
  export ADMIN_KEY=${ADMIN_KEY}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
}

function setup_operator_key {
  # echo  -ne '\n\n' | pelldvs keys create operator --key-type=ecdsa --insecure > /tmp/operator.key
  if ! pelldvs keys show operator --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure operator $OPERATOR_KEY --home $PELLDVS_HOME >/dev/null
  fi

  export OPERATOR_ADDRESS=$(pelldvs keys show operator --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)

  ## To register operator in the DVS, we need the operator's BLS key with the same name
  if ! pelldvs keys show operator --key-type bls --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo  -ne '\n\n' | pelldvs keys create operator --key-type=bls --insecure > /tmp/operator_bls.key
  fi

}

load_defaults
setup_operator_key
