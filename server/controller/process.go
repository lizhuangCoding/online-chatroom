package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"online-chatroom-2/common/message"
	"online-chatroom-2/common/utils"
	"online-chatroom-2/server/logics"
)

// Process 具体的控制器内容
func Process(conn net.Conn) error {
	// 注意 methodController 不能是全局变量，如果是全局变量，那么只能容纳一个用户在线，因为第二个用户上线后会把第一个用户的ctx中存储的userID信息覆盖掉。
	var methodController map[message.TY]logics.MessageHandler
	methodController = make(map[message.TY]logics.MessageHandler)
	methodController[message.Login] = &logics.LoginHandler{Conn: conn}
	methodController[message.Rank] = &logics.RankHandler{Conn: conn}
	methodController[message.Chat] = &logics.ChatHandler{}
	methodController[message.Heart] = &logics.HeartHandler{Conn: conn}

	// 初始化 MethodController，在 MethodController 中注册处理器。
	// logics.InitMethodController(conn)
	// 因为 logics.InitMethodController() 方法中赋值的是指针类型，所以直接修改 logics.VLoginHandler 就可以影响到 MethodController[message.Login] = VLoginHandler 中的 VLoginHandler 了。
	// logics.VLoginHandler.Conn = conn
	// logics.VHeartHandler.Conn = conn
	// logics.VRankHandler.Conn = conn

	// 上下文。等到登录之后，ctx会包含 userID
	// context.Background() 是一个空的上下文，通常用于程序的顶层，表示没有父上下文。它适用于没有父上下文的情况下，作为根上下文来启动协程、HTTP请求等。
	var ctx = context.Background()

	packetIO := utils.PacketIO{Conn: conn}
	// 循环读包
	for {
		mes, err := packetIO.ReadPkg()
		if err != nil {
			if err == io.EOF {
				return errors.New("用户已断开链接")
			} else {
				return errors.New("Process packetIO.ReadPkg() err = " + err.Error())
			}
		}
		// 从上下文中取出 userID。当第一次登录的时候取出的数据为 <nil>。等到该用户登录后，会使用 WithValue 携带 userID 数据，并返回一个 ctx，然后我们用这个 ctx 覆盖掉 context.Background() 就可以从 ctx 中取出 userID 了。
		userID := ctx.Value("userID")

		demo, _ := json.Marshal(mes)
		fmt.Printf("服务端接收客户端用户 %v 的信息为：%v\n", userID, string(demo))

		// 通过 mes.Type 来查找对应的处理器，并调用该处理器的DO方法进行相应的处理。
		// handler, exists := logics.MethodController[mes.Type]
		handler, exists := methodController[mes.Type]
		if !exists {
			return errors.New(fmt.Sprintf("message type %v is not exist\n", mes.Type))
		}
		// 执行对应的处理器
		if ctx, err = handler.DO(ctx, &mes); err != nil {
			return err
		}
	}
}
