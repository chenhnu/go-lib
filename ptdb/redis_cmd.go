package ptdb

const (
	//general command
	KEYS     = "KEYS"
	DEL      = "DEL"
	TYPE     = "TYPE"
	FlushALL = "FLUSHALL"
	OK       = "OK"
	PING     = "PING"
	AUTH     = "AUTH"
	EXPIRE   = "EXPIRE"
	EXISTS   = "EXISTS"

	//string
	SET = "SET"
	GET = "GET"
	EX  = "EX"

	//hash map
	HashGETALL = "HGETALL"
	HashMSET   = "HMSET"
	HashMGET   = "HMGET"
	HashSET    = "HSET"
	HashGET    = "HGET"
	HashDEL    = "HDEL"

	//list string
	ListPUSH = "LPUSH"
	ListPOP  = "LPOP"
)
