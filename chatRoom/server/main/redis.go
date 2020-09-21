package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var Pool *redis.Pool

//创建链接池
func InitPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	Pool = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", address)
		},
		TestOnBorrow:    nil,
		MaxIdle:         maxIdle,
		MaxActive:       maxActive,
		IdleTimeout:     idleTimeout,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}
