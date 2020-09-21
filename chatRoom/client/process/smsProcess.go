package process

import (
	"demo/chatRoom/client/utils"
	"demo/chatRoom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

type SmsTransferMes struct {
}

//发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个message
	var mes message.Message
	mes.Type = message.SmsMesType

	//2.创建一个SmsMes
	var sms message.SmsMes
	sms.Content = content               //内容
	sms.UserId = CurUser.UserId         //ID
	sms.UserStatus = CurUser.UserStatus //状态
	//3.序列化sms
	data, err := json.Marshal(sms)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(sms) Error = ", err)
		return
	}
	mes.Data = string(data)

	//4.序列化mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(mes) Error = ", err)
		return
	}
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes Error = ", err)
	}
	return
}

//发送点对点聊天
func (this *SmsProcess) SendMesToEachOther(id int, content string) (err error) {
	//1.创建一个message
	var mes message.Message
	mes.Type = message.SmsP2PMesType

	//2.创建一个SmsMes
	var spms message.SmsP2PMes
	spms.Content = content //内容
	spms.SendId = id
	spms.UserId = CurUser.UserId         //ID
	spms.UserStatus = CurUser.UserStatus //状态
	//3.序列化sms
	data, err := json.Marshal(spms)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(sms) Error = ", err)
		return
	}
	mes.Data = string(data)

	//4.序列化mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(mes) Error = ", err)
		return
	}
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
		Buf:  [8096]byte{},
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes Error = ", err)
	}
	return

}
