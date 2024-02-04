package wxhelper

import "io"

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

type MessageHandler func(msg *Message)
