package init

import (
	"github.com/garyburd/redigo/redis"
	"online-chatroom-2/server/pkg"
	"online-chatroom-2/server/repository"
	"time"
)

// RedisPoolInit 初始化 redis pool
func RedisPoolInit(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pkg.RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空闲链接数
		MaxActive:   maxActive,   // 表示和数据库的最大链接数， 0 表示没有限制
		IdleTimeout: idleTimeout, // 最大空闲时间
		Dial: func() (redis.Conn, error) { // 初始化链接的代码， 链接哪个ip的redis
			return redis.Dial("tcp", address)
		},
	}
}

// RedisRepInit 初始化 redis 查询全局变量
func RedisRepInit() {
	repository.RedisRep = repository.NewRedisRepository(pkg.RedisPool)
}
