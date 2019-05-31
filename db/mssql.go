package db

import (
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/go-xorm/xorm"
)

func Mssql(dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine("mssql", dsn)
}
