package main

import (
	"demo/chatRoom/server/dao"
	"fmt"
	"net"
	"time"
)

func init() {
	//服务器一启动，我们就初始化 redis 链接池
	InitPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

//创建客户端与服务器建立连接的进程
func process(conn net.Conn) {
	defer conn.Close()
	//获取客户端ip地址
	addr:=conn.RemoteAddr()
	fmt.Println("客户端ip地址：",addr)
	//读客户端发送的信息
	//这里调用总控
	//创建一个process结构体
	processor := &Processor{Conn: conn}
	//调用processor的方法，进行相应服务
	err := processor.Process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯协程错误=", err)
		return
	}
}

func initUserDao() {
	//初始化顺序n'n'n
	//先initPool，再initUserDao
	dao.MyUserDao = dao.NewUserDao(Pool)
}

func main() {
	fmt.Println("服务器在8889端口监听。。。。")
	listen, err := net.Listen("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.listen err = ", err)
		return
	}
	//关闭tcp连接
	defer listen.Close()
	//服务器端一直监听端口，等待客户端连接
	for {
		fmt.Println("等待客户端链接服务器")
		//获得链接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err = ", err)
			return
		}
		//一旦连接成功，就启动一个协程和客户端保持通讯
		go process(conn)
	}
}
