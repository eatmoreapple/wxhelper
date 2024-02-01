FROM golang:1.21 as builder

WORKDIR /app

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-s -w" -o /app/apiserver apiserver/cmd/main.go


FROM furacas/wine-vnc-box:latest

# 根据传入参数安装微信和wxhelper.dll
ARG WECHAT_URL=https://github.com/tom-snow/wechat-windows-versions/releases/download/v3.9.5.81/WeChatSetup-3.9.5.81.exe
ARG WXHELPER_URL=https://github.com/ttttupup/wxhelper/releases/download/3.9.5.81-v11/wxhelper.dll

WORKDIR /home/app/.wine/drive_c

# 加载注入器
ADD https://github.com/furacas/DllInjector/releases/download/v1.4.0/DllInjector64.exe DllInjector.exe
RUN sudo chown app:app DllInjector.exe && sudo chmod a+x DllInjector.exe

# 下载微信
ADD ${WECHAT_URL} WeChatSetup.exe
RUN sudo chown app:app WeChatSetup.exe  && sudo chmod a+x WeChatSetup.exe
RUN ls -lah

# 安装微信
COPY apiserver/cmd/install-wechat.sh install-wechat.sh
RUN bash -c 'nohup /entrypoint.sh 2>&1 &' && sleep 10
RUN sudo apt install -y net-tools
RUN sudo chmod a+x install-wechat.sh && ./install-wechat.sh
RUN rm -rf WeChatSetup.exe && rm -rf install-wechat.sh

# 下载wxhelper.dll
ADD ${WXHELPER_URL} wxhelper.dll
RUN sudo chown app:app wxhelper.dll

COPY --from=builder /app/apiserver /home/app/.wine/drive_c/apiserver

COPY apiserver/cmd/cmd.sh /cmd.sh
CMD ["/cmd.sh"]