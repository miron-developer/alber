package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"alber/pkg/orm"
)

func Travelers(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	// joins
	userJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "t.userID = u.id")
	fromJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS cf", "t.fromID = cf.id")
	toJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS ct", "t.toID = ct.id")
	topJ := orm.DoSQLJoin(orm.LOJOINQ, "TopTypes AS tt", "t.topTypeID = tt.id")
	typeJ := orm.DoSQLJoin(orm.LOJOINQ, "TravelTypes AS tRt", "t.travelTypeID = tRt.id")

	op := orm.DoSQLOption("", "t.creationDatetime DESC, tt.id DESC", "?,?")

	if r.FormValue("type") == "user" {
		userID := GetUserIDfromReq(w, r)
		if userID == -1 {
			return nil, errors.New("не зарегистрированы в сети")
		}
		op.Where = "t.userID = ? AND"
		op.Args = append(op.Args, userID)
	}

	// add filters
	searchGetCountFilter(" t.fromID = ?", "t.fromID > ?", r.FormValue("fromID"), 2, false, &op)
	searchGetCountFilter(" t.toID = ?", "t.toID > ?", r.FormValue("toID"), 1, false, &op)
	op.Where = removeLastFromStr(op.Where, "AND")

	first, count := getLimits(r)
	op.Args = append(op.Args, first, count)

	mainQ := orm.SQLSelectParams{
		Table:   "Travelers AS t",
		What:    "t.*, u.nickname, cf.name, ct.name, tt.name, tt.color, tRt.name, tRt.id",
		Options: op,
		Joins:   []orm.SQLJoin{userJ, fromJ, toJ, topJ, typeJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		nil,
		[]string{"nickname", "from", "to", "onTop", "color", "travelType", "travelTypeID"},
		nil,
		orm.Traveler{},
	)
}

// CreateTravel create one travel
func CreateTravel(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "GET" {
		return nil, errors.New("wrong method")
	}

	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("не зарегистрированы в сети")
	}

	contactNumber, description := r.PostFormValue("contactNumber"), r.PostFormValue("description")
	if e := CheckAllXSS(contactNumber, description); e != nil {
		return nil, errors.New("не корректный телефон или описание")
	}
	if description == "" {
		return nil, errors.New("заполните описание")
	}

	// check phone number
	if e := TestPhone(contactNumber, false); e != nil {
		return nil, e
	}

	from, e := strconv.Atoi(r.PostFormValue("fromID"))
	to, e2 := strconv.Atoi(r.PostFormValue("toID"))
	travelType, e3 := strconv.Atoi(r.PostFormValue("travelTypeID"))
	if e != nil || e2 != nil || e3 != nil || from*to*travelType == 0 {
		return nil, errors.New("не корректные точки отправки и прибытия, или транспорт")
	}

	t := &orm.Traveler{
		Description: description, IsHaveWhatsUp: "0", ContactNumber: contactNumber,
		UserID: userID, FromID: from, ToID: to, TravelTypeID: travelType,
		CreationDatetime: int(time.Now().Unix() * 1000),
	}

	if r.PostFormValue("isHaveWhatsUp") == "1" {
		t.IsHaveWhatsUp = "1"
	}

	travelID, e := t.Create()
	if e != nil {
		return nil, errors.New("не создан попутчик")
	}
	return travelID, nil
}

// ChangeTravel change one travel
func ChangeTravel(w http.ResponseWriter, r *http.Request) error {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("не зарегистрированы в сети")
	}
	travelID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("не корректный id")
	}

	contactNumber, description := r.PostFormValue("contactNumber"), r.PostFormValue("description")
	if e := CheckAllXSS(contactNumber, description); e != nil {
		return errors.New("не корректный телефон или описание")
	}

	// check phone number
	if e := TestPhone(contactNumber, false); e != nil && contactNumber != "" {
		return e
	}

	from, e := strconv.Atoi(r.PostFormValue("fromID"))
	to, e2 := strconv.Atoi(r.PostFormValue("toID"))
	travelType, e3 := strconv.Atoi(r.PostFormValue("travelTypeID"))
	if (r.PostFormValue("fromID") != "" && e != nil && from == 0) ||
		(r.PostFormValue("toID") != "" && e2 != nil && to == 0) ||
		(r.PostFormValue("toID") != "" && r.PostFormValue("fromID") != "" && from == to) ||
		(r.PostFormValue("travelTypeID") != "" && e3 != nil && travelType == 0) {
		return errors.New("не корректные точки отправки и прибытия, или транспорт")
	}

	isHaveWhatsUp := r.PostFormValue("isHaveWhatsUp")
	if isHaveWhatsUp != "1" && isHaveWhatsUp != "0" && isHaveWhatsUp != "" {
		return errors.New("не корректный ватсап")
	}

	t := &orm.Traveler{
		Description: description, IsHaveWhatsUp: isHaveWhatsUp, ContactNumber: contactNumber,
		UserID: userID, FromID: from, ToID: to, TravelTypeID: travelType, ID: travelID,
		CreationDatetime: int(time.Now().Unix() * 1000),
	}
	if e := t.Change(); e != nil {
		return errors.New("ошибка пунктов (пункты не должны быть одинаковы) или длины описания (макс 1000) ")
	}
	return nil
}

// RemoveTraveler remove one traveler
func RemoveTraveler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "GET" {
		return nil, errors.New("wrong method")
	}

	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("не зарегистрированы в сети")
	}
	ID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return nil, errors.New("не корректный id")
	}

	return nil, orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Travelers",
		Options: orm.DoSQLOption("id=? AND userID=?", "", "", ID, userID),
	})
}
