package wxhelper

import (
	"context"
	"github.com/eatmoreapple/wxhelper/apiclient"
)

type Bot struct {
	MessageHandler MessageHandler
	client         *Client
	ctx            context.Context
	stop           func()
}

func (b *Bot) Context() context.Context { return b.ctx }

func (b *Bot) GetLoginAccount() (*Account, error) {
	account, err := b.client.GetUserInfo(b.ctx)
	if err != nil {
		return nil, err
	}
	account.bot = b
	return account, nil
}

func (b *Bot) syncMessage() error {
	account, err := b.GetLoginAccount()
	if err != nil {
		return err
	}
	for {
		select {
		case <-b.ctx.Done():
			return b.ctx.Err()
		default:
		}
		message, err := b.client.SyncMessage(b.ctx)
		if err != nil {
			return err
		}
		for _, msg := range message {
			msg.account = account
			if b.MessageHandler != nil {
				go b.MessageHandler(msg)
			}
		}
	}
}

func (b *Bot) Run() error {
	return b.syncMessage()
	//messageChan, err := b.messageRetriever.RetrieveMessage()
	//if err != nil {
	//	return err
	//}
	//for {
	//	select {
	//	case <-b.ctx.Done():
	//		return b.Stop()
	//	case msg := <-messageChan:
	//		if b.MessageHandler != nil {
	//			go b.MessageHandler(msg)
	//		}
	//	}
	//}
}

func New(apiServerURL string) *Bot {
	bot := &Bot{
		client: &Client{
			apiclient: apiclient.New(apiServerURL),
		},
	}
	bot.ctx, bot.stop = context.WithCancel(context.Background())
	return bot
}
