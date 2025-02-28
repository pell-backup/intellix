FROM golang:1.23 AS build

ARG GITHUB_TOKEN
ARG HTTP_PROXY
ARG HTTPS_PROXY

## Setup env
WORKDIR /app

ARG GITHUB_TOKEN
ARG PELLDVS_VERSION
RUN if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN is not set" && exit 1; fi
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork"
RUN git clone https://github.com/0xPellNetwork/pelldvs.git --branch $PELLDVS_VERSION --depth 1

## Build pelldvs
WORKDIR /app/pelldvs
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    make pelldvs-build

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/IntelliXLabs".insteadOf "https://github.com/IntelliXLabs"

########## Build intellixd ##########
ARG DEBUG_ENABLED
WORKDIR /app/intellixd
COPY go.mod go.sum ./
COPY scripts ./scripts
RUN bash scripts/install-lib.sh
RUN --mount=type=cache,target="/go/pkg/mod" go mod download

COPY .. /app/intellixd/
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    if [ "$DEBUG_ENABLED" = "true" ]; then \
        make build-debug; \
    else \
        make build; \
    fi

########## Setup runtime env ##########
FROM golang:1.23 AS runtime
RUN apt-get update -yqq && apt-get install -yqq openssh-server curl jq less
RUN mkdir -p /run/sshd && chmod 0755 /run/sshd

RUN curl -L https://foundry.paradigm.xyz | bash
ENV PATH="$PATH:/root/.foundry/bin"
RUN foundryup

RUN go install github.com/mikefarah/yq/v4@latest
RUN go install github.com/tomwright/dasel/v2/cmd/dasel@master
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY --from=build /app/intellixd/build/intellixd /usr/bin/intellixd
COPY --from=build /app/intellixd/lib/* /usr/local/lib/
COPY --from=build /app/pelldvs/build/pelldvs /usr/bin/pelldvs
COPY ./docker/operator/scripts /root/scripts
COPY ./docker/ssh /root/.ssh

RUN ldconfig /usr/local/lib
RUN chmod 600 /root/.ssh/*

WORKDIR /root

# expose port
EXPOSE 22