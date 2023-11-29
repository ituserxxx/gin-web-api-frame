package model

type LoginRes struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Url   string `json:"url"`   //头像图片
	Phone string `json:"phone"` //手机号
}
