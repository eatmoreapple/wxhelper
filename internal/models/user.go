package models

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

type GroupMember struct {
	ChatRoomID     string `json:"chatRoomId"`
	Members        string `json:"members"`
	MemberNickname string `json:"memberNickname"`
	Admin          string `json:"admin"`
	AdminNickname  string `json:"adminNickname"`
}

type Profile struct {
	Account   string `json:"account"`
	HeadImage string `json:"headImage"`
	Nickname  string `json:"nickname"`
	V3        string `json:"v3"`
	Wxid      string `json:"wxid"`
}
