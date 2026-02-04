package logics

import (
	"context"
	"encoding/json"
	"fmt"
	"online-chatroom-2/common/message"
	"online-chatroom-2/server/repository"
)

// ChatHandler 具体的聊天处理器
type ChatHandler struct {
	// Conn   net.Conn
	// UserId int
}

func (c *ChatHandler) DO(ctx context.Context, mes *message.Message) (context.Context, error) {
	// 序列化mes
	data, _ := json.Marshal(mes)
	// 发布消息到 Redis
	if err := repository.RedisRep.Publish("chatroom", data); err != nil {
		fmt.Println("发布消息失败:", err)
		return ctx, err
	}

	// 从上下文中获取 userID
	userId, ok := ctx.Value("userID").(int)
	if !ok {
		return ctx, fmt.Errorf("ChatHandler.DO() err = UserId not found in context")
	}

	// 更改用户活跃度
	if err := repository.RedisRep.UpdateActiveById(userId, 1); err != nil {
		return ctx, err
	}
	return ctx, nil
}
