package api

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"alber/pkg/orm"
)

func Parsels(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// joins
	userJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "p.userID = u.id")
	fromJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS cf", "p.fromID = cf.id")
	toJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS ct", "p.toID = ct.id")
	topJ := orm.DoSQLJoin(orm.LOJOINQ, "TopTypes AS tt", "p.topTypeID = tt.id")

	op := orm.DoSQLOption("", "p.creationDatetime DESC, tt.id DESC", "?,?")

	if r.FormValue("type") == "user" {
		userID := GetUserIDfromReq(w, r)
		if userID == -1 {
			return nil, errors.New("не зарегистрированы в сети")
		}
		op.Where = "p.userID = ?"
		op.Args = append(op.Args, userID)
	} else {
		// add filters
		// from Almaty to Astana by default
		searchGetCountFilter(" p.fromID = ?", "p.fromID > ?", r.FormValue("fromID"), 0, true, &op)
		searchGetCountFilter(" p.toID = ?", "p.toID > ?", r.FormValue("toID"), 0, true, &op)

		// expires date between now and in 1 month
		searchGetCountFilter(" p.weight <= ?", " p.weight <= ?", r.FormValue("weight"), 100000, true, &op)
		searchGetCountFilter(" p.price >= ?", " p.price >= ?", r.FormValue("price"), 0, true, &op)
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
		return nil, errors.New("не зарегистрированы в сети")
	}

	description, contactNumber := r.PostFormValue("description"), r.PostFormValue("contactNumber")
	if CheckAllXSS(description, contactNumber) != nil {
		return nil, errors.New("содержимое не корректно")
	}
	if description == "" {
		return nil, errors.New("заполните описание")
	}

	// check phone number
	if e := TestPhone(contactNumber); e != nil {
		return nil, e
	}

	price, e := strconv.Atoi(r.PostFormValue("price"))
	weight, e2 := strconv.Atoi(r.PostFormValue("weight"))
	if e != nil || e2 != nil ||
		price*weight == 0 {
		return nil, errors.New("не корректные цена или вес")
	}

	from, e := strconv.Atoi(r.PostFormValue("fromID"))
	to, e2 := strconv.Atoi(r.PostFormValue("toID"))
	if e != nil || e2 != nil || from*to == 0 {
		return nil, errors.New("не корректные точки отправки или прибытия")
	}

	now := int(time.Now().Unix() * 1000)
	p := &orm.Parsel{
		Description: description, ContactNumber: contactNumber,
		Price: price, Weight: weight, IsHaveWhatsUp: "0",
		UserID: userID, FromID: from, ToID: to,
		CreationDatetime: now, ExpireDatetime: now + 86400*1000*30,
	}

	if r.PostFormValue("isHaveWhatsUp") == "1" {
		p.IsHaveWhatsUp = "1"
	}

	parselID, e := p.Create()
	if e != nil {
		return nil, errors.New("не создана посылка")
	}
	return parselID, nil
}

// ChangeParsel change one parsel
func ChangeParsel(w http.ResponseWriter, r *http.Request) error {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("не зарегистрированы в сети")
	}
	parselID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("не корректная посылка")
	}

	description, contactNumber := r.PostFormValue("description"), r.PostFormValue("contactNumber")
	if CheckAllXSS(description, contactNumber) != nil {
		return errors.New("содежимое не корректно")
	}
	if description == "" {
		return errors.New("заполните описание")
	}
	// check phone number
	if e := TestPhone(contactNumber); e != nil {
		return e
	}

	price, e := strconv.Atoi(r.PostFormValue("price"))
	weight, e2 := strconv.Atoi(r.PostFormValue("weight"))
	if (r.PostFormValue("price") != "" && e != nil && price == 0) ||
		(r.PostFormValue("weight") != "" && e2 != nil && weight == 0) {
		return errors.New("не корректные цена или вес")
	}

	from, e := strconv.Atoi(r.PostFormValue("fromID"))
	to, e2 := strconv.Atoi(r.PostFormValue("toID"))
	if (r.PostFormValue("fromID") != "" && e != nil && from == 0) ||
		(r.PostFormValue("toID") != "" && e2 != nil && to == 0) ||
		(r.PostFormValue("toID") != "" && r.PostFormValue("fromID") != "" && from == to) {
		return errors.New("не корректные точки отправки или прибытия place")
	}

	now := int(time.Now().Unix() * 1000)

	isHaveWhatsUp := r.PostFormValue("isHaveWhatsUp")
	if isHaveWhatsUp != "1" && isHaveWhatsUp != "0" && isHaveWhatsUp != "" {
		return errors.New("не корректный ватсап")
	}

	p := &orm.Parsel{
		Description: description, ContactNumber: contactNumber, IsHaveWhatsUp: isHaveWhatsUp,
		Price: price, Weight: weight,
		UserID: userID, FromID: from, ToID: to, ID: parselID,
		CreationDatetime: now, ExpireDatetime: now * 86400 * 1000 * 30,
	}
	return p.Change()
}

// RemoveParsel remove one parsel
func RemoveParsel(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("не зарегистрированы в сети")
	}
	parselID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return nil, errors.New("не корректная посылка")
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
