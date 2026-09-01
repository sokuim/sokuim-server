FROM golang:1.25-alpine AS Builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG GATEWAY_PORT=8080
RUN mkdir -p /app/etc && \
    sed "s/\${GATEWAY_PORT}/${GATEWAY_PORT}/g" /app/app/gateway/etc/gateway.build.yaml > /app/etc/gateway.yaml

RUN go build -ldflags="-s -w" -o gateway ./app/gateway

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=Builder /app/gateway .
COPY --from=Builder /app/etc/gateway.yaml ./etc/gateway.yaml

EXPOSE 8080

CMD ["./gateway", "-f", "etc/gateway.yaml"]