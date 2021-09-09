package api

import (
	"errors"
	"net/http"
	"strconv"
	"zhibek/pkg/orm"

	"golang.org/x/crypto/bcrypt"
)

func User(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ID, e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		return nil, errors.New("wrong id")
	}

	mainQ := orm.SQLSelectParams{
		Table:   "Users AS u",
		What:    "u.*",
		Options: orm.DoSQLOption("u.id=?", "", "", ID),
	}
	parselsQ := orm.SQLSelectParams{
		Table:   "Parsels",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("userID=?", "", "", ID),
	}
	travelsQ := orm.SQLSelectParams{
		Table:   "Travelers",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("userID=?", "", "", ID),
	}

	querys := []orm.SQLSelectParams{parselsQ, travelsQ}
	as := []string{"parselsCount", "travelsCount"}

	return orm.GetWithSubqueries(
		mainQ,
		querys,
		[]string{},
		as,
		orm.User{},
	)
}

func ChangeProfile(w http.ResponseWriter, r *http.Request) error {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("not logged")
	}

	nickname, phone, pass := r.PostFormValue("nickname"), r.PostFormValue("phone"), r.PostFormValue("password")
	if CheckAllXSS(nickname, phone) != nil {
		return errors.New("wrong content")
	}
	u := &orm.User{
		ID: userID, Nickname: nickname, PhoneNumber: phone,
	}

	if pass != "" {
		if hashPass, e := bcrypt.GenerateFromPassword([]byte(pass), 4); e != nil {
			u.Password = string(hashPass)
		}
	}

	return u.Change()
}
