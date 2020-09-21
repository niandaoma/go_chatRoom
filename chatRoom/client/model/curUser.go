package model

import (
	"demo/chatRoom/common/message"
	"net"
)

//在客户端很多地方都会使用，这里做一个 全局CurUser
type CurUser struct {
	Conn net.Conn
	message.User
}
