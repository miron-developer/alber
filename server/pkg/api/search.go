package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"alber/pkg/orm"
)

func doSearch(r *http.Request, q orm.SQLSelectParams, sampleStruct interface{}, addQs []orm.SQLSelectParams, joinAs, qAs []string) (interface{}, error) {
	first, count := getLimits(r)
	q.Options.Args = append(q.Options.Args, first, count)
	return orm.GetWithSubqueries(q, addQs, joinAs, qAs, sampleStruct)
}

func removeLastFromStr(src, delim string) string {
	splitted := strings.Split(src, delim)
	return strings.Join(splitted[:len(splitted)-1], delim)
}

// searchGetCountFilter add filter where to options
// 	var where - add where condition to options if var formVal is correct
// 	var defWhere - add default where condition to options if var formVal empty or not correct
// 	var formVal - value from request form
// 	var defVal - value, add when formVal empty or wrong
// 	var omitEmptyVal - if formVal is empty shoud continue add condition to options or not
func searchGetCountFilter(where, defWhere, formVal string, defVal int, omitEmptyVal bool, op *orm.SQLOption) {
	if formVal == "" && !omitEmptyVal {
		return
	}

	val, e := strconv.Atoi(formVal)
	if e != nil {
		val = defVal
		op.Where += defWhere
	} else {
		op.Where += where
	}
	op.Where += " AND "
	op.Args = append(op.Args, val)
}

func searchGetTextFilter(q string, searchFields []string, op *orm.SQLOption) error {
	if q == "" {
		return nil
	}
	if xss(q) != nil {
		return errors.New("не корректный поисковый запрос")
	}

	op.Where += "("
	for _, v := range searchFields {
		op.Where += v + " LIKE '%" + q + "%' OR "
	}
	op.Where = removeLastFromStr(op.Where, "OR ")
	op.Where += ")"
	return nil
}

// controller for future search
func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	return SearchCity(r)
}
