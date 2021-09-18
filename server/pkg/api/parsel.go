package api

import (
	"errors"
	"net/http"
	"os"
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

	// get not my
	userID := GetUserIDfromReq(w, r)
	op := orm.DoSQLOption("", "p.creationDatetime DESC, tt.id DESC", "?,?")
	if userID != -1 {
		op.Where = "p.userID != ? AND"
		op.Args = append(op.Args, userID)
	}

	if r.FormValue("type") == "user" {
		if userID == -1 {
			return nil, errors.New("not logged")
		}
		op.Where = "p.userID = ?"
	} else {
		// add filters
		// from Almaty to Astana by default
		searchGetCountFilter(" p.fromID = ?", "p.fromID > ?", r.FormValue("fromID"), 0, true, &op)
		searchGetCountFilter(" p.toID = ?", "p.toID > ?", r.FormValue("toID"), 0, true, &op)

		// expires date between now and in 1 month
		searchGetCountFilter(" p.expireDatetime >= ?", " p.expireDatetime >= ?", r.FormValue("startDT"), int(time.Now().Unix())*1000, true, &op)
		searchGetCountFilter(" p.expireDatetime <= ?", " p.expireDatetime <= ?", r.FormValue("endDT"), int(time.Now().Unix())*1000+86400000*30, true, &op)
		op.Where = removeLastFromStr(op.Where, "AND")
	}

	first, count := getLimits(r)
	op.Args = append(op.Args, first, count)

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

	from, e := strconv.Atoi(r.PostFormValue("fromID"))
	to, e2 := strconv.Atoi(r.PostFormValue("toID"))
	if e != nil || e2 != nil || from*to == 0 {
		return nil, errors.New("wrong from or to")
	}

	p := &orm.Parsel{
		Title: title, ContactNumber: countryCode + contactNumber,
		Price: price, Weight: weight, IsHaveWhatsUp: "0",
		UserID: userID, FromID: from, ToID: to,
		CreationDatetime: int(time.Now().Unix() * 1000),
	}

	expire, e := strconv.Atoi(r.PostFormValue("expireDatetime"))
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

	from, e := strconv.Atoi(r.PostFormValue("fromID"))
	to, e2 := strconv.Atoi(r.PostFormValue("toID"))
	if (r.PostFormValue("fromID") != "" && e != nil && from == 0) ||
		(r.PostFormValue("toID") != "" && e2 != nil && to == 0) ||
		(r.PostFormValue("toID") != "" && r.PostFormValue("fromID") != "" && from == to) {
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

	// removing clipped photos
	if photos, e := orm.GetFrom(orm.SQLSelectParams{
		Table:   "Images",
		What:    "src",
		Options: orm.DoSQLOption("parselID = ?", "", "", parselID),
	}); e == nil && len(photos) > 0 {
		wd, _ := os.Getwd()
		for _, src := range photos {
			os.Remove(wd + src[0].(string))
		}
	}

	return nil, orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Parsels",
		Options: orm.DoSQLOption("id=? AND userID = ?", "", "", parselID, userID),
	})
}
