package controller

import (
	"fmt"
	"net"
	"online-chatroom-2/server/pkg"
	"time"
)

// Controller 控制器
func Controller(conn net.Conn) {
	defer func() {
		// 删除map中的用户
		pkg.UserMap.DelByConn(conn)
		fmt.Println("从map中删除了一个用户...")

		// 关闭链接
		conn.Close()
	}()

	// 关于心跳检测：
	// 客户端写一个周期定时器，每隔5s发送一次PING。
	// 服务端设置：conn.SetReadDeadline(time.Now().Add(6 * time.Second))，每隔6s接收一次，接收到 PING 后(然后给客户端发送了PONG，其实客户端那里应该也要设置一下如果XXs内没有接收到PONG就退出，但是我没有设置)， 重置conn.SetReadDeadline(time.Now().Add(6 * time.Second))，一直这样循环。
	// 我写的是只要服务端接收到客户端的 PING，就重置超时时间，但是其实我们也可以只要接收到任意数据就重置超时时间。
	// （如果超出了超时时间，服务端还是没有接收到客户端的数据的话，conn会自动断开链接，会报：i/o timeout）
	if err := conn.SetReadDeadline(time.Now().Add(6 * time.Second)); err != nil {
		fmt.Println("设置读取超时失败:", err)
		return
	}

	// 具体的控制器内容
	if err := Process(conn); err != nil {
		fmt.Println("Controller Process() err = ", err)
		return
	}
}
