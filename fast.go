package fast

import (
	"flag"

	logger "github.com/go-logger/logger"
	"github.com/gookit/ini"
)

type Fast struct {
	Conf   *ini.Ini
	Db     *DbModel
	Http   *Ghttp
	Logger *logger.Logger
}

func New(cfile string) (f *Fast, err error) {
	if cfile == "" {
		cfile = "config.ini"
	}
	flag.StringVar(&cfile, "c", cfile, "配置文件路径.")
	flag.Parse()

	var (
		c *ini.Ini
		d *DbModel
		g *Ghttp
		l *logger.Logger
	)
	f = new(Fast)

	c, err = ParseIni(cfile)
	if err != nil {
		return
	}
	f.Conf = c

	l, err = Logger(c.DefInt("log.level", 0),
		c.DefString("log.output", ""),
		c.DefString("log.dir", ""))
	if err != nil {
		return
	}
	f.Logger = l

	d, err = DbOpen(c.DefString("db.driver", "mysql"),
		c.DefString("db.dsn", ""),
		c.DefBool("db.log", false),
		c.DefBool("db.cacherOpen", false),
		c.DefInt("db.cacherMaxSize", 0),
		c.DefString("db.cacherExpired", "6m"))
	if err != nil {
		return
	}
	f.Db = d

	g, err = NewHttp(c.DefString("app.mode", "debug"))
	if err != nil {
		return
	}
	f.Http = g

	return
}
