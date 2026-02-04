package pkg

import (
	"fmt"
	"net"
	"online-chatroom-2/common/utils"
	"sync"
)

// UserMap 全局在线用户变量
var UserMap *OnlineMap

// OnlineMap 存放在线用户(并发安全)
type OnlineMap struct {
	mu      sync.Mutex      // 锁（最好改为读写锁）
	clients map[int]*Client // key:用户id
}

// Client 结构体
type Client struct {
	UserId int
	Conn   net.Conn
}

func NewOnlineMap() *OnlineMap {
	return &OnlineMap{
		clients: map[int]*Client{},
	}
}

// Set 设置数据
func (o *OnlineMap) Set(id int, conn net.Conn) {
	o.mu.Lock()
	defer o.mu.Unlock()

	client := Client{
		UserId: id,
		Conn:   conn,
	}
	o.clients[id] = &client
}

// Get 获取数据
func (o *OnlineMap) Get(id int) (conn net.Conn, ok bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	client, ok := o.clients[id] // 说明该id在线
	if ok {
		conn = client.Conn
	}

	// for _, v := range o.clients {
	// 	fmt.Println("Map 里面的数据为", v.UserId)
	// }

	return conn, ok
}

// DelById 根据 id 删除数据
func (o *OnlineMap) DelById(id int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.clients, id)
}

// DelByConn 根据 Conn 删除数据
func (o *OnlineMap) DelByConn(conn net.Conn) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for i, v := range o.clients {
		if conn == v.Conn {
			delete(o.clients, i)
			break
		}
	}
}

// RangeClientSendMes 遍历每一个Client，并给每一个Client发信息
func (o *OnlineMap) RangeClientSendMes(data []byte) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, v := range o.clients {
		packetIO := utils.PacketIO{Conn: v.Conn}
		if err := packetIO.WritePkg(data); err != nil {
			fmt.Printf("给用户 %d 发送聊天信息错误，err = %v\n", v.UserId, err)
		}
	}
}
