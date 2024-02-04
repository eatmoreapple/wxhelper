package models

type ChatRoomInfo struct {
	ChatRoomID string `json:"chatRoomId"`
	Notice     string `json:"notice"`
	Admin      string `json:"admin"`
	XML        string `json:"xml"`
}
