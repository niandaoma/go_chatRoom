package dao

import (
	"demo/chatRoom/common/message"
	"demo/chatRoom/server/model"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//定义一个UserDao 结构体 完成对User结构体的各种操作

//在服务器启动之后就初始化一个userDao实例
//它做成全局变量，和redis操作的时候，可以直接用

var (
	MyUserDao *UserDao
)

type UserDao struct {
	pool *redis.Pool
}

//工厂模式得到一个UserDao实例
//连接池在其他地方初始化
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{pool: pool}
	return
}

//增删改查

//1.根据用户 id 返回 一个User实例+err

func (this *UserDao) getUserById(conn redis.Conn, id int) (user model.User, err error) {
	//在 redis 查询这个用户
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		//如果没有该id
		if err == redis.ErrNil { //err == redis.ErrNil 表示没有找到对应id
			err = model.ERROR_USER_NOTEXISTS
		}
		return
	}

	//现在的 res 还需要反序列化成User实例
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	return
}

//完成登陆的校验  Login
//1. Login完成对用户的验证
//2. 如果 id 和 Pwd 正确，返回一个user实例
//3. 如果有错，返回一个错误信息

func (this *UserDao) Login(userId int, userPwd string) (user model.User, err error) {
	//在 dao 连接池取一根链接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}
	//现在id已经没有问题了

	//从redis获取的密码与用户输入的密码不对
	if user.UserPwd != userPwd {
		err = model.ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	//在 dao 连接池取一根链接
	conn := this.pool.Get()
	defer conn.Close()
	//向redis中查询是否有该用户存在
	_, err = this.getUserById(conn, user.UserId)
	//如果没有错误，说明有该用户存在，返回错误
	if err == nil {
		err = model.ERROR_USER_EXISTS
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		return
	}
	return
}
