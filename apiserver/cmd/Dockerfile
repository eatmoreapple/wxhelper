FROM golang:1.21 as builder

WORKDIR /app

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build --ldflags="-s -w" -o /app/apiserverd ./apiserver/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/apiserverd /app/apiserver

RUN apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /application/zoneinfo.zip

ENV ZONEINFO=/app/zoneinfo.zip

RUN chmod +x /app/apiserver

EXPOSE 19089

CMD ["/app/apiserver"]