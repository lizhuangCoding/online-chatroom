package pkg

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// RedisPool redis 池
var RedisPool *redis.Pool

// 关于 redis pool 的常量
const (
	REDISADDRESS = "localhost:6379"
	MAXIDLE      = 16
	MAXACTIVE    = 0
	IDLETIMEOUT  = 300 * time.Second
)

// redis 数据库的表名
const (
	ProjectName = "online-chatroom" // 项目名称
	UserTable   = "users"           // 用户表名称
	RankTable   = "ranks"           // 用户表名称
)
