FROM golang:1.12.4-alpine3.9 AS builder

WORKDIR $GOPATH/src/github.com/v2af/aliyun_ddns

COPY . .

ARG VERSION="unset"

RUN apk add --no-cache \ 
    git \
    tzdata \
    gcc \
    g++ && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone && \
    apk del tzdata

RUN DATE="$(date -u +%Y-%m-%d-%H:%M:%S-%Z)" && GO111MODULE=on CGO_ENABLED=0 GOPROXY="https://proxy.golang.org" go build -ldflags "-X github.com/v2af/aliyun_ddns/build.version=$VERSION -X github.com/v2af/aliyun_ddns/build.buildDate=$DATE" -o /bin/aliddns .

FROM alpine

COPY --from=builder /bin/aliddns /bin/aliddns

RUN apk add --no-cache \ 
    git \
    tzdata \
    gcc \
    g++ && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone && \
    apk del tzdata

CMD ["aliddns", "-c=/config/cfg.json"]