package logics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"online-chatroom-2/common/message"
	"online-chatroom-2/common/utils"
	"time"
)

// HeartHandler 具体的心跳处理器
type HeartHandler struct {
	Conn net.Conn
	// UserId int
}

func (h *HeartHandler) DO(ctx context.Context, mes *message.Message) (context.Context, error) {
	// 收到 PONG 后重置心跳检测
	if mes.Data == "PING" {
		if err := h.Conn.SetReadDeadline(time.Now().Add(6 * time.Second)); err != nil {
			fmt.Println("设置读取超时失败:", err)
			return ctx, err
		}
		// fmt.Println("服务端重置心跳时间...")

		// 给客户端发送 PONG
		resMes := message.Message{
			Type: message.Heart,
			Data: "PONG",
		}
		data, _ := json.Marshal(resMes)

		packetIO := utils.PacketIO{Conn: h.Conn}
		if err := packetIO.WritePkg(data); err != nil {
			return ctx, errors.New(fmt.Sprintln("Error sending PONG:", err))
		}
	}
	return ctx, nil
}
