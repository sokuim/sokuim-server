FROM golang:1.25-alpine AS Builder

WORKDIR /app

ENV GC_ENABLE=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o comet ./app/comet

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=Builder /app/comet .

EXPOSE 6070

CMD ["./comet"]