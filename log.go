package fast

import (
	"os"
	"time"

	log "github.com/go-logger/logger"
)

func Logger(level int, output string, dir string) (l *log.Logger, err error) {
	var log_flag = log.Lmodule | log.Llevel | log.Lshortfile | log.Ldate | log.Ltime

	if output == "file" {
		if dir == "" {
			dir = "logs"
		}
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			return
		} else {
			fpath := dir + "/" + "app" + time.Now().Format("2006-01") + ".log"
			var file *os.File
			file, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

			if err != nil {
				return
			} else {
				l = log.New(file, "", log_flag)
			}
		}
	}

	if l == nil {
		l = log.New(os.Stderr, "", log_flag)
	}

	l.SetOutputLevel(level)

	return
}
