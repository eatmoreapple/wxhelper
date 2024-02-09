package wxhelper

import (
	"github.com/eatmoreapple/wxhelper/apiclient"
	"io"
	"strings"
)

type empty struct{}

var (
	// internalUsers is a map of internal users.
	internalUsers = map[string]empty{
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

	// owner returns the owner of the user.
	owner func() *Account
}

// IsGroup returns whether the user is a group.
func (u *User) IsGroup() bool {
	return strings.Contains(u.Wxid, "@chatroom")
}

// IsFriend returns whether the user is a friend.
func (u *User) IsFriend() bool {
	return u.Type == 3 && u.VerifyFlag == 0 && !u.IsGroup() && !u.IsInternal()
}

// IsInternal returns whether the user is an internal user.
func (u *User) IsInternal() bool {
	_, ok := internalUsers[u.Wxid]
	return ok
}

// IsFileHelper returns whether the user is a file helper.
func (u *User) IsFileHelper() bool {
	return u.Wxid == "filehelper"
}

// Owner returns the owner of the user.
func (u *User) Owner() *Account {
	return u.owner()
}

func (u *User) SendText(content string) error {
	return u.Owner().sendText(u.Wxid, content)
}

func (u *User) SendImage(img io.Reader) error {
	return u.Owner().sendImage(u.Wxid, img)
}

func (u *User) SendFile(file io.Reader) error {
	return u.Owner().sendFile(u.Wxid, file)
}

type Friend struct{ *User }

func (f *Friend) SendText(content string) error {
	return f.Owner().SendTextToFriend(f, content)
}

func (f *Friend) SendImage(img io.Reader) error {
	return f.Owner().SendImageToFriend(f, img)
}

func (f *Friend) SendFile(file io.Reader) error {
	return f.Owner().SendFileToFriend(f, file)
}

type Friends []*Friend

func (f Friends) Search(limit uint, searchFunc func(friend *Friend) bool) Friends {
	var search = make(Friends, 0)
	for _, friend := range f {
		if searchFunc(friend) {
			search = append(search, friend)
			if uint(len(search)) == limit {
				break
			}
		}
	}
	return search
}

func (f Friends) SearchByWxID(wxID string) (*Friend, bool) {
	search := f.Search(1, func(friend *Friend) bool { return friend.Wxid == wxID })
	if len(search) == 0 {
		return nil, false
	}
	return search[0], true
}

func (f Friends) SearchByNickname(nickname string, limit uint) Friends {
	return f.Search(limit, func(friend *Friend) bool { return friend.Nickname == nickname })
}

func (f Friends) SearchByRemark(remark string, limit uint) Friends {
	return f.Search(limit, func(friend *Friend) bool { return friend.Remark == remark })
}

type Group struct{ *User }

// IsInContactList returns whether the group is in the contact list.
func (g *Group) IsInContactList() bool {
	return len(g.EncryptName) == 0
}

func (g *Group) SendText(content string) error {
	return g.Owner().SendTextToGroup(g, content)
}

func (g *Group) SendImage(img io.Reader) error {
	return g.Owner().SendImageToGroup(g, img)
}

func (g *Group) SendFile(file io.Reader) error {
	return g.Owner().SendFileToGroup(g, file)
}

func (g *Group) Members() ([]*Profile, error) {
	return g.Owner().bot.client.GetChatRoomMembers(g.Owner().bot.Context(), g.Wxid)
}

func (g *Group) SendAtText(content string, memberIDs ...string) error {
	return g.Owner().bot.client.SendAtText(g.Owner().bot.Context(), apiclient.SendAtTextOption{
		GroupID: g.Wxid,
		AtList:  memberIDs,
		Content: content,
	})
}

func (g *Group) SendAtALLTextMsg(content string) error {
	return g.SendAtText(content, "notify@all")
}

type GroupInfo struct {
	ChatRoomID string
	Notice     string
	Admin      string
	XML        string
}

func (g *Group) Info() (*GroupInfo, error) {
	return g.Owner().bot.client.GetChatRoomInfo(g.Owner().bot.Context(), g.Wxid)
}

func (g Groups) Search(limit uint, searchFunc func(group *Group) bool) Groups {
	var search = make(Groups, 0)
	for _, group := range g {
		if searchFunc(group) {
			search = append(search, group)
			if uint(len(search)) == limit {
				break
			}
		}
	}
	return search
}

type Groups []*Group

type Members []*User

func (l Members) Friends() Friends {
	members := l.Search(uint(len(l)), func(user *User) bool { return user.IsFriend() })
	var friends = make(Friends, 0, len(members))
	for _, user := range members {
		friends = append(friends, &Friend{User: user})
	}
	return friends
}

func (l Members) Groups() Groups {
	members := l.Search(uint(len(l)), func(user *User) bool { return user.IsGroup() })
	var groups = make(Groups, 0, len(members))
	for _, user := range members {
		groups = append(groups, &Group{User: user})
	}
	return groups
}

func (l Members) Len() int {
	return len(l)
}

func (l Members) Less(i, j int) bool {
	return l[i].Pinyin < l[j].Pinyin
}

func (l Members) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Members) Search(limit uint, searchFunc func(user *User) bool) Members {
	var search = make(Members, 0)
	for _, user := range l {
		if searchFunc(user) {
			search = append(search, user)
			if uint(len(search)) == limit {
				break
			}
		}
	}
	return search
}

type Profile struct {
	Account   string `json:"account"`
	HeadImage string `json:"headImage"`
	Nickname  string `json:"nickname"`
	V3        string `json:"v3"`
	Wxid      string `json:"wxid"`
}
