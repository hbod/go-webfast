package fast

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type Ghttp struct {
	*gin.Engine
}

func NewHttp(mode string) (g *Ghttp, err error) {
	gin.SetMode(mode)
	g = new(Ghttp)
	g.Engine = gin.New()

	return
}

func (g *Ghttp) Router(routers ...Router) {
	parseRouter(g, routers...)
}

type Router struct {
	RelPath  string
	Method   string
	Handlers []interface{}
	Groups   []Router
}

func parseRouter(o interface{}, rs ...Router) {
	l := len(rs)
	rv := reflect.ValueOf(o)

	for i := 0; i < l; i++ {
		r := rs[i]

		var rv_o_c = []reflect.Value{
			reflect.ValueOf(r.RelPath),
		}

		if r.Handlers != nil {
			for _, h := range r.Handlers {
				rv_o_c = append(rv_o_c, reflect.ValueOf(h))
			}
		}

		mName := "Group"
		if r.Method != "" {
			mName = r.Method
		}

		rv_o := rv.MethodByName(mName)
		rv_o_c_v := rv_o.Call(rv_o_c)

		if r.Groups != nil {
			if len(rv_o_c_v) > 0 {
				parseRouter(rv_o_c_v[0].Interface(), r.Groups...)
			}
		}
	}

}
