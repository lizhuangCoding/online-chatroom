package message

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
