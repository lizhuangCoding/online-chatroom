package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

type TY int

const (
	Login TY = iota // 登录
	Rank            // 排行
	Chat            // 聊天
	Heart           // 心跳
)

type Message struct {
	Type  TY     `json:"type"`  // 消息类型
	Data  string `json:"data"`  // 消息内容
	Error string `json:"error"` // 错误信息
}

// MessageHandler 处理器接口
type MessageHandler interface {
	DO(mes Message) error
}

// MethodController 存储消息类型与处理器的映射
var MethodController = make(map[TY]MessageHandler)

// LoginHandler 具体的登录处理器
type LoginHandler struct{}

func (h *LoginHandler) DO(mes Message) error {
	// 处理登录逻辑
	fmt.Println("Login Handler: ", mes.Data)
	// 这里可以根据具体的业务逻辑做处理
	return nil
}

// RankHandler 具体的排行榜处理器
type RankHandler struct{}

func (h *RankHandler) DO(mes Message) error {
	// 处理排行榜逻辑
	fmt.Println("Rank Handler: ", mes.Data)
	// 这里可以根据具体的业务逻辑做处理
	return nil
}

// ChatHandler 具体的聊天处理器
type ChatHandler struct{}

func (h *ChatHandler) DO(mes Message) error {
	// 处理聊天逻辑
	fmt.Println("Chat Handler: ", mes.Data)
	// 这里可以根据具体的业务逻辑做处理
	return nil
}

// HeartHandler 具体的心跳处理器
type HeartHandler struct{}

func (h *HeartHandler) DO(mes Message) error {
	// 处理心跳检测逻辑
	fmt.Println("Heart Handler: ", mes.Data)
	// 这里可以根据具体的业务逻辑做处理
	return nil
}

// 初始化方法映射
func init() {
	MethodController[Login] = &LoginHandler{}
	MethodController[Rank] = &RankHandler{}
	MethodController[Chat] = &ChatHandler{}
	MethodController[Heart] = &HeartHandler{}
}

func Process(conn net.Conn) error {
	packetIO := PacketIO{Conn: conn}

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

		// 打印接收到的消息
		demo, _ := json.Marshal(mes)
		fmt.Println("服务端接收客户端的信息为：", string(demo))

		// 通过消息类型获取相应的处理器
		handler, exists := MethodController[mes.Type]
		if !exists {
			return errors.New("wrong message type")
		}

		// 执行对应的处理器
		if err = handler.DO(mes); err != nil {
			return err
		}

	}
}

// PacketIO 处理数据包的输入输出
type PacketIO struct {
	Conn net.Conn   // 服务端与客户端链接
	Buf  [8096]byte // 缓冲
}

// WritePkg 写包。先发送信息的长度，然后发送真正的信息。
func (p *PacketIO) WritePkg(data []byte) error {
	pkgLen := uint32(len(data))
	// 把 pkgLen 存储到 Buf 的前四位中
	binary.BigEndian.PutUint32(p.Buf[:4], pkgLen)

	// 先发送长度
	n, err := p.Conn.Write(p.Buf[:4])
	if n != 4 {
		return errors.New("发送信息错误")
	} else if err != nil {
		return err
	}

	// 发送真正的信息
	n, err = p.Conn.Write(data)
	// 如果理论上发送信息的长度和实际上发送的长度不相等
	if n != int(pkgLen) {
		return errors.New("发送信息错误")
	} else if err != nil {
		return err
	}
	return nil
}

// ReadPkg 读包。先读取信息的长度，然后读取真正的信息。
func (p *PacketIO) ReadPkg() (Message, error) {
	// 读取信息的长度
	_, err := p.Conn.Read(p.Buf[:4])
	if err != nil {
		return Message{}, err
	}

	// 转换为数字
	pkgLen := binary.BigEndian.Uint32(p.Buf[:4])

	// 读取真正的信息
	n, err := p.Conn.Read(p.Buf[:pkgLen])
	if n != int(pkgLen) {
		return Message{}, errors.New("读取信息错误")
	} else if err != nil {
		return Message{}, err
	}

	// 反序列化
	var mes Message
	if err = json.Unmarshal(p.Buf[:pkgLen], &mes); err != nil {
		return Message{}, err
	}
	return mes, nil
}
