package process2

import (
	"demo/chatRoom/common/message"
	"demo/chatRoom/server/dao"
	"demo/chatRoom/server/model"
	"demo/chatRoom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//分析有什么字段
	Conn net.Conn
	//增加一个字段表明是哪个用户登陆
	UserId int
}

//用来记录用户登陆的信息，方便当用户退出登陆的时候删除在线信息

//处理用户登陆
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//将mes中data反序列化得到LoginMes结构体，再将结构体中字段与数据库中数据进行比对
	var loginMes message.LoginMes
	var loginResMes message.LoginResMes
	//返回的信息
	var res message.Message
	//返回处理登录信息的类型
	res.Type = message.LoginResMesType
	//对数据进行解码
	err1 := json.Unmarshal([]byte(mes.Data), &loginMes)
	if err1 != nil {
		fmt.Println("json.Unmarshal err = ", err1)
		return
	}

	//判断用户id和密码是否正确    在redis数据库里面
	//从redis中获取信息，判断用户输入的用户名和id是否正确
	user, err := dao.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500 //500状态码表示用户不存在
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403//用户密码错误
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else { //表示用户登陆成功
		loginResMes.Code = 200
		//将从redis中获取到的id添加到当前处理进程中
		this.UserId = loginMes.UserId
		//把该用户添加到在线用户列表
		userMgr.AddOnlineUserProcess(this) //把这个用户放入在线列表中
		//遍历服务器中的在线用户列表，将所有在线用户的id存入返回信息中的切片
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}
		fmt.Println(user)
	}

	//开始序列化,这里要先发送数据给客户端
	//将登陆处理信息序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		return
	}
	//把序列化之后的存入res，res有type字段标记该数据类型
	res.Data = string(data)
	data, err = json.Marshal(res)
	if err != nil {
		fmt.Println("json.Marshal fail err = ", err)
		return
	}
	//我们因为使用了分层的模式（MVC）我们先创建一个Transfer实例，然后读取
	//所有客户端与服务器之间的信息交互都通过transfer实现
	tf := &utils.Transfer{
		Conn: this.Conn,
		Buf:  [8096]byte{},
	}
	//把序列化后的res发送给客户端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("json.Marshal fail err = ", err)
		return
	}
	if loginResMes.Code == 200 {
		//通知其他用户，该用户上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)

	}
	return
}

//处理用户注册
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	var registerResMes message.RegisterResMes
	var res message.Message

	//从mes中读取传入数据
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}
	//向redis中创建用户
	err = dao.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 403
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "服务器内部错误"
		}
	} else {
		registerResMes.Code = 200
	}

	//开始序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		return
	}
	res.Data = string(data)
	data, err = json.Marshal(res)
	if err != nil {
		fmt.Println("json.Marshal fail err = ", err)
		return
	}
	//我们因为使用了分层的模式（MVC）我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("json.Marshal fail err = ", err)
		return
	}
	return
}

//处理用户状态发生改变
func (this *UserProcess) ServerProcessChangeStatus(mes *message.Message) (err error) {
	var res message.NotifyUserStatusMes
	err = json.Unmarshal([]byte(mes.Data), &res)
	if err != nil {
		return
	}
	switch res.UserStatus {
	//如果是下线
	case message.UserOffline:
		//通知其他用户该用户下线
		this.NotifyOthersOfflineUser(res.UserId)
	}
	return
}

//向其他在线用户通知有人上线
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历所有在线用户
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//开始逐一通知各个用户
		up.NotifyMeOnline(userId)
	}
}

//向其他在线用户通知有人下线
func (this *UserProcess) NotifyOthersOfflineUser(userId int) {
	//遍历所有在线用户
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//开始逐一通知各个用户
		up.NotifyMeOffline(userId)
	}
}

//发送该用户上线信息给其他用户
func (this *UserProcess) NotifyMeOnline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.UserStatus = message.UserOnline

	//序列化notifyUserStatusMes
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal(notifyUserStatusMes) err =", err)
		return
	}
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) err =", err)
		return
	}
	//把data传给客户端
	tf := &utils.Transfer{
		Conn: this.Conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err =", err)
		return
	}
	return
}

//发送该用户离线信息给其他用户
func (this *UserProcess) NotifyMeOffline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.UserStatus = message.UserOffline

	//序列化notifyUserStatusMes
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal(notifyUserStatusMes) err =", err)
		return
	}
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) err =", err)
		return
	}
	//把data传给客户端
	tf := &utils.Transfer{
		Conn: this.Conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err =", err)
		return
	}
	return
}
