package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"online-chatroom-2/common/message"
	"online-chatroom-2/common/model"
	"online-chatroom-2/common/utils"
	"os"
	"strconv"
	"time"
)

// 全局变量
var (
	UserID   int
	packetIO utils.PacketIO
)

// InitPacketIO 初始化 packetIO
func InitPacketIO(conn net.Conn) {
	packetIO = utils.PacketIO{Conn: conn}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8890")
	if err != nil {
		fmt.Println("链接服务器失败 err = ", err)
		return
	}
	defer conn.Close()
	fmt.Println("链接成功...")

	// 初始化 packetIO
	InitPacketIO(conn)

	// 开协程，心跳检测
	go HeartBeat()

	// 登录
	Login()

	// 开协程，读消息
	go func() {
		for {
			Read()
		}
	}()

	// 写消息
	fmt.Print("请输入内容(exit退出，rank排行榜):")
	for {
		Write()
	}
}

// Login 登录（输入用户ID）
func Login() {
	fmt.Print("请输入用户ID:")
	if _, err := fmt.Scanf("%d\n", &UserID); err != nil {
		fmt.Println("输入用户ID错误，err = ", err)
		os.Exit(0)
	}
	if UserID <= 0 {
		fmt.Println("ID应为正数")
		os.Exit(0)
	}

	// 给服务端发送登录消息（包含用户ID）
	mes := message.Message{
		Type: message.Login,
		Data: strconv.Itoa(UserID),
	}
	data, err := json.Marshal(mes)
	if err != nil {
		return
	}
	// fmt.Println("客户端发送的信息为：", mes, data)
	if err = packetIO.WritePkg(data); err != nil {
		fmt.Println("客户端发送信息错误，err = ", err)
		return
	}

	// 读取服务端的消息，判断是否可以成功登录
	mes, err = packetIO.ReadPkg()
	if err == io.EOF {
		fmt.Println("服务端主动关闭了链接，客户端也退出...")
		os.Exit(0)
	} else if err != nil {
		fmt.Println(mes, err)
		return
	}
	// fmt.Println("客户端接收的信息为：", mes)

	if mes.Type == message.Login && mes.Error == "" { // 成功登录
		// loop = false
		fmt.Println("登录成功")
	} else {
		fmt.Println("登录失败 ", mes.Error)
		os.Exit(0)
	}
}

// Read 读
func Read() {
	mes, err := packetIO.ReadPkg()
	if err == io.EOF {
		fmt.Println("服务端主动关闭了链接，客户端也退出...")
		os.Exit(0)
	} else if err != nil {
		fmt.Println("客户端读取信息错误，err = ", err)
		return
	}

	switch mes.Type {
	case message.Rank: // 排行榜
		// 反序列化
		var rankList = make([]model.Rank, 0)
		if err := json.Unmarshal([]byte(mes.Data), &rankList); err != nil {
			fmt.Println("反序列化错误，err = ", err)
			return
		}

		fmt.Println("===== 排行榜 =====")
		fmt.Println(" ID  : 活跃度")
		for _, v := range rankList {
			fmt.Printf("%-4d : %-10d\n", v.ID, v.Activity)
		}

	case message.Chat: // 聊天
		fmt.Println(mes.Data)

	case message.Heart: // 心跳
		// fmt.Println("客户端收到心跳检测：", mes)

	default:
		fmt.Println("无法识别的类型")
	}
}

// Write 写内容
func Write() {
	// 从键盘读取输入
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n') // 直到遇到换行符 \n，才将读取到的内容赋值给变量input
	input = input[:len(input)-1]        // 删除字符串最后面的换行符
	// fmt.Println("id 为", UserID, "的用户输入了：", input)

	mes := message.Message{}
	if input == "rank" { // 排行榜
		mes.Type = message.Rank
	} else if input == "exit" || input == "quit" { // 退出
		os.Exit(0)
	} else { // 聊天
		mes.Type = message.Chat
		mes.Data += fmt.Sprintf("用户 %d 发送：", UserID)
	}

	mes.Data += input
	// 序列化
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化错误，err = ", err)
		return
	}

	// 写包
	if err = packetIO.WritePkg(data); err != nil {
		fmt.Println("客户端发送信息错误，err = ", err)
		return
	}
}

// HeartBeat 心跳检测，每隔5s发送一次PING
func HeartBeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // 发送数据
			mes := message.Message{
				Type: message.Heart,
				Data: "PING",
			}
			data, _ := json.Marshal(mes)
			if err := packetIO.WritePkg(data); err != nil {
				fmt.Println("客户端发送信息错误，err = ", err)
				return
			}
			// fmt.Println("客户端发送PING...")
		}
	}

}
