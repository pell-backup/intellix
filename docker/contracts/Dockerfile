FROM node:18

ARG HTTP_PROXY
ARG HTTPS_PROXY

RUN --mount=type=secret,id=github_token \
    if [ ! -s /run/secrets/github_token ]; then echo "github token is not set via secrets" >&2; exit 1; fi; \
    git config --global url."https://$(cat /run/secrets/github_token):x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork" ; \
    git config --global url."https://$(cat /run/secrets/github_token):x-oauth-basic@github.com/IntelliXLabs".insteadOf "https://github.com/IntelliXLabs"

RUN apt-get update -yqq && apt-get install -yqq openssh-server curl jq
RUN mkdir -p /run/sshd && chmod 0755 /run/sshd

########## Setup env ##########
WORKDIR /app

ENV PATH="$PATH:/root/.foundry/bin"
RUN curl -L https://foundry.paradigm.xyz | bash
RUN foundryup


ARG CONTRACT_VERSION
RUN git clone https://github.com/IntelliXLabs/price-oracle-dvs.git -b $CONTRACT_VERSION ./price-oracle-dvs


########## Install contracts ##########
WORKDIR /app/price-oracle-dvs

# install deps
RUN --mount=type=cache,target="/usr/local/share/.cache/yarn/v6" \
    forge install --shallow && yarn install --frozen-lockfile
RUN --mount=type=cache,target="/usr/local/share/.cache/yarn/v6" \
    cd lib/pell-middleware-contracts && yarn install --frozen-lockfile
# install 0.6.6 compiler through hardhat compile as Uniswap v2 contracts are compiled with 0.6.6
RUN --mount=type=cache,target="/usr/local/share/.cache/yarn/v6" \
    cd lib/pell-middleware-contracts/lib/pell-contracts && yarn install --frozen-lockfile && npx hardhat compile
RUN export ACCOUNT_SECRETKEY=0x0000000000000000000000000000000000000000000000000000000000000000 &&\
    forge compile && npx hardhat compile

########## Install dev env ##########
COPY ./contracts/scripts/install_dev_env.sh /app/scripts/install_dev_env.sh
RUN bash /app/scripts/install_dev_env.sh

COPY ./contracts/.env.contracts .env
COPY ./contracts/scripts/docker_init.sh ./scripts/docker_init.sh

COPY ./ssh /root/.ssh
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/*

# expose port
EXPOSE 8545
EXPOSE 22

# run hardhat node
CMD /app/price-oracle-dvs/scripts/docker_init.sh

