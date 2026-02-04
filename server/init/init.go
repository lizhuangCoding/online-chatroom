package init

import (
	"online-chatroom-2/server/pkg"
)

// 总初始化
func init() {
	// 初始化 redis pool
	RedisPoolInit(pkg.REDISADDRESS, pkg.MAXIDLE, pkg.MAXACTIVE, pkg.IDLETIMEOUT)

	// 初始化 redis 查询全局变量
	RedisRepInit()

	// 全局在线用户变量
	UserMapInit()
}
