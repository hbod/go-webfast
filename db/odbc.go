package db

import (
	_ "github.com/lunny/godbc"

	"github.com/go-xorm/xorm"
)

func Odbc(dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine("odbc", dsn)
}
