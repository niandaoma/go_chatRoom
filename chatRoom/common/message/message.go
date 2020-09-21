package message

const (
	LoginMesType            string = "LoginMes"
	LoginResMesType         string = "LoginResMes"
	RegisterMesType         string = "RegisterMes"
	RegisterResMesType      string = "RegisterMes"
	NotifyUserStatusMesType string = "NotifyUserStatusMes"
	SmsMesType              string = "SmsMes"
	SmsP2PMesType           string = "SmsP2PMes"
)

const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息的类型
	Data string `json:"data"` //数据的类型
}

//定义两个消息，后面需要再增加

type LoginMes struct {
	UserId   int    `json:"userId"`   //用户ID
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户姓名
}

type LoginResMes struct {
	Code    int    `json:"code"` //返回的状态码   500表示该用户未注册	200表示登陆成功
	UserIds []int  //增加字段，保存用户id的切片
	Error   string `json:"error"` //返回错误信息
}

type RegisterMes struct {
	User User
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回的状态码   500表示该用户未注册	200表示登陆成功
	Error string `json:"error"` //返回错误信息
}

//新定一个类型来推送用户状态
type NotifyUserStatusMes struct {
	UserId     int `json:"userId"`     //用户ID
	UserStatus int `json:"userStatus"` //用户状态
}

//增加一个 SmsMes //发送的信息
type SmsMes struct {
	Content string `json:"content"` //发送的消息内容
	User           //继承
}

type SmsP2PMes struct {
	Content string `json:"content"` //发送的消息内容
	SendId  int
	User    //继承
}
