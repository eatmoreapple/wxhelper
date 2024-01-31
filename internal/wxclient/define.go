package wxclient

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
}

var (
	// internalUsers is a map of internal users.
	internalUsers = map[string]struct{}{
		"filehelper":            {},
		"newsapp":               {},
		"fmessage":              {},
		"weibo":                 {},
		"qqmail":                {},
		"tmessage":              {},
		"qmessage":              {},
		"qqsync":                {},
		"floatbottle":           {},
		"lbsapp":                {},
		"shakeapp":              {},
		"medianote":             {},
		"qqfriend":              {},
		"readerapp":             {},
		"blogapp":               {},
		"facebookapp":           {},
		"masssendapp":           {},
		"meishiapp":             {},
		"feedsapp":              {},
		"voip":                  {},
		"blogappweixin":         {},
		"weixin":                {},
		"brandsessionholder":    {},
		"weixinreminder":        {},
		"officialaccounts":      {},
		"wxitil":                {},
		"userexperience_alarm":  {},
		"notification_messages": {},
		"exmail_tool":           {},
	}
)

type User struct {
	Reserved1     int    `json:"reserved1"`
	Reserved2     int    `json:"reserved2"`
	Type          int    `json:"type"`
	VerifyFlag    int    `json:"verifyFlag"`
	CustomAccount string `json:"customAccount"`
	EncryptName   string `json:"encryptName"`
	Nickname      string `json:"nickname"`
	Pinyin        string `json:"pinyin"`
	PinyinAll     string `json:"pinyinAll"`
	Remark        string `json:"remark"`
	RemarkPinyin  string `json:"remarkPinyin"`
	LabelIds      string `json:"labelIds"`
	Wxid          string `json:"wxid"`
}

type Members []*User

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
}
