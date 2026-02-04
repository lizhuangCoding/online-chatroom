package main

import (
	"fmt"
	"net"
	"online-chatroom-2/server/controller"
	_ "online-chatroom-2/server/init" // 空导入
	"online-chatroom-2/server/repository"
)

func main() {
	fmt.Println("服务器监听中...")
	listen, err := net.Listen("tcp", "0.0.0.0:8890")
	if err != nil {
		fmt.Println("listen err = ", err)
		return
	}
	defer listen.Close()

	// 开始广播
	go repository.BroadcastMessages()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() err = ", err)
			continue
		}
		fmt.Println("链接成功...")

		// 开启业务逻辑协程
		go func(conn net.Conn) {
			controller.Controller(conn)
		}(conn)
	}
}
