package main

import (
	"demo/chatRoom/common/message"
	"demo/chatRoom/server/process2"
	"demo/chatRoom/server/utils"
	"fmt"
	"io"
	"net"
)

//先创建一个Processor结构体
type Processor struct {
	Conn net.Conn
}

func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {
	//创建一个UserProcess实例
	fmt.Println("mes=", mes)

	switch mes.Type {
	//处理登陆，如果mes的类型是LoginMesType，判断为登陆请求
	case message.LoginMesType:
		up := &process2.UserProcess{
			Conn:   this.Conn,
			UserId: 0,
		}
		//处理用户登陆
		err = up.ServerProcessLogin(mes)
		if err != nil {
			fmt.Println("ServerProcessLogin error = ", err)
		}
	//处理注册，如果mes的类型是RegisterMesType，判断为注册请求
	case message.RegisterMesType:
		up := &process2.UserProcess{
			Conn:   this.Conn,
			UserId: 0,
		}
		//处理用户注册
		err = up.ServerProcessRegister(mes)
	//处理群发，如果mes的类型是SmsMesType，判断为群发请求
	case message.SmsMesType:
		sp := &process2.SmsProcess{}
		sp.SendGroupMes(mes)
	//处理点对点聊天
	case message.SmsP2PMesType:
		sp := &process2.SmsProcess{}
		sp.SendP2PMes(mes)
	case message.NotifyUserStatusMesType:
		up := &process2.UserProcess{
			Conn:   this.Conn,
			UserId: 0,
		}
		err = up.ServerProcessChangeStatus(mes)
	default:
		fmt.Println("选项有误")
	}
	return
}

func (this *Processor) Process2() (err error) {
	//和客户端保持连接，每当客户端发送消息就进行处理
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg()，返回Message，Err
		//创建一个Transfer实例完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
			Buf:  [8096]byte{},
		}
		//从连接中读取数据
		mes, err1 := tf.ReadPkg()
		//var loginMes message.LoginMes
		//err:=json.Unmarshal([]byte(mes.Data),loginMes)

		if err1 != nil {
			fmt.Println("Process2 ReadPkg() err = ", err)
			return
		}
		if err1 != nil {
			//io.EOF说明客户端断开，数据读取结束
			if err1 == io.EOF {
				fmt.Println("客户端退出，服务器端也退出")
				err = err1
				return
			} else {
				return
			}
		}
		//处理从连接中读取到的数据
		err1 = this.ServerProcessMes(&mes)
		err = err1
		if err1 != nil {
			return
		}
	}
}
