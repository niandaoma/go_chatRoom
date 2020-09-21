package process

import (
	"demo/chatRoom/client/utils"
	"demo/chatRoom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

//显示登陆成功后的界面。。
func ShowLoginMenu() {
	for {

		var key int
		var content string
		var id int

		//总会使用到SmsProcess实例，定义swtich在外部节约空间
		smsProcess := &SmsProcess{}
		up := &UserProcess{}
		fmt.Println("-----恭喜***登陆成功-----")
		fmt.Println("----1	显示在线用户列表----")
		fmt.Println("----2	发送消息----")
		fmt.Println("----3	单独发送消息----")
		fmt.Println("----4	信息列表----")
		fmt.Println("----5	退出系统----")
		fmt.Println("----请选择(1-4)----")
		fmt.Scanf("%d\n", &key)
		//fmt.Scanf(" %c")
		switch key {
		case 1:
			outputOnlineUser()
		case 2:

			fmt.Println("发送消息")
			fmt.Println("请输入你想广播的消息")
			fmt.Scanf("%s\n", &content)
			time.Sleep(time.Microsecond*300)
			err := smsProcess.SendGroupMes(content)
			if err != nil {
				fmt.Println("smsProcess.SendGroupMes(content) Error = ", err)
				return
			}
		case 3:
			fmt.Println("点对点聊天")
			fmt.Println("请输入你要发送的id:")
			fmt.Scanf("%d\n", &id)
			time.Sleep(time.Microsecond*300)
			fmt.Println("请输入你想广播的消息")
			time.Sleep(time.Microsecond*300)
			fmt.Scanf("%s\n", &content)
			err := smsProcess.SendMesToEachOther(id, content)
			if err != nil {
				fmt.Println("smsProcess.SendMesToEachOther() Error = ", err)
				return
			}
		case 4:
			fmt.Println("查看信息列表")
		case 5:
			fmt.Println("退出系统")
			up.ClientExit()
			os.Exit(0)
		default:
			fmt.Println("你输入的选项不正确。。")
		}
	}
}

func ShowRegisterMenu() {
	for {
		var key int

		fmt.Println("-----恭喜***注册成功-----")
		fmt.Println("----1	显示在线用户列表----")
		fmt.Println("----2	发送消息----")
		fmt.Println("----3	信息列表----")
		fmt.Println("----4	退出系统----")
		fmt.Println("----请选择(1-4)----")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("显示在线用户列表")
		case 2:
			fmt.Println("发送列表")
		case 3:
			fmt.Println("查看信息列表")
		case 4:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("你输入的选项不正确。。")
		}
	}
}

//和服务器端保持通讯
func serverProcessMes(conn net.Conn) {
	//创建一个Transfer实例，不停的读取服务器发送的消息
	tf := utils.Transfer{
		Conn: conn,
		Buf:  [8096]byte{},
	}
	for {
		mes, err := tf.ReadPkg()
		//如果读取出错，比如服务器端或客户端退出

		if err != nil {
			fmt.Println("tf.ReadPkg err = ", err)
			return
		}
		//如果读取到消息，，又是下一步处理逻辑
		switch mes.Type {
		//有人的状态发生改变，比如上线
		case message.NotifyUserStatusMesType:
			//1.取出NotifyUserStatusMes
			//2.把这个人加入客户端Online map里面
			//
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			upDataUserStatus(&notifyUserStatusMes)
		case message.SmsMesType: //有人群发消息
			outputGroupMes(&mes)
		case message.SmsP2PMesType:
			outputP2PMes(&mes)
		default:
			fmt.Println("返回了一个暂时还不能识别的消息")
		}

	}
}
