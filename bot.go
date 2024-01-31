package wxhelper

import (
	"context"
	"github.com/eatmoreapple/wxhelper/apiclient"
)

type Bot struct {
	MessageHandler   MessageHandler
	messageRetriever MessageRetriever
	apiclient        *apiclient.Client
	ctx              context.Context
	stop             func()
}

func (b *Bot) Context() context.Context { return b.ctx }

func (b *Bot) GetLoginAccount() (*Account, error) {
	account, err := b.apiclient.GetUserInfo(b.ctx)
	if err != nil {
		return nil, err
	}
	return &Account{
		Account:         account.Account,
		City:            account.City,
		Country:         account.Country,
		CurrentDataPath: account.CurrentDataPath,
		DataSavePath:    account.DataSavePath,
		DbKey:           account.DbKey,
		HeadImage:       account.HeadImage,
		Mobile:          account.Mobile,
		Name:            account.Name,
		Province:        account.Province,
		Signature:       account.Signature,
		Wxid:            account.Wxid,
		PrivateKey:      account.PrivateKey,
		PublicKey:       account.PublicKey,
		bot:             b,
	}, nil
}

func (b *Bot) Run() error {
	messageChan, err := b.messageRetriever.RetrieveMessage()
	if err != nil {
		return err
	}
	for {
		select {
		case <-b.ctx.Done():
			return b.Stop()
		case msg := <-messageChan:
			if b.MessageHandler != nil {
				go b.MessageHandler(msg)
			}
		}
	}
}

func (b *Bot) Stop() error {
	b.stop()
	return b.ctx.Err()
}

//func New(ctx context.Context, client Client, options ...HttpMessageRetrieverOptionFunc) *Bot {
//	ctx, cancel := context.WithCancel(ctx)
//	messageRetriever := NewHttpMessageRetriever(ctx, client, options...)
//	return &Bot{
//		ctx:              ctx,
//		stop:             cancel,
//		client:           client,
//		messageRetriever: messageRetriever,
//	}
//}
