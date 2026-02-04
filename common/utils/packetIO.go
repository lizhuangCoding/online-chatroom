package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"online-chatroom-2/common/message"
)

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
func (p *PacketIO) ReadPkg() (message.Message, error) {
	// 读取信息的长度
	_, err := p.Conn.Read(p.Buf[:4])
	if err != nil {
		return message.Message{}, err
	}

	// 转换为数字
	pkgLen := binary.BigEndian.Uint32(p.Buf[:4])

	// 读取真正的信息
	n, err := p.Conn.Read(p.Buf[:pkgLen])
	if n != int(pkgLen) {
		return message.Message{}, errors.New("读取信息错误")
	} else if err != nil {
		return message.Message{}, err
	}

	// 反序列化
	var mes message.Message
	if err = json.Unmarshal(p.Buf[:pkgLen], &mes); err != nil {
		return message.Message{}, err
	}
	return mes, nil
}
