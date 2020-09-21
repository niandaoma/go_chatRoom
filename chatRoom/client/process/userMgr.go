//用来维护客户端的 OnlineUsers Map

package process

import (
	"demo/chatRoom/client/model"
	"demo/chatRoom/common/message"
	"fmt"
)

var onlineUsers map[int]*message.User = make(map[int]*message.User, 1024)
var CurUser model.CurUser //在用户登陆后，完成CurUser初始化

//在 客户端 显示当前 在线的用户
func outputOnlineUser() {
	//遍历onlineUsers
	fmt.Println("当前在线用户列表:")
	for id, _ := range onlineUsers {
		fmt.Println("用户id\t", id)
	}
}

//编写一个方法处理返回的信息
func upDataUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	switch notifyUserStatusMes.UserStatus {
	case message.UserOnline:
		user, ok := onlineUsers[notifyUserStatusMes.UserId]
		if !ok { //原来没有
			user = &message.User{
				UserId: notifyUserStatusMes.UserId,
			}
		}
		user.UserStatus = notifyUserStatusMes.UserStatus
		onlineUsers[notifyUserStatusMes.UserId] = user
	case message.UserOffline:
		delete(onlineUsers, notifyUserStatusMes.UserId)
	}
	outputOnlineUser()
}
