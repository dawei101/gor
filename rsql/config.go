package rsql

/*
对db进行简单的封装，实现一次实例，永久使用。
使用时不需要关心连接池和重连的问题。
*/

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/dawei101/gor/rlog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	DataSource   string `yaml:"dataSource"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	DBType       string `yaml:"dbType"`
}

type DBConn struct {
	*sql.DB
	dbType string
}

var _dbs map[string]*DBConn = make(map[string]*DBConn)
var _db_mu sync.RWMutex

// 注册DB
//
// 获得前必须保证 InitDB 过，否则会 panic
//
func Reg(name, dbType string, dataSource string, maxOpenConns int) error {
	db, err := sql.Open(dbType, dataSource)
	if err != nil {
		rlog.Error(context.Background(), "could not connect to db:", dataSource)
		return err
	}
	db.SetMaxOpenConns(maxOpenConns)
	_db_mu.Lock()
	defer _db_mu.Unlock()
	_dbs[name] = &DBConn{db, dbType}
	return nil
}

func getDB(name string) (*DBConn, bool) {
	_db_mu.RLock()
	defer _db_mu.RUnlock()
	db, ok := _dbs[name]
	return db, ok
}

// 获得 name 的 *sql.DB
//
// 获得前必须保证 InitDB 过，否则会 panic
func Conn(name string) *sql.DB {
	db, ok := getDB(name)
	if !ok {
		rlog.Error(context.Background(), " no database config named:", name)
		return nil
	}
	return db.DB
}

// 获得 name 的 *sql.DB
//
// 获得前必须保证 `InitDB("default", "xxxx")` 过，否则会 panic
func DefConn() *sql.DB {
	return Conn("default")
}

// 获得*sqlx.DB
//
// 获得前必须保证 InitDB 过，否则会 panic
//
// 请不要使用migration特性
//
// @see github.com/jmoiron/sqlx
func XConn(name string) *sqlx.DB {
	rdb, ok := getDB(name)
	if !ok {
		panic(fmt.Sprintf(" no database config named:%s", name))
	}
	return sqlx.NewDb(rdb.DB, rdb.dbType)
}

// 获得default *sqlx.DB
//
// 获得前必须保证 Reg 过，否则会 panic
//
// 请不要使用migration特性
//
// @see github.com/jmoiron/sqlx
func DefXConn() *sqlx.DB {
	return XConn("default")
}
