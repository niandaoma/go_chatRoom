package main

import (
	"demo/chatRoom/client/process"
	"fmt"
	"os"
)

var userId int
var userPwd string
var userName string

func main() {
	var key int //接受用户的选择
	//var loop = true //loop判断是否退出菜单循环
	for {
		fmt.Println("----------欢迎登陆多人聊天系统---------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3)")

		fmt.Scanln(&key)
		up := process.UserProcess{}
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanln(&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n",&userPwd)
			//loop=false
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户的id")
			fmt.Scanln(&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf(" %s\n", &userPwd)
			fmt.Println("请输入用户的姓名")
			fmt.Scanf(" %s\n", &userName)
			up.Register(userId, userPwd, userName)
			//loop=false
		case 3:
			fmt.Println("退出系统")
			//loop=false
			os.Exit(0)
		default:
			fmt.Println("你的输入有误请重新输入")
		}

	}
}
