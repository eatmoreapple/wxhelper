FROM golang:1.21 as builder

WORKDIR /app

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-s -w" -o /app/apiserver apiserver/cmd/main.go


FROM lxh01/wxhelper-docker:3.9.5.81-v11

COPY --from=builder /app/apiserver /home/app/.wine/drive_c/apiserver

COPY apiserver/cmd/cmd.sh /cmd.sh

RUN chmod +x /cmd.sh

CMD ["/cmd.sh"]