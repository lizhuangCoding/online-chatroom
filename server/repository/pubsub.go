package repository

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"online-chatroom-2/server/pkg"
)

// 聊天消息的发布订阅

// BroadcastMessages 从这里开启广播消息
func BroadcastMessages() {
	msgChan, err := RedisRep.Subscribe("chatroom")
	if err != nil {
		fmt.Println("订阅频道失败:", err)
		return
	}

	// 遍历每一个Client，并给每一个Client发信息
	for msg := range msgChan {
		pkg.UserMap.RangeClientSendMes([]byte(msg))
	}
}

// Publish 发布消息到 Redis
func (r *RedisRepository) Publish(channel string, message []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PUBLISH", channel, message)
	return err
}

// Subscribe 订阅 Redis 通道
func (r *RedisRepository) Subscribe(channel string) (chan string, error) {
	conn := r.pool.Get()
	// defer conn.Close() // 注意不要在这里关闭，否则下面的协程sub.Receive()的打印结果是 redigo: connection closed

	sub := redis.PubSubConn{Conn: conn}
	if err := sub.Subscribe(channel); err != nil {
		return nil, err
	}

	msgChan := make(chan string)
	go func(conn redis.Conn) {
		defer func() {
			conn.Close()
			close(msgChan)
		}()

		for {
			switch v := sub.Receive().(type) {
			case redis.Message:
				msgChan <- string(v.Data)

			case redis.Subscription:
				continue

			default:
				fmt.Println("无法识别的 sub.Receive() : ", sub.Receive())
				return
			}
		}
	}(conn)
	return msgChan, nil
}
