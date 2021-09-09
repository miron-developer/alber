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
	userJ := orm.DoSQLJoin(orm.LOJOINQ, "Users AS u", "p.userID = u.id")
	fromJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS cf", "p.fromID = cf.id")
	toJ := orm.DoSQLJoin(orm.LOJOINQ, "Cities AS ct", "p.toID = ct.id")
	topJ := orm.DoSQLJoin(orm.LOJOINQ, "TopTypes AS tt", "p.topTypeID = tt.id")
	typeJ := orm.DoSQLJoin(orm.LOJOINQ, "TravelTypes AS tRt", "p.topTypeID = tRt.id")

	first, count := getLimits(r)
	op := orm.DoSQLOption("", "creationDatetime DESC", "?,?", first, count)
	if r.FormValue("type") == "user" {
		userID, e := GetUserID(w, r, "id")
		if e != nil {
			return nil, errors.New("not logged")
		}
		op.Where = "t.userID = ?"
		op.Args = append([]interface{}{}, userID, op.Args)
	}

	mainQ := orm.SQLSelectParams{
		Table:   "Travelers as t",
		What:    "t.*, u.nickname, cf.name, ct.name, tt.name, tt.color, tRt.name",
		Options: op,
		Joins:   []orm.SQLJoin{userJ, fromJ, toJ, topJ, typeJ},
	}

	return orm.GetWithSubqueries(
		mainQ,
		nil,
		[]string{"nickname", "from", "to", "onTop", "color", "travelType"},
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

	weight, e := strconv.Atoi(r.PostFormValue("weight"))
	if e != nil || weight == 0 {
		return nil, errors.New("wrong weigth")
	}

	from, e := strconv.Atoi(r.PostFormValue("from"))
	to, e2 := strconv.Atoi(r.PostFormValue("to"))
	travelType, e3 := strconv.Atoi(r.PostFormValue("to"))
	if e != nil || e2 != nil || e3 != nil || from*to*travelType == 0 {
		return nil, errors.New("wrong from or to place, or travel type")
	}

	t := &orm.Traveler{
		Weight: weight, IsHaveWhatsUp: "0",
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
