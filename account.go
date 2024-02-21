package wxhelper

import (
	"context"
	"errors"
	"io"
	"strconv"
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
}

func (a *Account) Friends() (Friends, error) {
	members, err := a.bot.client.GetContactList(context.Background())
	if err != nil {
		return nil, err
	}
	friends := members.Friends()
	for _, friend := range friends {
		friend.User.owner = func() *Account { return a }
	}
	return friends, nil
}

func (a *Account) Groups() (Groups, error) {
	members, err := a.bot.client.GetContactList(context.Background())
	if err != nil {
		return nil, err
	}
	groups := members.Groups()
	for _, group := range groups {
		group.User.owner = func() *Account { return a }
	}
	return groups, nil
}

func (a *Account) FileHelper() *User {
	return &User{Wxid: "filehelper", owner: func() *Account { return a }}
}

func (a *Account) sendText(wxID string, content string) error {
	return a.bot.client.SendText(a.bot.Context(), wxID, content)
}

func (a *Account) sendImage(account string, img io.Reader) error {
	return a.bot.client.SendImage(a.bot.Context(), account, img)
}

func (a *Account) sendFile(account string, file io.Reader) error {
	return a.bot.client.SendFile(a.bot.Context(), account, file)
}

func (a *Account) SendTextToFriend(friend *Friend, content string) error {
	return a.sendText(friend.User.Wxid, content)
}

func (a *Account) SendImageToFriend(friend *Friend, img io.Reader) error {
	return a.sendImage(friend.User.Wxid, img)
}

func (a *Account) SendFileToFriend(friend *Friend, file io.Reader) error {
	return a.sendFile(friend.User.Wxid, file)
}

func (a *Account) SendTextToGroup(group *Group, content string) error {
	return a.sendText(group.User.Wxid, content)
}

func (a *Account) SendImageToGroup(group *Group, img io.Reader) error {
	return a.sendImage(group.User.Wxid, img)
}

func (a *Account) SendFileToGroup(group *Group, file io.Reader) error {
	return a.sendFile(group.User.Wxid, file)
}

func (a *Account) AddMemberIntoChatRoom(group *Group, users ...*Friend) error {
	if len(users) == 0 {
		return errors.New("no user to add")
	}
	// 判断群聊人数有没有超过40人
	members, err := group.Members()
	if err != nil {
		return err
	}
	var wxIds = make([]string, 0, len(users))
	for _, user := range users {
		wxIds = append(wxIds, user.Wxid)
	}
	if len(members) > 40 {
		return a.bot.client.InviteMemberToChatRoom(a.bot.Context(), group.User.Wxid, wxIds)
	} else {
		return a.bot.client.AddMemberIntoChatRoom(a.bot.Context(), group.User.Wxid, wxIds)
	}
}

func (a *Account) ForwardMessage(msg *Message, user *User) error {
	return a.bot.client.ForwardMsg(a.bot.Context(), user.Wxid, strconv.FormatInt(msg.MsgId, 10))
}
