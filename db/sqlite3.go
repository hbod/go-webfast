package db

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-xorm/xorm"
)

func Sqlite3(dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine("sqlite3", dsn)
}
