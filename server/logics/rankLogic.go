package logics

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"online-chatroom-2/common/message"
	"online-chatroom-2/common/utils"
	"online-chatroom-2/server/repository"
)

// RankHandler 具体的排行榜处理器
type RankHandler struct {
	Conn net.Conn
	// UserId int
}

func (r *RankHandler) DO(ctx context.Context, mes *message.Message) (context.Context, error) {
	limit := 10
	page := 1
	rankList, err := repository.RedisRep.QueryRankList(limit, page)
	if err != nil {
		return ctx, err
	}

	// 回送消息
	resMes := message.Message{Type: message.Rank}
	// 将排行数据序列化
	data, _ := json.Marshal(rankList)
	resMes.Data = string(data)
	data, _ = json.Marshal(resMes)

	fmt.Println("发送的信息为：", string(data))

	packetIO := utils.PacketIO{Conn: r.Conn}
	if err = packetIO.WritePkg(data); err != nil {
		return ctx, err
	}
	return ctx, nil
}
