package fast

import (
	"flag"

	"github.com/gookit/ini"
)

var Conf *ini.Ini

func init() {
	var (
		cfile string
		err   error
	)
	flag.StringVar(&cfile, "c", "config.ini", "配置文件路径.")
	flag.Parse()

	Conf = ini.NewWithOptions(func(o *ini.Options) {
		o.Readonly = true
	})
	err = Conf.LoadFiles(cfile)
	if err != nil {
		L.Fatalf("配置文件(ERROR:%s).", err)
	}
}
