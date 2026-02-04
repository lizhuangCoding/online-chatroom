package logics

import (
	"context"
	"online-chatroom-2/common/message"
)

// MessageHandler 接口定义了一个DO方法，所有具体的消息处理器（如LoginHandler、RankHandler等）都实现该接口。
type MessageHandler interface {
	DO(ctx context.Context, mes *message.Message) (context.Context, error)
}

// 我这里为什么要返回一个context.Context？因为我在聊天的时候需要用到 UserID，我通过context来传递UserID，也可以把UserID存储到 LoginHandler 结构体里

// MethodController 全局变量。
// （后续可以把 map[message.TY]MessageHandler 类型改为 map[message.TY][]MessageHandler，一个消息类型对应一个切片。比如某个消息类型有多个业务逻辑）
// var MethodController map[message.TY]MessageHandler

// var (
// 	VLoginHandler = &LoginHandler{}
// 	VRankHandler  = &RankHandler{}
// 	VChatHandler  = &ChatHandler{}
// 	VHeartHandler = &HeartHandler{}
// )

// InitMethodController 初始化 MethodController，将消息类型（Login、Rank等）映射到对应的处理器（在 MethodController 中注册处理器）
// func InitMethodController(conn net.Conn) {
// 	MethodController = make(map[message.TY]MessageHandler)
// 	// MethodController[message.Login] = VLoginHandler
// 	// MethodController[message.Rank] = VRankHandler
// 	// MethodController[message.Chat] = VChatHandler
// 	// MethodController[message.Heart] = VHeartHandler
// 	MethodController[message.Login] = &LoginHandler{Conn: conn}
// 	MethodController[message.Rank] = &RankHandler{Conn: conn}
// 	MethodController[message.Chat] = &ChatHandler{}
// 	MethodController[message.Heart] = &HeartHandler{Conn: conn}
// }
