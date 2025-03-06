FROM golang:1.23-bullseye AS build

## Setup env
WORKDIR /app
RUN --mount=type=secret,id=github_token \
    git config --global url."https://$(cat /run/secrets/github_token):x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork" ; \
    git config --global url."https://$(cat /run/secrets/github_token):x-oauth-basic@github.com/IntelliXLabs".insteadOf "https://github.com/IntelliXLabs"

######### Build pelldvs ######
ARG PELLDVS_VERSION=v0.2.2
RUN git clone https://github.com/0xPellNetwork/pelldvs.git --branch $PELLDVS_VERSION --depth 1
WORKDIR /app/pelldvs
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    make build

########## Build intellixd ##########
WORKDIR /app/intellixd
COPY go.mod go.sum ./
COPY scripts ./scripts
RUN bash scripts/install-lib.sh
RUN --mount=type=cache,target="/go/pkg/mod" go mod download

COPY . /app/intellixd/
RUN --mount=type=cache,target="/go/pkg/mod" \
    --mount=type=cache,target="/root/.cache/go-build" \
    make build

########## Setup runtime env ##########
FROM golang:1.23-bullseye AS runtime
RUN apt-get update -yqq && apt-get install -yqq curl jq less
COPY --from=build /app/intellixd/build/intellixd /usr/bin/intellixd
COPY --from=build /app/pelldvs/build/pelldvs /usr/bin/pelldvs
COPY --from=build /app/intellixd/lib/* /usr/local/lib/
RUN ldconfig /usr/local/lib

WORKDIR /root

ENV PELLDVS_HOME=/root/.pelldvs
VOLUME [ "$PELLDVS_HOME" ]

EXPOSE 26657

ENTRYPOINT ["intellixd"]
