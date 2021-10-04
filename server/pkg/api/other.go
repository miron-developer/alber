package api

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"alber/pkg/orm"
)

func SearchCity(r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	op := orm.DoSQLOption("", "name DESC", "?,?")
	if e := searchGetTextFilter(r.FormValue("q"), []string{"c.name"}, &op); e != nil {
		return nil, e
	}

	q := orm.SQLSelectParams{
		Table:   "Cities AS c",
		What:    "c.*",
		Options: op,
	}
	return doSearch(r, q, orm.City{}, nil, nil, nil)
}

func Images(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("не корректный id")
	}

	mainQ := orm.SQLSelectParams{
		Table:   "Images AS i",
		What:    "i.*",
		Options: orm.DoSQLOption("i.parselID=?", "", "", ID),
	}
	if datas := orm.GeneralGet(mainQ, nil, orm.Image{}); datas != nil {
		return datas, nil
	}
	return nil, errors.New("н/д")
}

func TravelTypes(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	return orm.GeneralGet(orm.SQLSelectParams{
		Table:   "TravelTypes AS tRt",
		What:    "tRt.*",
		Options: orm.DoSQLOption("", "id ASC", ""),
	}, nil, orm.TravelType{}), nil
}

func TopTypes(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	return orm.GeneralGet(orm.SQLSelectParams{
		Table:   "TopTypes AS tt",
		What:    "tt.*",
		Options: orm.DoSQLOption("", "id ASC", ""),
	}, nil, orm.TopType{}), nil
}

func CountryCodes(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	return orm.GeneralGet(orm.SQLSelectParams{
		Table:   "CountryCodes AS cc",
		What:    "cc.*, c.name",
		Options: orm.DoSQLOption("", "", ""),
		Joins:   []orm.SQLJoin{orm.DoSQLJoin(orm.INJOINQ, "Countries AS c", "cc.countryID = c.id")},
	}, []string{"country"}, orm.CountryCode{}), nil
}

// CreateImage create one image
func CreateImage(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("не зарегистрированы в сети")
	}

	link, name := r.PostFormValue("link"), r.PostFormValue("filename")
	i := &orm.Image{
		Source: link, Name: name,
		UserID: userID,
	}

	parselID, e := strconv.Atoi(r.PostFormValue("whomID"))
	if e != nil {
		return nil, errors.New("не корректная посылка")
	}
	i.ParselID = parselID

	if _, e = i.Create(); e != nil {
		wd, _ := os.Getwd()
		os.Remove(wd + i.Source)
		return nil, errors.New("не удалось прикрепить фото")
	}
	return nil, nil
}

// ChangeTop change one parsel's or travel's expire on top
func ChangeTop(w http.ResponseWriter, r *http.Request) error {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("не зарегистрированы в сети")
	}
	ID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("не корректный id")
	}

	table := "Parsels"
	if r.PostFormValue("type") == "traveler" {
		table = "Travelers"
	}

	topID, e := strconv.Atoi(r.PostFormValue("topID"))
	if e != nil {
		return errors.New("wrong try to up")
	}
	duration, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   "TopTypes",
		What:    "duration",
		Options: orm.DoSQLOption("id = ?", "", "", topID),
	})
	if e != nil {
		return errors.New("ошибка сервера: toptype")
	}

	newExpire := int(time.Now().Unix()*1000) + duration[0].(int)
	expire, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   table,
		What:    "expireDatetime",
		Options: orm.DoSQLOption("userID = ? AND id = ?", "", "1", userID, ID),
	})
	if e != nil {
		return errors.New("не корректный id")
	}

	if expire[0].(int) < newExpire {
		newExpire = expire[0].(int)
	}

	if table == "Parsels" {
		p := &orm.Parsel{
			UserID: userID, ID: ID, TopTypeID: topID, ExpireOnTopDatetime: newExpire,
		}
		return p.Change()
	} else {
		t := &orm.Traveler{
			UserID: userID, ID: ID, TopTypeID: topID, ExpireOnTopDatetime: newExpire,
		}
		return t.Change()
	}
}

// ItemUp change one parsel's or travel's creation date
func ItemUp(w http.ResponseWriter, r *http.Request) error {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("не зарегистрированы в сети")
	}
	ID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("не корректный id")
	}

	table := "Parsels"
	if r.PostFormValue("type") == "traveler" {
		table = "Travelers"
	}

	now := int(time.Now().Unix() * 1000)

	// check one day between
	cdate, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   table,
		What:    "creationDatetime",
		Options: orm.DoSQLOption("id=?", "", "", ID),
	})
	if e != nil || orm.FromINT64ToINT(cdate[0]) > now-84600000 {
		return errors.New("последнее поднятие было раньше чем день")
	}

	if table == "Parsels" {
		p := &orm.Parsel{
			UserID: userID, ID: ID, CreationDatetime: now,
		}
		return p.Change()
	} else {
		t := &orm.Traveler{
			UserID: userID, ID: ID, CreationDatetime: now,
		}
		return t.Change()
	}
}

// RemoveImage remove one image
func RemoveImage(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("не зарегистрированы в сети")
	}
	imgID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return nil, errors.New("не корректный id  фото")
	}

	wd, _ := os.Getwd()
	os.Remove(wd + r.PostFormValue("src"))

	return nil, orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Images",
		Options: orm.DoSQLOption("id=? AND userID=?", "", "", imgID, userID),
	})
}
