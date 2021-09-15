package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"zhibek/pkg/orm"
)

func Travelers(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// joins
	userJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "t.userID = u.id")
	fromJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS cf", "t.fromID = cf.id")
	toJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS ct", "t.toID = ct.id")
	topJ := orm.DoSQLJoin(orm.LOJOINQ, "TopTypes AS tt", "t.topTypeID = tt.id")
	typeJ := orm.DoSQLJoin(orm.LOJOINQ, "TravelTypes AS tRt", "t.travelTypeID = tRt.id")

	op := orm.DoSQLOption("", "t.creationDatetime DESC AND tt.id DESC", "?,?")
	if r.FormValue("type") == "user" {
		userID := GetUserIDfromReq(w, r)
		if userID == -1 {
			return nil, errors.New("not logged")
		}
		op.Where = "t.userID = ?"
		op.Args = append([]interface{}{userID}, op.Args...)
	}

	// add filters
	searchGetCountFilter(" t.fromID =", r.FormValue("from"), 1, &op)
	searchGetCountFilter(" t.toID =", r.FormValue("to"), 2, &op)
	searchGetCountFilter(" t.departureDatetime >=", r.FormValue("departure"), 0, &op)
	searchGetCountFilter(" t.arrivalDatetime <=", r.FormValue("arrival"), int(time.Now().Unix())*1000, &op)

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
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}

	contactNumber, countryCode := r.PostFormValue("contactNumber"), r.PostFormValue("countryCode")
	if e := CheckAllXSS(contactNumber, countryCode); e != nil {
		return nil, errors.New("wrong number")
	}

	weight, e := strconv.Atoi(r.PostFormValue("weight"))
	if e != nil || weight == 0 {
		return nil, errors.New("wrong weigth")
	}

	from, e := strconv.Atoi(r.PostFormValue("from"))
	to, e2 := strconv.Atoi(r.PostFormValue("to"))
	travelType, e3 := strconv.Atoi(r.PostFormValue("travelType"))
	if e != nil || e2 != nil || e3 != nil || from*to*travelType == 0 {
		return nil, errors.New("wrong from or to place, or travel type")
	}

	t := &orm.Traveler{
		Weight: weight, IsHaveWhatsUp: "0", ContactNumber: countryCode + contactNumber,
		UserID: userID, FromID: from, ToID: to, TravelTypeID: travelType,
		CreationDatetime: int(time.Now().Unix() * 1000),
	}

	departure, e := strconv.Atoi(r.PostFormValue("departure"))
	arrival, e2 := strconv.Atoi(r.PostFormValue("arrival"))
	if e != nil || e2 != nil ||
		departure < t.CreationDatetime ||
		arrival < t.CreationDatetime ||
		arrival < departure {
		return nil, errors.New("wrong departure or arrival time")
	}
	t.DepartureDatetime = departure
	t.ArrivalDatetime = arrival

	if r.PostFormValue("isHaveWhatsUp") == "1" {
		t.IsHaveWhatsUp = "1"
	}

	travelID, e := t.Create()
	if e != nil {
		return nil, errors.New("not create travel")
	}
	return travelID, nil
}

// ChangeTravel change one travel
func ChangeTravel(w http.ResponseWriter, r *http.Request) error {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("not logged")
	}
	travelID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return errors.New("wrong id")
	}

	contactNumber, countryCode := r.PostFormValue("contactNumber"), r.PostFormValue("countryCode")
	if e := CheckAllXSS(contactNumber, countryCode); e != nil {
		return errors.New("wrong number")
	}

	weight, e := strconv.Atoi(r.PostFormValue("weight"))
	if r.PostFormValue("weight") != "" && e != nil || weight == 0 {
		return errors.New("wrong weigth")
	}

	from, e := strconv.Atoi(r.PostFormValue("from"))
	to, e2 := strconv.Atoi(r.PostFormValue("to"))
	travelType, e3 := strconv.Atoi(r.PostFormValue("travelType"))
	if (r.PostFormValue("from") != "" && e != nil && from == 0) ||
		(r.PostFormValue("to") != "" && e2 != nil && to == 0) ||
		(r.PostFormValue("to") != "" && r.PostFormValue("from") != "" && from == to) ||
		(r.PostFormValue("travelType") != "" && e3 != nil && travelType == 0) {
		return errors.New("wrong from or to place, or travel type")
	}

	arrivalDatetime, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   "Travelers",
		What:    "arrivalDatetime",
		Options: orm.DoSQLOption("userID = ? AND id = ?", "", "1", userID, travelID),
	})
	if e != nil {
		return errors.New("not found parsel: wrong id")
	}

	now := int(time.Now().Unix() * 1000)
	departure, e := strconv.Atoi(r.PostFormValue("departure"))
	arrival, e2 := strconv.Atoi(r.PostFormValue("arrival"))
	if (r.PostFormValue("departure") != "" && e != nil && departure < now) ||
		(r.PostFormValue("arrival") != "" && e2 != nil && arrival < now) ||
		(r.PostFormValue("arrival") != "" && r.PostFormValue("departure") != "" && arrival < departure) ||
		(r.PostFormValue("arrival") == "" && r.PostFormValue("departure") != "" && arrivalDatetime[0].(int) < departure) {
		return errors.New("wrong departure or arrival time")
	}

	isHaveWhatsUp := r.PostFormValue("isHaveWhatsUp")
	if isHaveWhatsUp != "1" && isHaveWhatsUp != "0" && isHaveWhatsUp != "" {
		return errors.New("wrong whatsup")
	}

	t := &orm.Traveler{
		Weight: weight, IsHaveWhatsUp: isHaveWhatsUp, ContactNumber: countryCode + contactNumber,
		UserID: userID, FromID: from, ToID: to, TravelTypeID: travelType,
		CreationDatetime: now, DepartureDatetime: departure, ArrivalDatetime: arrival,
	}
	return t.Change()
}

// RemoveTraveler remove one traveler
func RemoveTraveler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// get general ids
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return nil, errors.New("not logged")
	}
	parselID, e := strconv.Atoi(r.PostFormValue("id"))
	if e != nil {
		return nil, errors.New("wrong traveler")
	}

	return nil, orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Travelers",
		Options: orm.DoSQLOption("id=? AND userID=?", "", "", parselID, userID),
	})
}
