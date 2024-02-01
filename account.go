package wxhelper

import (
	"context"
	"io"
)

type Account struct {
	Account         string `json:"account"`
	City            string `json:"city"`
	Country         string `json:"country"`
	CurrentDataPath string `json:"currentDataPath"`
	DataSavePath    string `json:"dataSavePath"`
	DbKey           string `json:"dbKey"`
	HeadImage       string `json:"headImage"`
	Mobile          string `json:"mobile"`
	Name            string `json:"name"`
	Province        string `json:"province"`
	Signature       string `json:"signature"`
	Wxid            string `json:"wxid"`
	PrivateKey      string `json:"privateKey"`
	PublicKey       string `json:"publicKey"`
	bot             *Bot
	friends         Friends
	groups          Groups
	fileHelper      *User
}

func (a *Account) Friends(update ...bool) (Friends, error) {
	if (len(update) > 0 && update[0]) || a.friends == nil {
		members, err := a.bot.client.GetContactList(context.Background())
		if err != nil {
			return nil, err
		}
		a.friends = members.Friends()
		for _, friend := range a.friends {
			friend.User.owner = func() *Account { return a }
		}
	}
	return a.friends, nil
}

func (a *Account) Groups(update ...bool) (Groups, error) {
	if (len(update) > 0 && update[0]) || a.groups == nil {
		members, err := a.bot.client.GetContactList(context.Background())
		if err != nil {
			return nil, err
		}
		a.groups = members.Groups()
		for _, group := range a.groups {
			group.User.owner = func() *Account { return a }
		}
	}
	return a.groups, nil
}

func (a *Account) FileHelper() *User {
	if a.fileHelper == nil {
		a.fileHelper = &User{Wxid: "filehelper", owner: func() *Account { return a }}
	}
	return a.fileHelper
}

func (a *Account) sendText(wxID string, content string) error {
	return a.bot.client.SendText(a.bot.Context(), wxID, content)
}

func (a *Account) sendImage(account string, img io.Reader) error {
	return a.bot.client.SendImage(a.bot.Context(), account, img)
}

func (a *Account) SendTextToFriend(friend *Friend, content string) error {
	return a.sendText(friend.User.Wxid, content)
}

func (a *Account) SendImageToFriend(friend *Friend, img io.Reader) error {
	return a.sendImage(friend.User.Wxid, img)
}

func (a *Account) SendTextToGroup(group *Group, content string) error {
	return a.sendText(group.User.Wxid, content)
}

func (a *Account) SendImageToGroup(group *Group, img io.Reader) error {
	return a.sendImage(group.User.Wxid, img)
}
