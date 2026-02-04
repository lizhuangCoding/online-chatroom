package repository

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"online-chatroom-2/common/model"
	"online-chatroom-2/server/pkg"
	"strconv"
)

// RedisRep redis 查询全局变量
var RedisRep *RedisRepository

type RedisRepository struct {
	pool *redis.Pool
}

func NewRedisRepository(p *redis.Pool) *RedisRepository {
	return &RedisRepository{pool: p}
}

// QueryUserById 通过id获取用户详细信息
func (r *RedisRepository) QueryUserById(id int) (*model.User, error) {
	// 从redis链接池中取出一条链接
	conn := r.pool.Get()
	defer conn.Close()

	// 定义用户键
	userKey := fmt.Sprintf("%s:%s:%d", pkg.ProjectName, pkg.UserTable, id)

	// 从Redis中获取用户数据
	data, err := redis.String(conn.Do("GET", userKey))
	if err != nil {
		return nil, err
	}

	// 反序列化
	var user *model.User
	if err = json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return user, nil
}

// AddUser 添加用户（使用 hash 存储用户信息）
func (r *RedisRepository) AddUser(user *model.User) error {
	conn := r.pool.Get()
	defer conn.Close()

	// 定义用户键
	userKey := fmt.Sprintf("%s:%s:%d", pkg.ProjectName, pkg.UserTable, user.ID)

	// 序列化
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if _, err = conn.Do("set", userKey, string(data)); err != nil {
		return err
	}
	return nil
}

// QueryRankList 查询用户的活跃度排行（使用 zset 存储用户的活跃度）
func (r *RedisRepository) QueryRankList(limit, page int) ([]model.Rank, error) {
	conn := r.pool.Get()
	defer conn.Close()

	rankKey := fmt.Sprintf("%s:%s", pkg.ProjectName, pkg.RankTable)

	// 计算起始位置和结束位置
	start := (page - 1) * limit
	end := start + limit - 1

	// 从Redis中获取排名数据
	data, err := redis.Strings(conn.Do("ZREVRANGE", rankKey, start, end, "WITHSCORES"))
	if err != nil {
		return nil, fmt.Errorf("error retrieving rank data from Redis: %w", err)
	}

	// 存储排行信息
	rank := make([]model.Rank, 0, len(data)/2) // 预分配切片容量以提高效率
	for i := 0; i < len(data); i += 2 {
		userId, err := strconv.Atoi(data[i]) // 每对中的第一个是用户ID
		if err != nil {
			return nil, fmt.Errorf("error converting user id: %w", err)
		}
		active, err := strconv.Atoi(data[i+1]) // 活跃度
		if err != nil {
			return nil, fmt.Errorf("error converting active score: %w", err)
		}
		rank = append(rank, model.Rank{
			ID:       userId,
			Activity: active,
		})
	}
	return rank, nil
}

// UpdateActiveById 更改用户活跃度
func (r *RedisRepository) UpdateActiveById(id int, active int) error {
	conn := r.pool.Get()
	defer conn.Close()

	// 定义活跃度排名的键
	rankKey := fmt.Sprintf("%s:%s", pkg.ProjectName, pkg.RankTable)

	_, err := conn.Do("ZINCRBY", rankKey, active, id)
	if err != nil {
		return fmt.Errorf("error updating active score for user %d: %w", id, err)
	}
	return nil
}
