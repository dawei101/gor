// +build !testing

package rsql

import (
	"context"
	"fmt"

	"github.com/dawei101/gor/rconfig"
	"github.com/dawei101/gor/rlog"
)

func init() {
	dbs := map[string]DBConfig{}
	rconfig.DefConf().ValTo("rsql", &dbs)
	for name, dbc := range dbs {
		if err := Reg(name, dbc.DBType, dbc.DataSource, dbc.MaxOpenConns); err != nil {
			panic(fmt.Sprint("could not create db:", dbc.DataSource, err.Error()))
		}
	}
	if conn := DefConn(); conn == nil {
		rlog.Warning(context.Background(), "no default db found in config!")
	}
	SetLogging(rconfig.DefConf().IsDev())
}
