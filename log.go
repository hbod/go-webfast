package fast

import (
	"os"
	"time"

	"github.com/go-logger/logger"
)

var L *log.Logger = log.Std

func init() {
	L.SetPrefix("[APP] ")
	L.SetFlags(log.Lmodule | log.Llevel | log.Lshortfile | log.Ldate | log.Ltime)

	L.SetOutputLevel(Conf.DefInt("log.level", 0))
	output := Conf.DefString("log.output", "")

	if output == "file" {
		dir := Conf.DefString("log.dir", "logs")
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			L.Errorf("日志目录创建失败(%s)[%s]", dir, err)
		} else {
			fpath := dir + "/" + "app" + time.Now().Format("2006-01") + ".log"
			var file *os.File
			file, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

			if err != nil {
				L.Errorf("日志打开失败(%s)[%s]", fpath, err)
			} else {
				log.SetOutput(file)
			}
		}
	}
}
