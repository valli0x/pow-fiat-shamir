FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS build-env

WORKDIR /go/src/github.com/pow-fiat-shamir

COPY . .

ARG TARGETARCH=amd64
ARG TARGETOS=linux
ARG CGO_ENABLED=0

RUN GOARCH=${TARGETARCH} CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} go build -o build/pow-fiat-shamir main.go 

FROM alpine:edge

ENV STORAGE_PASS=""
ENV VAULT_ADDR=""
ENV VAULT_TOKEN=""

RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=build-env /go/src/github.com/pow-fiat-shamir/build/pow-fiat-shamir /usr/bin/pow-fiat-shamir

COPY --from=build-env /go/src/github.com/pow-fiat-shamir/example/config-alice.yml /etc/pow-fiat-shamir-config.yml

ENTRYPOINT ["pow-fiat-shamir"]