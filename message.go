package wxhelper

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

type Message struct {
	Content            string `json:"content"`
	CreateTime         int    `json:"createTime"`
	DisplayFullContent string `json:"displayFullContent"`
	FromUser           string `json:"fromUser"`
	MsgId              int64  `json:"msgId"`
	MsgSequence        int    `json:"msgSequence"`
	Pid                int    `json:"pid"`
	Signature          string `json:"signature"`
	ToUser             string `json:"toUser"`
	Type               int    `json:"type"`
	Base64Img          string `json:"base64Img,omitempty"`

	account *Account
}

func (m Message) IsText() bool {
	return m.Type == 1
}

func (m Message) IsImage() bool {
	return m.Type == 3
}

func (m Message) IsVideo() bool {
	return m.Type == 43
}

func (m Message) IsEmoticon() bool {
	return m.Type == 47
}

func (m Message) IsVoice() bool {
	return m.Type == 34
}

func (m Message) IsAtMe() bool {
	return strings.HasSuffix(m.DisplayFullContent, "在群聊中@了你")
}

func (m Message) IsGroupMessage() bool {
	return strings.HasSuffix(m.FromUser, "@chatroom")
}

func (m Message) Owner() *Account { return m.account }

func (m Message) ReplyText(text string) error {
	return m.Owner().sendText(m.FromUser, text)
}

func (m Message) ReplyImage(img io.Reader) error {
	return m.Owner().sendImage(m.FromUser, img)
}

func (m Message) ReplyFile(file io.Reader) error {
	return m.Owner().sendFile(m.FromUser, file)
}

func (m Message) SaveImage(writer io.Writer) error {
	if !m.IsImage() {
		return errors.New("not an image message")
	}
	data, err := base64.StdEncoding.DecodeString(m.Base64Img)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (m Message) ForwardTo(u *User) error {
	return m.Owner().ForwardMessage(&m, u)
}

func (m Message) Sender() (*User, error) {
	members, err := m.Owner().bot.client.GetContactList(m.Owner().bot.Context())
	if err != nil {
		return nil, err
	}
	result := members.Search(1, func(user *User) bool { return user.Wxid == m.FromUser })
	if len(result) == 0 {
		return nil, ErrNoSuchUserFound
	}
	return result[0], nil
}

type MessageHandler func(msg *Message)
