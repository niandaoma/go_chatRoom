package process2

import (
	"demo/chatRoom/common/message"
	"demo/chatRoom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
	//conn net.Conn
}

//转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的map,得到所有在线的用户的id
	var smsMes message.SmsMes
	//反序列化客户端传来的信息
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Unmarshal() Error = ", err)
		return
	}
	//遍历所有的在线用户列表，发送信息
	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(mes, up.Conn)
	}
}

//发送给信息给在线用户
func (this *SmsProcess) SendMesToEachOnlineUser(mes *message.Message, conn net.Conn) {
	//序列化信息
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("SendMesToEachOnlineUser json.Marshal(mes) Error = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendMesToEachOnlineUser  tf.WritePkg(data) Error = ", err)
		return
	}
}

//点对点发送消息
func (this *SmsProcess) SendP2PMes(mes *message.Message) {
	//遍历所有在线用户列表，找到我们要发送的id
	var spms message.SmsP2PMes
	err := json.Unmarshal([]byte(mes.Data), &spms)
	if err != nil {
		fmt.Println("SendP2PMes json.Unmarshal() Error = ", err)
		return
	}
	//遍历所有在线用户id，发现要通信的id发送信息
	for id, up := range userMgr.onlineUsers {
		if id == spms.SendId {
			this.SendMesToEachOnlineUser(mes, up.Conn)
		}
	}
}
