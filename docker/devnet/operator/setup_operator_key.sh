
set -x

function load_defaults {
  export ADMIN_KEY=${ADMIN_KEY}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
}

function setup_operator_key {

  # if "$PELLDVS_HOME"/keys/operator.ecdsa.key.json not exists then import it
  if [ ! -f "$PELLDVS_HOME"/keys/operator.ecdsa.key.json ]; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure operator $OPERATOR_KEY --home $PELLDVS_HOME >/dev/null
  fi

  export OPERATOR_ADDRESS=$(pelldvs keys show operator --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)

  ## To register operator in the DVS, we need the operator's BLS key with the same name
  # if "$PELLDVS_HOME"/keys/operator.bls.key.json not exists then create it
  if [ ! -f "$PELLDVS_HOME"/keys/operator.bls.key.json ]; then
    echo  -ne '\n\n' | pelldvs keys create operator --key-type=bls --insecure > /tmp/operator_bls.key
  fi

}

load_defaults
setup_operator_key
