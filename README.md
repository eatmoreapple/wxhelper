## 🚀 wxhelper: 微信机器人，让你的微信更机智！ 🤖

想让你的微信自动回复消息吗？想要拥有一个个人微信助手，让你的朋友眼前一亮吗？来试试wxhelper，它能让你的微信聊天像打乒乓球一样有来有回！

### 🎮 玩法指南：


1. **先行一步**：跳到这个神秘的地方 [这里](https://github.com/eatmoreapple/wxapiserver) 。像宝藏猎人一样，找到那份藏着无尽可能的docker-compose.yml文件，拷贝它到你的本地宝藏库（也就是你的电脑）。

2. **唤醒机器人**：在你的电脑上施展docker-compose up -d咒语，召唤出服务的灵魂。

3. **耐心等待**：服务启动需要一点时间，就像煮面条一样，你不能急。等面条煮好了，就可以访问 http://localhost:8080 ，那里有一个神奇的页面。在页面上找到那个vnc的秘密通道，点击进入，然后像电影里的特工一样，扫码登录。

4. **操纵机器人**：复制下面的魔法代码，它是用Go语言编写的。这段代码能让你的微信机器人活起来，它会在你的命令下工作。

    ```go
    package main
    
    import "github.com/eatmoreapple/wxhelper"
    
    func main() {
        bot := wxhelper.New("http://localhost:19089")
        bot.MessageHandler = func(msg *wxhelper.Message) {
            if msg.Content == "ping" {
                msg.ReplyText("pong")
            }
        }
        bot.Run()
    }
    ```

5. **开始游戏**：给你登录的微信号发一个ping，看看会发生什么？没错，你会得到一个pong的回复，就像魔法一样！

### 🎉 来吧，开始你的微信机器人之旅吧！
把你的微信变成一个不眠不休的回复机器，让它在你忙于拯救世界的时候，继续和你的朋友们愉快地聊天吧！别忘了，用这个超级简单的wxhelper，你的微信将变得比任何时候都更加机智。现在就开始吧，成为你朋友圈里的科技达人！💡👾