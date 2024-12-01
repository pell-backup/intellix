FROM golang:1.22-bullseye AS go
FROM rust:1.82.0-bullseye

ENV PATH=$PATH:/usr/local/go/bin
COPY --from=go /usr/local/go /usr/local/go

ARG HTTP_PROXY
ARG HTTPS_PROXY

## Setup env
WORKDIR /app

ARG GITHUB_TOKEN
ARG IWASM_VERSION
RUN if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN is not set" && exit 1; fi
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/IntelliXLabs".insteadOf "https://github.com/IntelliXLabs"
RUN git clone https://github.com/IntelliXLabs/iwasm.git --branch $IWASM_VERSION --depth 1


## Build iwasm
WORKDIR /app/iwasm
ENV CARGO_HOME=/root/.cargo
RUN --mount=type=cache,target="/root/.cargo" \
    make build

WORKDIR /root

# expose port
EXPOSE 22