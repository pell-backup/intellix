#!/bin/bash
# Wait for hardhat to start
sleep 3

rm -rf /app/deployments/*

# deploy pell evm
cd ./lib/pell-middleware-contracts/lib/pell-contracts
npx hardhat deploy --deploy-scripts deploy_restaking --network localhost
npx hardhat deploy --deploy-scripts deploy_pell --network localhost
npx hardhat deploy --deploy-scripts deploy_service_omni --network localhost
npx hardhat update-delegation-connector --network localhost

# deploy pell middleware
cd ../..
npx hardhat --network localhost deploy

# deploy price dvs
cd ../..
npx hardhat --network localhost deploy

# listen pell events
cd lib/pell-middleware-contracts/lib/pell-contracts
npx hardhat --network localhost listen-pell-events
