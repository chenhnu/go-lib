package ptdb

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

type ConnInfo struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
}

func NewConnInfo(host string, port int, dbName string, user string, pwd string, timeoutMs int) (*ConnInfo, error) {
	res := strings.Split(host, ".")
	if len(res) != 4 {
		return nil, errors.New("host error")
	}
	for _, s := range res {
		num, e := strconv.Atoi(s)
		if e != nil || num < 0 || num > 255 {
			return nil, errors.New("host error")
		}
	}
	return &ConnInfo{
		Host:     host,
		Port:     port,
		DbName:   dbName,
		User:     user,
		Password: pwd,
	}, nil
}

type MysqlPool struct {
	connInfo           *ConnInfo
	MaxSize            int
	ActiveSize         int
	TimeoutMillisecond int
	db                 *sql.DB
}

func NewMysqlPool(info *ConnInfo, maxSize int, activeSize int, timeoutms int) (*MysqlPool, error) {
	if maxSize <= 0 || activeSize <= 0 || activeSize > maxSize {
		return nil, errors.New("config error")
	}
	dataSourceName := info.User + ":" + info.Password + "@tcp(" + info.Host + ")/" + info.DbName + "?charset=utf8"
	db, e := sql.Open("mysql", dataSourceName)
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(activeSize)
	db.SetMaxIdleConns(maxSize)
	db.SetConnMaxLifetime(time.Duration(timeoutms) * time.Millisecond)
	return &MysqlPool{
		connInfo:           info,
		MaxSize:            maxSize,
		ActiveSize:         activeSize,
		TimeoutMillisecond: timeoutms,
		db:                 db,
	}, nil
}

//sql查询函数
func (pool *MysqlPool) Query(query string) ([]map[string]interface{}, error) {
	rows, e := pool.db.Query(query)
	if e != nil {
		return nil, e
	}
	colFields, e := rows.Columns()
	colTypes, e := rows.ColumnTypes()
	scanArgs := make([]interface{}, len(colFields))
	val := make([]interface{}, len(colFields))
	for i := range val {
		scanArgs[i] = &val[i]
	}
	var result []map[string]interface{}
	for rows.Next() {
		e = rows.Scan(scanArgs...)
		if e != nil {
			return nil, e
		}
		record := make(map[string]interface{})
		for i, col := range val {
			switch colTypes[i].ScanType().Name() {
			case "int32":
				temp, _ := strconv.Atoi(string(col.([]uint8)))
				record[colFields[i]] = temp
			case "RawBytes":
				record[colFields[i]] = string(col.([]uint8))
			default:
				record[colFields[i]] = col
			}
		}
		result = append(result, record)
	}
	return result, e
}
