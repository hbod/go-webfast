package fast

import (
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-xorm/xorm"
)

var DB = new(db)

func init() {
	DB.db_model = db_model{DB}
}

type db struct {
	db_model
	*xorm.Engine
}

func (d *db) Open() {
	dsn := Conf.DefString("db.dsn", "")

	gdb, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		L.Fatalf("db.dsn(%s)", err)
	}

	if Conf.DefBool("db.cacher.open", false) {
		maxSize := Conf.DefInt("db.cacher.maxSize", 1000)
		expired, _ := time.ParseDuration(Conf.DefString("db.cacher.expired", "6m"))
		cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), expired, maxSize)
		gdb.SetDefaultCacher(cacher)
	}

	d.Engine = gdb
	log := Conf.DefBool("db.log", true)
	d.ShowExecTime(log)
	d.ShowSQL(log)
}
