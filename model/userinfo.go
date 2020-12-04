package model

type UserInfo struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Id        string    `json:"id"`
	Roles     [1]string `json:"roles"`
}

type UserInfos []UserInfo
