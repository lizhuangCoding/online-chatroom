package logics

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"online-chatroom-2/common/message"
	"online-chatroom-2/common/model"
	"online-chatroom-2/common/utils"
	"online-chatroom-2/server/pkg"
	"online-chatroom-2/server/repository"
	"strconv"
)

// LoginHandler 具体的登录处理器
type LoginHandler struct {
	Conn net.Conn
	// UserId int
}

func (l *LoginHandler) DO(ctx context.Context, mes *message.Message) (context.Context, error) {
	// 判断登录数据是否合法
	code, err := IsValidLogin(mes) // 这里不能return错误

	// 回送消息
	resMes := message.Message{Type: message.Login}
	// 如果数据没问题，将这个用户加入到 map 中，并且把这个id赋值给 UserContext 对象，否则发送错误
	if code == 200 {
		id, _ := strconv.Atoi(mes.Data) // 数据不合法
		pkg.UserMap.Set(id, l.Conn)
		// l.UserId = id

		// 注意：这里一定要赋值（感觉应该可以用 context 包来传递 userId 参数，比这种更专业）
		// VChatHandler.UserId = id

		// 使用现有的 ctx 作为基础上下文
		ctx = context.WithValue(ctx, "userID", id)

		// 使用 context.Background() 作为基础上下文。context.Background() 是 Go 中的一个顶层（root）上下文，通常用于程序的最外层，通常是程序的起始点。
		// ctx = context.WithValue(context.Background(), "userID", id)

	} else {
		resMes.Error = err.Error()
	}

	// 序列化
	data, err := json.Marshal(resMes)
	if err != nil {
		return ctx, err
	}
	// fmt.Println("发送给客户端的信息为：", string(data))

	packetIO := utils.PacketIO{Conn: l.Conn}
	if err := packetIO.WritePkg(data); err != nil {
		return ctx, err
	}
	return ctx, nil
}

// IsValidLogin 检查登录数据是否合法
func IsValidLogin(mes *message.Message) (int, error) {
	id, err := strconv.Atoi(mes.Data) // 数据不合法
	if err != nil {
		return 400, pkg.ErrorIllegalData
	}

	// 判断该id是否存在
	user, _ := repository.RedisRep.QueryUserById(id)
	if user == nil {
		// 添加该用户
		if err = repository.RedisRep.AddUser(&model.User{ID: id}); err != nil {
			fmt.Println(err)
		}

		// 添加该用户活跃度
		if err = repository.RedisRep.UpdateActiveById(id, 0); err != nil {
			fmt.Println(err)
		}
	}

	// 判断该用户是否已经在线
	_, ok := pkg.UserMap.Get(id)
	if ok {
		return 400, pkg.ErrorIDOnline
	}
	return 200, nil
}
