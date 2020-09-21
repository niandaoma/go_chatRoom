package process

import (
	"demo/chatRoom/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) {
	var sms message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &sms)
	if err != nil {
		return
	}

	//显示信息
	fmt.Printf("用户%d群发消息:%s\n", sms.UserId, sms.Content)
}

func outputP2PMes(mes *message.Message) {
	var spms message.SmsP2PMes
	err := json.Unmarshal([]byte(mes.Data), &spms)
	if err != nil {
		return
	}

	//显示信息
	fmt.Printf("用户%d向你发送消息:%s\n", spms.UserId, spms.Content)
}
