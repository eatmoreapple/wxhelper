package wxhelper

import "context"

type Bot struct {
	MessageHandler   MessageHandler
	messageRetriever MessageRetriever
	client           Client
	ctx              context.Context
	stop             func()
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

func New(ctx context.Context, client Client, options ...HttpMessageRetrieverOptionFunc) *Bot {
	ctx, cancel := context.WithCancel(ctx)
	messageRetriever := NewHttpMessageRetriever(ctx, client, options...)
	return &Bot{
		ctx:              ctx,
		stop:             cancel,
		client:           client,
		messageRetriever: messageRetriever,
	}
}
