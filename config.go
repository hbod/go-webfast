package fast

import (
	"github.com/gookit/ini"
)

func ParseIni(file string) (c *ini.Ini, err error) {
	c = ini.NewWithOptions(func(o *ini.Options) {
		o.Readonly = true
	})
	err = c.LoadFiles(file)
	return
}
