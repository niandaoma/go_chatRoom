package process2

import (
	"fmt"
)

//定成全局变量
var (
	userMgr *UserMgr
)

//在线列表
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//增
func (this *UserMgr) AddOnlineUserProcess(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

//删
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

//返回所有在线的用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

//根据id返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userId] //待检测的从map中取出一个值
	if !ok {                           //!ok说明userId对应的用户当前不在线
		err = fmt.Errorf("用户不在线或不存在")
		return
	}
	return
}
