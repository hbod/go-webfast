package db

import (
	_ "github.com/lib/pq"

	"github.com/go-xorm/xorm"
)

func Postgres(dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine("postgres", dsn)
}
