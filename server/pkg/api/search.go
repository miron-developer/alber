package api

import (
	"errors"
	"net/http"
	"strings"
	"zhibek/pkg/orm"
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

func searchGetTextFilter(q string, searchFields []string, op *orm.SQLOption) error {
	if q == "" {
		return nil
	}
	if xss(q) != nil {
		return errors.New("danger search text")
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
	return SearchCity(r)
}
