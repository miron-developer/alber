package app

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"wnet/pkg/app"
	"wnet/pkg/orm"
)

func searchGetCountFilter(where, formVal string, defVal int, op *orm.SQLOption) {
	if formVal != "" {
		val, e := strconv.Atoi(formVal)
		if e != nil {
			val = defVal
		}
		op.Where += where + " ? AND"
		op.Args = append(op.Args, val)
	}
}

func searchGetBtnFilter(val string, choises [][]string, isOrder bool, op *orm.SQLOption) {
	for _, choise := range choises {
		if val == choise[0] {
			if isOrder {
				op.Order = choise[1] + " DESC"
				return
			}
			op.Where += choise[1] + " ? AND"
			op.Args = append(op.Args, choise[0])
			return
		}
	}
}

func searchGetSwitchFilter(where, formVal, onVal, offVal string, op *orm.SQLOption) {
	if formVal == "1" {
		op.Where += where + " ? AND"
		if offVal == "" {
			op.Args = append(op.Args, onVal)
		} else {
			op.Args = append(op.Args, offVal)
		}
		return
	}
}

func removeLastFromStr(src, delim string) string {
	splitted := strings.Split(src, delim)
	return strings.Join(splitted[:len(splitted)-1], delim)
}

func searchGetTextFilter(q string, searchFields []string, op *orm.SQLOption) error {
	if q != "" {
		if app.XCSSOther(q) != nil {
			return errors.New("Danger search text")
		}

		op.Where += "("
		for _, v := range searchFields {
			op.Where += v + " LIKE '%" + q + "%' OR "
		}
		op.Where = removeLastFromStr(op.Where, "OR ")
		op.Where += ")"
		return nil
	}
	op.Where = removeLastFromStr(op.Where, "AND")
	return nil
}

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := app.GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}

	searchType := r.FormValue("type")
	if searchType != "all" &&
		searchType != "user" &&
		searchType != "group" &&
		searchType != "post" &&
		searchType != "video" {
		return nil, errors.New("wrong type")
	}

	if searchType == "user" {
		return searchUser(userID, r)
	}
	if searchType == "group" {
		return searchGroup(userID, r)
	}
	if searchType == "post" {
		return searchPost(userID, r)
	}
	if searchType == "video" {
		return searchVideo(userID, r)
	}
	return nil, errors.New("'All' not supported yet")
}
