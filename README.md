# Online Chatroom

一个基于 Go 语言开发的在线聊天室应用，支持多用户实时聊天、排行榜等功能。

## 功能特性

- 用户登录（通过用户 ID）
- 群聊功能 - 所有在线用户可以互相聊天
- 用户排行榜 - 基于活跃度排序
- 心跳检测 - 客户端每 5 秒发送一次心跳保持连接
- Redis 数据持久化

## 项目结构

```
online-chatroom-2/
├── client/          # 客户端代码
│   └── cmd/
│       └── main.go  # 客户端入口
├── server/          # 服务端代码
│   ├── cmd/
│   │   └── main.go  # 服务端入口
│   ├── controller/  # 控制器
│   ├── logics/      # 业务逻辑
│   ├── init/        # 初始化
│   ├── repository/  # 数据访问层
│   └── pkg/         # 工具包
├── common/          # 公共代码
│   ├── message/     # 消息定义
│   ├── model/       # 数据模型
│   └── utils/       # 工具函数
└── demo/            # 示例代码
```

## 快速开始

### 环境要求

- Go 1.16+
- Redis（用于数据持久化）

### 安装

1. 克隆项目

```bash
git clone https://github.com/lizhuangCoding/online-chatroom.git
cd online-chatroom-2
```

2. 安装依赖

```bash
go mod download
```

3. 确保 Redis 服务已启动

### 运行

1. 启动服务端（监听端口 8890）

```bash
go run server/cmd/main.go
```

2. 启动客户端（新开终端）

```bash
go run client/cmd/main.go
```

### 使用说明

1. 客户端启动后，输入用户 ID 进行登录
2. 登录成功后，输入消息即可发送到聊天室
3. 输入 `rank` 查看用户排行榜
4. 输入 `exit` 或 `quit` 退出

## 消息类型

| 类型 | 说明 |
|------|------|
| Login | 登录消息 |
| Chat | 聊天消息 |
| Rank | 排行榜查询 |
| Heart | 心跳检测 |

## 技术栈

- **语言**: Go
- **网络通信**: TCP
- **数据存储**: Redis
- **数据格式**: JSON

## 开源协议

MIT License