FROM golang:1.25-alpine AS Builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o comet ./app/comet

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=Builder /app/comet .
COPY --from=Builder /app/app/comet/etc/comet.build.yaml ./etc/comet.yaml

EXPOSE 8082

CMD ["./comet", "-f", "etc/comet.yaml"]