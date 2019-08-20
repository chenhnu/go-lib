package ptdb

import (
	"github.com/chenhnu/go-lib/ptlog"
	"errors"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type RedisTool struct {
	conn redis.Conn
}

func NewRedisTool(host string, port int, pwd string) *RedisTool {
	rt := &RedisTool{}
	e := rt.initTool(host, port, pwd)
	if e != nil {
		ptlog.Error(e)
	} else {
		ptlog.Debug("redis connect success")
	}
	return rt
}

func (tool *RedisTool) initTool(host string, port int, pwd string) error {
	if !isHost(host) {
		return errors.New("host illegal")
	}
	conn, e := redis.Dial("tcp", host+":"+strconv.Itoa(port))
	if e != nil {
		return e
	}
	tool.conn = conn
	if pwd != "" {
		reply, e := tool.doCmd(AUTH, pwd)
		if e != nil {
			return e
		}
		if reply != "OK" {
			return errors.New("password error")
		}

	}
	reply, e := tool.doCmd(PING)
	if e != nil {
		return e
	}
	if reply != "PONG" {
		return errors.New("password required")
	}
	return e
}

//将do和send函数封装一层，后期可以以在此处选择要使用的连接（redis pool）
//后期可以在此处添加连接中断的重新恢复（recovery connect）
func (tool *RedisTool) doCmd(cmd string, args ...interface{}) (interface{}, error) {
	return tool.conn.Do(cmd, args...)
}

//Send()函数只是将命令写进buffer里，并不会执行，而Flush()就可以执行
//可以在使用事务的时候使用，一个事务一次提交
func (tool *RedisTool) sendCmd(cmd string, args ...interface{}) error {
	e := tool.conn.Send(cmd, args...)
	e = tool.conn.Flush()
	return e
}

func (tool *RedisTool) CloseTool() error {
	e := tool.conn.Close()
	if e != nil {
		return e
	}
	ptlog.Debug("redis disconnect success")
	return e
}
func (tool *RedisTool) Keys() []string {
	var keys []string
	rp, e := tool.doCmd(KEYS, "*")
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

func (tool *RedisTool) delKeys(args ...interface{}) error {
	_, e := tool.doCmd(DEL, args...)
	return e
}

func (tool *RedisTool) ExitKey(key string) bool {
	rp, e := tool.doCmd(EXISTS, key)
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
func (tool *RedisTool) FlushAll() {
	_, e := tool.doCmd(FlushALL)
	if e != nil {
		ptlog.Error(e)
	}
}

func (tool *RedisTool) SetExpire(key string, expire int) error {
	rp, e := tool.doCmd(EXPIRE, key, expire)
	if rp == 1 {
		return nil
	}
	return e
}
func (tool *RedisTool) SetString(key, str string, expire int) error {
	var e error
	if expire > 0 {
		_, e = tool.doCmd(SET, key, str, EX, expire)
	} else {
		_, e = tool.doCmd(SET, key, str)
	}
	return e
}
func (tool *RedisTool) GetString(key string) string {
	rp, e := tool.doCmd(GET, key)
	if e != nil {
		ptlog.Error(e)
		return ""
	}
	if rp == nil {
		return ""
	}
	return string(rp.([]uint8))
}
func (tool *RedisTool) DelString(key string) error {
	return tool.delKeys(key)
}
func (tool *RedisTool) SetMHash(key string, hashMap map[string]string) error {
	var args = make([]interface{}, 0, len(hashMap))
	args = append(args, key)
	for k, v := range hashMap {
		args = append(args, k)
		args = append(args, v)
	}
	_, e := tool.doCmd(HashMSET, args...)
	return e
}
func (tool *RedisTool) GetMHash(key string, field ...interface{}) map[string]string {
	res := make(map[string]string)
	args := make([]interface{}, 0, len(field)+1)
	args = append(args, key)
	args = append(args, field...)
	rp, e := tool.doCmd(HashMGET, args...)
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
func (tool *RedisTool) SetHash(key string, filed string, value string) error {
	_, e := tool.doCmd(HashSET, key, filed, value)
	return e
}

func (tool *RedisTool) GetHash(key, field string) string {
	rp, e := tool.doCmd(HashGET, key, field)
	if e != nil || rp == nil {
		ptlog.Error(e)
		return ""
	}
	return string(rp.([]uint8))
}
func (tool *RedisTool) GetHashAll(key string) map[string]string {
	res := make(map[string]string)
	rp, e := tool.doCmd(HashGETALL, key)
	if e != nil {
		ptlog.Error(e)
	} else {
		for i := 0; i < len(rp.([]interface{})); i = i + 2 {
			res[string(rp.([]interface{})[i].([]uint8))] = string(rp.([]interface{})[i+1].([]uint8))
		}
	}
	return res
}
func (tool *RedisTool) DelHashField(key string, field ...string) error {
	var args = make([]interface{}, 0, len(field)+1)
	args = append(args, key)
	for _,v:=range field{
		args=append(args, v)
	}
	_, e := tool.doCmd(HashDEL, args...)
	return e
}
func (tool *RedisTool) ListPush(key string, value ...string) error{
	var args = make([]interface{}, 0, len(value)+1)
	args = append(args, key)
	for _,v:=range value{
		args=append(args, v)
	}
	_,e:=tool.doCmd(ListPUSH,args...)
	return e
}
func (tool *RedisTool) ListPop(key string) string {
	rp,e:=tool.doCmd(ListPOP,key)
	if e!=nil{
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
