package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"zhibek/pkg/orm"
)

func Parsels(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// joins
	userJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "p.userID = u.id")
	fromJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS cf", "p.fromID = cf.id")
	toJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS ct", "p.toID = ct.id")
	topJ := orm.DoSQLJoin(orm.LOJOINQ, "TopTypes AS tt", "p.topTypeID = tt.id")

	first, count := getLimits(r)
	op := orm.DoSQLOption("", "creationDatetime DESC", "?,?", first, count)
	if r.FormValue("type") == "user" {
		userID := GetUserIDfromReq(w, r)
		if userID == -1 {
			return nil, errors.New("not logged")
		}
		op.Where = "p.userID = ?"
		op.Args = append([]interface{}{userID}, op.Args...)
	}

	mainQ := orm.SQLSelectParams{
		Table:   "Parsels AS p",
		What:    "p.*, u.nickname, cf.name, ct.name, tt.name, tt.color",
		Options: op,
		Joins:   []orm.SQLJoin{userJ, fromJ, toJ, topJ},
	}
	return orm.GetWithSubqueries(
		mainQ,
		nil,
		[]string{"nickname", "from", "to", "onTop", "color"},
		nil,
		orm.Parsel{},
	)
}

// CreateParsel create one parsel
func CreateParsel(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}

	title, contactNumber, countryCode := r.PostFormValue("title"), r.PostFormValue("contactNumber"), r.PostFormValue("countryCode")
	if CheckAllXSS(title, contactNumber, countryCode) != nil {
		return nil, errors.New("wrong content")
	}

	price, e := strconv.Atoi(r.PostFormValue("price"))
	weight, e2 := strconv.Atoi(r.PostFormValue("weight"))
	if e != nil || e2 != nil ||
		price*weight == 0 {
		return nil, errors.New("wrong price or weigth")
	}

	from, e := strconv.Atoi(r.PostFormValue("from"))
	to, e2 := strconv.Atoi(r.PostFormValue("to"))
	if e != nil || e2 != nil || from*to == 0 {
		return nil, errors.New("wrong from or to")
	}

	p := &orm.Parsel{
		Title: title, ContactNumber: countryCode + contactNumber,
		Price: price, Weight: weight, IsHaveWhatsUp: "0",
		UserID: userID, FromID: from, ToID: to,
		CreationDatetime: int(time.Now().Unix() * 1000),
	}

	expire, e := strconv.Atoi(r.PostFormValue("expire"))
	if e != nil || expire < p.CreationDatetime {
		return nil, errors.New("wrong expire")
	}
	p.ExpireDatetime = expire

	if r.PostFormValue("isHaveWhatsUp") == "1" {
		p.IsHaveWhatsUp = "1"
	}

	parselID, e := p.Create()
	if e != nil {
		return nil, errors.New("not create parsel")
	}
	return parselID, nil
}

// ChangeParsel change one parsel
func ChangeParsel(w http.ResponseWriter, r *http.Request) error {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("not logged")
	}
	parselID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("wrong parsel")
	}

	title, contactNumber, countryCode := r.PostFormValue("title"), r.PostFormValue("contactNumber"), r.PostFormValue("countryCode")
	if CheckAllXSS(title, contactNumber, countryCode) != nil {
		return errors.New("danger content")
	}

	price, e := strconv.Atoi(r.PostFormValue("price"))
	weight, e2 := strconv.Atoi(r.PostFormValue("weight"))
	if (r.PostFormValue("price") != "" && e != nil && price == 0) ||
		(r.PostFormValue("weight") != "" && e2 != nil && weight == 0) {
		return errors.New("wrong price or weigth")
	}

	from, e := strconv.Atoi(r.PostFormValue("from"))
	to, e2 := strconv.Atoi(r.PostFormValue("to"))
	if (r.PostFormValue("from") != "" && e != nil && from == 0) ||
		(r.PostFormValue("to") != "" && e2 != nil && to == 0) ||
		(r.PostFormValue("to") != "" && r.PostFormValue("from") != "" && from == to) {
		return errors.New("wrong from or to place")
	}

	now := int(time.Now().Unix() * 1000)
	expire, e := strconv.Atoi(r.PostFormValue("expire"))
	if r.PostFormValue("expire") != "" && e != nil && expire < now {
		return errors.New("wrong expire")
	}

	isHaveWhatsUp := r.PostFormValue("isHaveWhatsUp")
	if isHaveWhatsUp != "1" && isHaveWhatsUp != "0" && isHaveWhatsUp != "" {
		return errors.New("wrong whatsup")
	}

	p := &orm.Parsel{
		Title: title, ContactNumber: countryCode + contactNumber, IsHaveWhatsUp: isHaveWhatsUp,
		Price: price, Weight: weight,
		UserID: userID, FromID: from, ToID: to, ID: parselID,
		CreationDatetime: now, ExpireDatetime: expire,
	}
	return p.Change()
}

// RemoveParsel remove one parsel
func RemoveParsel(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}
	parselID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return nil, errors.New("wrong parsel")
	}

	return nil, orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Parsels",
		Options: orm.DoSQLOption("id=? AND userID = ?", "", "", parselID, userID),
	})
}
