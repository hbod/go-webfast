package fast

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

type db_model struct {
	db *db
}

func (d db_model) MFilter(s *xorm.Session, search map[string]interface{}, filters map[string]string) {
	nsearch := make(map[string]interface{}, 0)
	for k, se := range search {
		switch k {
		case "<order>", "<group>":
			rt := reflect.TypeOf(se)
			if rt.Kind() == reflect.Slice {
				rv := reflect.ValueOf(se)
				l_se := rv.Len()

				var s_arr []string

				for i := 0; i < l_se; i++ {
					rvIndex := rv.Index(i)

					vs := fmt.Sprintf("%s", rvIndex.Interface())

					if k == "<order>" {
						ord := "ASC"
						if strings.HasPrefix(vs, "-") {
							vs = strings.TrimPrefix(vs, "-")
							ord = "DESC"
						}
						if filter_s, ok := filters[vs]; ok {
							if filter_s != "" {
								vs = filter_s
							}
							s_arr = append(s_arr, vs+" "+ord)
						}
					} else {
						if filter_s, ok := filters[vs]; ok {
							if filter_s != "" {
								vs = filter_s
							}
							s_arr = append(s_arr, vs)
						}
					}
				}

				if len(s_arr) != 0 {
					s_str := strings.Join(s_arr, ",")
					if k == "<order>" {
						s.OrderBy(s_str)
					} else {
						s.GroupBy(s_str)
					}
				}
			}

		default:
			nsearch[k] = se
		}
	}

	s.And(d.MSfilter(nsearch, filters, true))
}

func (d db_model) MSfilter(search interface{}, filters map[string]string, isand bool) builder.Cond {
	var b = builder.NewCond()

	rts := reflect.TypeOf(search)

	if rts.Kind() == reflect.Map {
		rvs := reflect.ValueOf(search)
		map_keys := rvs.MapKeys()
		l_map_key := len(map_keys)

		for rvi := 0; rvi < l_map_key; rvi++ {
			var k = map_keys[rvi].String()
			var v = rvs.MapIndex(map_keys[rvi]).Interface()

			var bcond builder.Cond

			if strings.HasPrefix(k, "<or>") {
				bcond = d.MSfilter(v, filters, false)
			} else {
				if strings.HasPrefix(k, "<and>") {
					bcond = d.MSfilter(v, filters, true)
				} else {
					ks := strings.Split(k, " ")

					prefix := "="
					if len(ks) > 1 {
						prefix = ks[0]
						k = ks[1]
					}

					if filter_s, ok := filters[k]; ok {
						if filter_s != "" {
							k = filter_s
						}
					} else {
						continue
					}

					switch prefix {
					case "<>": // Neq
						bcond = builder.Neq{
							k: v,
						}
					case "<": // Lt
						bcond = builder.Lt{
							k: v,
						}
					case "<=": // Lte
						bcond = builder.Lte{
							k: v,
						}
					case ">": // Gt
						bcond = builder.Gt{
							k: v,
						}
					case ">=": // Gte
						bcond = builder.Gte{
							k: v,
						}
					case "<-": // In
						bcond = builder.In(k, v)
					case "!<-": // NotIn
						bcond = builder.NotIn(k, v)
					case "-": // IsNull
						bcond = builder.IsNull{k}
					case "!-": // NotNull
						bcond = builder.NotNull{k}
					case "%": // Like
						bcond = builder.Like{k, fmt.Sprintf("%s", v)}
					case "<->": // Between
						switch reflect.TypeOf(v).Kind() {
						case reflect.Slice:
							rv := reflect.ValueOf(v)
							if rv.Len() > 1 {
								bcond = builder.Between{
									Col:     k,
									LessVal: rv.Index(0).Interface(),
									MoreVal: rv.Index(1).Interface(),
								}
							}
						}
					default:
						bcond = builder.Eq{
							k: v,
						}
					}
				}
			}

			if isand {
				b = b.And(bcond)
			} else {
				b = b.Or(bcond)
			}
		}
	}

	return b
}

func (d db_model) MFind(beans interface{}, bean interface{}, s map[string]interface{}, f map[string]string, after func(*xorm.Session, map[string]interface{})) (count int64, page int, limit int, err error) {
	if s == nil {
		s = make(map[string]interface{}, 0)
	}
	session := d.db.NewSession()
	defer session.Close()
	filter_data := make(map[string]interface{}, 0)
	if filterData, ok := s["filter"]; ok {
		rv_f_d := reflect.ValueOf(filterData)
		if rv_f_d.Kind() == reflect.Map {
			keys := rv_f_d.MapKeys()
			l := len(keys)
			for i := 0; i < l; i++ {
				filter_data[keys[i].String()] = rv_f_d.MapIndex(keys[i]).Interface()
			}
		}
		d.MFilter(session, filter_data, f)
	}

	if after != nil {
		after(session, filter_data)
	}

	count_session := session.Clone()
	defer count_session.Close()
	count, _ = count_session.Count(bean)

	var l_all int
	if all, ok := s["all"]; ok {
		l_all, _ = strconv.Atoi(fmt.Sprintf("%d", all))
	}

	if l_all < 1 {
		if s_page, ok := s["page"]; ok {
			page, _ = strconv.Atoi(fmt.Sprintf("%d", s_page))
		}
		if s_limit, ok := s["limit"]; ok {
			limit, _ = strconv.Atoi(fmt.Sprintf("%d", s_limit))
		}

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		session.Limit(limit, (page-1)*limit)
	}

	err = session.Find(beans)
	if err != nil {
		return
	}
	return
}

func (d db_model) MFirst(bean interface{}, w map[string]interface{}, f map[string]string) (bool, error) {
	session := d.db.NewSession()
	defer session.Close()
	d.MFilter(session, w, f)
	return session.Get(bean)
}
