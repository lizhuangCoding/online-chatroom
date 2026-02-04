package init

import "online-chatroom-2/server/pkg"

// UserMapInit 初始化全局在线用户变量
func UserMapInit() {
	pkg.UserMap = pkg.NewOnlineMap()
}
