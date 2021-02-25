package rsession

import (
	"github.com/go-session/cookie"
	"github.com/go-session/mysql"
	"github.com/go-session/redis"
	"github.com/go-session/session"
)

var manager session.Manager

func init() {
	// TODO read from config
	manager = session.SetStore(mysql.NewDefaultStore(mysql.NewConfig(dsn)))
	session.InitManager(manager)
}
