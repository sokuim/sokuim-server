FROM golang:1.25-alpine AS Builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o job ./app/job

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=Builder /app/job .

EXPOSE 6073

CMD ["./job"]