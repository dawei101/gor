package rdb

/*
对db进行简单的封装，实现一次实例，永久使用。
使用时不需要关心连接池和重连的问题。
*/

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/dawei101/gor/rlog"
)

type rDB struct {
	*sql.DB
	dbType string
}

var _dbs map[string]*rDB = make(map[string]*rDB)
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
	_dbs[name] = &rDB{db, dbType}
	return nil
}

func getDB(name string) (*rDB, bool) {
	_db_mu.RLock()
	defer _db_mu.RUnlock()
	db, ok := _dbs[name]
	return db, ok
}

// 获得 name 的 *sql.DB
//
// 获得前必须保证 InitDB 过，否则会 panic
func DB(name string) *sql.DB {
	db, ok := getDB(name)
	if !ok {
		panic(fmt.Sprintf(" no database config named:%s", name))
	}
	return db.DB
}

// 获得 name 的 *sql.DB
//
// 获得前必须保证 `InitDB("default", "xxxx")` 过，否则会 panic
func DefaultDB() *sql.DB {
	return DB("default")
}

// 获得*sqlx.DB
//
// 获得前必须保证 InitDB 过，否则会 panic
//
// 请不要使用migration特性
//
// @see github.com/jmoiron/sqlx
func DBX(name string) *sqlx.DB {
	rdb, ok := getDB(name)
	if !ok {
		panic(fmt.Sprintf(" no database config named:%s", name))
	}
	return sqlx.NewDb(rdb.DB, rdb.dbType)
}

// 获得default *sqlx.DB
//
// 获得前必须保证 InitDB 过，否则会 panic
//
// 请不要使用migration特性
//
// @see github.com/jmoiron/sqlx
func DefaultDBX() *sqlx.DB {
	return DBX("default")
}
