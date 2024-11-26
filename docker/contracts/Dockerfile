FROM node:18

ARG GITHUB_TOKEN
ARG HTTP_PROXY
ARG HTTPS_PROXY

RUN apt-get update -yqq && apt-get install -yqq openssh-server curl jq
RUN mkdir -p /run/sshd && chmod 0755 /run/sshd

########## Setup env ##########
WORKDIR /app

ENV PATH="$PATH:/root/.foundry/bin"
RUN curl -L https://foundry.paradigm.xyz | bash
RUN foundryup

ARG GITHUB_TOKEN
RUN if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN is not set" && exit 1; fi
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork"
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/IntelliXLabs".insteadOf "https://github.com/IntelliXLabs"
RUN git clone https://github.com/IntelliXLabs/price-oracle-dvs.git -b v0.2.2 ./price-oracle-dvs


########## Install contracts ##########
WORKDIR /app/price-oracle-dvs

# install deps
RUN forge install && yarn install --frozen-lockfile
RUN cd lib/pell-middleware-contracts && yarn install --frozen-lockfile
# install 0.6.6 compiler through hardhat compile as Uniswap v2 contracts are compiled with 0.6.6
RUN cd lib/pell-middleware-contracts/lib/pell-contracts && yarn install --frozen-lockfile && npx hardhat compile
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

