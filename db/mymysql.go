package db

import (
	_ "github.com/ziutek/mymysql/godrv"

	"github.com/go-xorm/xorm"
)

func Mymysql(dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine("mymysql", dsn)
}
