package ptdb

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	//pt lib
	"github.com/chenhnu/go-lib/ptlog"
)

type redisPoolInfo struct {
	//Host redis ip 地址
	Host string
	//Port redis 端口
	Port int
	//Auth redis 密码
	Auth string
	//MaxIdle 连接池的最大连接数
	MaxIdle int
	//TimeoutSec 超时时间（单位：秒），该参数需要小于redis数据库的超时时间
	TimeoutSec int
}

func NewRedisPoolInfo(host string, port int, auth string, maxIdle int, timeoutSec int) *redisPoolInfo {
	return &redisPoolInfo{
		Host:       host,
		Port:       port,
		Auth:       auth,
		MaxIdle:    maxIdle,
		TimeoutSec: timeoutSec,
	}
}

type ptRedisPool struct {
	pool *redis.Pool
}

func NewRedisPool(info *redisPoolInfo) *ptRedisPool {
	return &ptRedisPool{
		pool: &redis.Pool{
			MaxIdle:     info.MaxIdle,
			IdleTimeout: time.Duration(info.TimeoutSec) * time.Second,
			Dial: func() (redis.Conn, error) {
				if !isHost(info.Host) {
					return nil, errors.New("host illegal")
				}
				conn, e := redis.Dial("tcp", info.Host+":"+strconv.Itoa(info.Port))
				if e != nil {
					return nil, e
				}
				reply, e := conn.Do(AUTH, info.Auth)
				if e != nil {
					return nil, e
				}
				if reply != "OK" {
					return nil, errors.New("password error")
				}
				return conn, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}

//将do和send函数封装一层，后期可以以在此处选择要使用的连接（redis pool）
//后期可以在此处添加连接中断的重新恢复（recovery connect）
func (redisPool *ptRedisPool) doCmd(cmd string, args ...interface{}) (interface{}, error) {
	conn := redisPool.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	return conn.Do(cmd, args...)
}

//Send()函数只是将命令写进buffer里，并不会执行，而Flush()就可以执行
//可以在使用事务的时候使用，一个事务一次提交
func (redisPool *ptRedisPool) sendCmd(cmd string, args ...interface{}) error {
	conn := redisPool.pool.Get()
	e := conn.Send(cmd, args...)
	e = conn.Flush()
	return e
}

func (redisPool *ptRedisPool) Keys() []string {
	var keys []string
	rp, e := redisPool.doCmd(KEYS, "*")
	if e != nil {
		ptlog.Error(e)
		return nil
	}
	if rp != nil {
		for _, v := range rp.([]interface{}) {
			keys = append(keys, string(v.([]uint8)))
		}
	}
	return keys
}

func (redisPool *ptRedisPool) delKeys(args ...interface{}) error {
	_, e := redisPool.doCmd(DEL, args...)
	return e
}

func (redisPool *ptRedisPool) ExistsKey(key string) bool {
	rp, e := redisPool.doCmd(EXISTS, key)
	if e != nil {
		ptlog.Error(e)
		return false
	}
	if rp.(int64) > 0 {
		return true
	}
	return false
}

//敏感操作，会删除所有的缓存，建议隐藏
func (redisPool *ptRedisPool) FlushAll() {
	_, e := redisPool.doCmd(FlushALL)
	if e != nil {
		ptlog.Error(e)
	}
}

func (redisPool *ptRedisPool) SetExpire(key string, expire int) error {
	rp, e := redisPool.doCmd(EXPIRE, key, expire)
	if rp == 1 {
		return nil
	}
	return e
}

func (redisPool *ptRedisPool) SetString(key, str string, expire int) error {
	var e error
	if expire > 0 {
		_, e = redisPool.doCmd(SET, key, str, EX, expire)
	} else {
		_, e = redisPool.doCmd(SET, key, str)
	}
	return e
}

func (redisPool *ptRedisPool) GetString(key string) string {
	rp, e := redisPool.doCmd(GET, key)
	if e != nil {
		ptlog.Error(e)
		return ""
	}
	if rp == nil {
		return ""
	}
	return string(rp.([]uint8))
}

func (redisPool *ptRedisPool) DelString(key string) error {
	return redisPool.delKeys(key)
}

func (redisPool *ptRedisPool) SetMHash(key string, hashMap map[string]string) error {
	var args = make([]interface{}, 0, len(hashMap))
	args = append(args, key)
	for k, v := range hashMap {
		args = append(args, k)
		args = append(args, v)
	}
	_, e := redisPool.doCmd(HashMSET, args...)
	return e
}

func (redisPool *ptRedisPool) GetMHash(key string, field ...interface{}) map[string]string {
	res := make(map[string]string)
	args := make([]interface{}, 0, len(field)+1)
	args = append(args, key)
	args = append(args, field...)
	rp, e := redisPool.doCmd(HashMGET, args...)
	if e != nil {
		ptlog.Error(e)
	} else {
		for k, v := range field {
			value := rp.([]interface{})[k]
			switch value.(type) {
			case nil: //该field不存在
			case []uint8:
				res[v.(string)] = string(value.([]uint8))
			default:
				res[v.(string)] = "UNDEFINED TYPE"
			}
		}
	}
	return res
}

func (redisPool *ptRedisPool) SetHash(key string, filed string, value string) error {
	_, e := redisPool.doCmd(HashSET, key, filed, value)
	return e
}

func (redisPool *ptRedisPool) GetHash(key, field string) string {
	rp, e := redisPool.doCmd(HashGET, key, field)
	if e != nil || rp == nil {
		ptlog.Error(e)
		return ""
	}
	return string(rp.([]uint8))
}

func (redisPool *ptRedisPool) GetHashAll(key string) map[string]string {
	res := make(map[string]string)
	rp, e := redisPool.doCmd(HashGETALL, key)
	if e != nil {
		ptlog.Error(e)
	} else {
		for i := 0; i < len(rp.([]interface{})); i = i + 2 {
			res[string(rp.([]interface{})[i].([]uint8))] = string(rp.([]interface{})[i+1].([]uint8))
		}
	}
	return res
}

func (redisPool *ptRedisPool) DelHashField(key string, field ...string) error {
	var args = make([]interface{}, 0, len(field)+1)
	args = append(args, key)
	for _, v := range field {
		args = append(args, v)
	}
	_, e := redisPool.doCmd(HashDEL, args...)
	return e
}

func (redisPool *ptRedisPool) ListPush(key string, value ...string) error {
	var args = make([]interface{}, 0, len(value)+1)
	args = append(args, key)
	for _, v := range value {
		args = append(args, v)
	}
	_, e := redisPool.doCmd(ListPUSH, args...)
	return e
}

func (redisPool *ptRedisPool) ListPop(key string) string {
	rp, e := redisPool.doCmd(ListPOP, key)
	if e != nil {
		ptlog.Error(e)
		return ""
	}
	switch rp.(type) {
	case nil:
		return ""
	default:
		return string(rp.([]uint8))
	}
}

//SetAdd 必须要添加一个成员
func (redisPool *ptRedisPool) SetAdd(key string, member1 string, members ...string) error {
	var args = make([]interface{}, 0, len(members)+2)
	args = append(args, key, member1)
	for _, v := range members {
		args = append(args, v)
	}
	_, e := redisPool.doCmd(SetADD, args...)
	return e
}

func (redisPool *ptRedisPool) SetMembers(key string) []string {
	var members []string
	rp, e := redisPool.doCmd(SetMEMBERS, key)
	if e != nil {
		ptlog.Error(e)
		return members
	}
	for _, v := range rp.([]interface{}) {
		members = append(members, string(v.([]uint8)))
	}
	return members
}
