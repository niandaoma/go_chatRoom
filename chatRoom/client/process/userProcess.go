package process

import (
	"demo/chatRoom/client/utils"
	"demo/chatRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//暂时不需要字段
}

func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	//下一步就要定协议

	//获得链接
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	//准备通过conn发送消息给服务器

	var mes message.Message
	mes.Type = message.LoginMesType

	//创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//将login序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	//把data赋给mes.Data
	mes.Data = string(data)
	//把mes.Data序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	//先把data转成[]byte
	var pkgLen uint32
	//获得data长度
	pkgLen = uint32(len(data))
	var buf = make([]byte, 4)
	//将pkhLen转成[]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	n, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	tf := &utils.Transfer{
		Conn: conn,
		Buf:  [8096]byte{},
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg() fail err=", err)
		return
	}
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		return
	}
	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		fmt.Println("当前在线用户列表如下")
		for _, v := range loginResMes.UserIds {
			//如果要求不显示自己在线，可以加以下代码
			if v == userId {
				continue
			}
			fmt.Println("在线用户id:\t", v)
			//完成 客户端 OnlineUsers 初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		//这里我们还需要在客户端启动一个协程
		//这里协程保持和服务器端的通讯，如果服务器有数据推送给客户端
		//则接受并显示在客户端

		go serverProcessMes(conn)
		ShowLoginMenu()
		//登陆成功，显示我们登陆成功的菜单

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	//获得链接
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()
	//准备通过conn发送消息给服务器

	var mes message.Message
	mes.Type = message.RegisterMesType

	//创建一个RegisterMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//将login序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	//把data赋给mes.Data
	mes.Data = string(data)

	//把mes.Data序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
		Buf:  [8096]byte{},
	}
	//发送数据给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("WritePkg() fail err=", err)
		return
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg() fail err=", err)
		return
	}
	var RegisterResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterResMes)
	if err != nil {
		return
	}
	if RegisterResMes.Code == 200 {

		//这里我们还需要在客户端启动一个协程
		//这里协程保持和服务器端的通讯，如果服务器有数据推送给客户端
		//则接受并显示在客户端
		go serverProcessMes(conn)
		ShowRegisterMenu()
	} else {
		fmt.Println(RegisterResMes.Error)
	}
	return

}

func (this *UserProcess) ClientExit() (err error) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = CurUser.UserId
	notifyUserStatusMes.UserStatus = message.UserOffline
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		return
	}
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		return
	}
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		return
	}
	return
}
