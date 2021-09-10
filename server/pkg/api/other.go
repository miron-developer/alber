package api

import (
	"errors"
	"net/http"
	"strconv"
	"zhibek/pkg/orm"
)

func SearchCity(r *http.Request) (interface{}, error) {
	op := orm.DoSQLOption("", "name DESC", "?,?")
	if e := searchGetTextFilter(r.FormValue("q"), []string{"c.name"}, &op); e != nil {
		return nil, e
	}

	q := orm.SQLSelectParams{
		Table:   "Cities AS c",
		What:    "c.name",
		Options: op,
	}
	return doSearch(r, q, orm.City{}, nil, nil, nil)
}

func Images(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("wrong id")
	}

	mainQ := orm.SQLSelectParams{
		Table:   "Images AS i",
		What:    "i.*",
		Options: orm.DoSQLOption("i.parselID=?", "", "", ID),
	}
	return orm.GeneralGet(mainQ, nil, orm.Image{}), nil
}

func TopTypes(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return orm.GeneralGet(orm.SQLSelectParams{
		Table:   "TopTypes AS tt",
		What:    "tt.*",
		Options: orm.DoSQLOption("", "id DESC", ""),
	}, nil, orm.TopType{}), nil
}

// CreateImage create one image
func CreateImage(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}

	link, name := r.PostFormValue("link"), r.PostFormValue("name")
	i := &orm.Image{
		Source: link, Name: name,
		UserID: userID,
	}

	parselID, e := strconv.Atoi(r.PostFormValue("parselID"))
	if e != nil {
		return nil, errors.New("wrong parsel")
	}
	i.ParselID = parselID

	if _, e = i.Create(); e != nil {
		return nil, errors.New("not create clipped image")
	}
	return nil, nil
}

// ChangeTop change one parsel's or travel's expire on top
func ChangeTop(w http.ResponseWriter, r *http.Request) error {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("not logged")
	}
	ID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("wrong id")
	}

	table := r.PostFormValue("type")
	expire, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   table,
		What:    "expireDatetime",
		Options: orm.DoSQLOption("userID = ? AND id = ?", "", "1", userID, ID),
	})
	if e != nil {
		return errors.New("wrong type")
	}

	expireOnTop, e := strconv.Atoi(r.PostFormValue("expireOnTop"))
	topID, e2 := strconv.Atoi(r.PostFormValue("topID"))
	if e != nil || e2 != nil || expire[0].(int) < expireOnTop {
		return errors.New("wrong try to up")
	}

	if table == "Parsels" {
		p := &orm.Parsel{
			UserID: userID, ID: ID, TopTypeID: topID, ExpireOnTopDatetime: expireOnTop,
		}
		return p.Change()
	} else {
		t := &orm.Traveler{
			UserID: userID, ID: ID, TopTypeID: topID, ExpireOnTopDatetime: expireOnTop,
		}
		return t.Change()
	}
}