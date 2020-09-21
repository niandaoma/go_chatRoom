package message

type User struct {
	UserId     int    `json:"userId"`     //用户ID
	UserPwd    string `json:"userPwd"`    //用户密码
	UserName   string `json:"userName"`   //用户姓名
	UserStatus int    `json:"userStatus"` //用户状态
}
