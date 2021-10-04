package api

import (
	"errors"
	"net/http"
	"strconv"

	"alber/pkg/orm"

	"golang.org/x/crypto/bcrypt"
)

func User(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

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

	data, e := orm.GetWithSubqueries(
		mainQ,
		querys,
		[]string{},
		as,
		orm.User{},
	)
	if e != nil {
		return nil, e
	}

	if data[0]["phoneNumber"] == "+77759339540" {
		data[0]["isAdmin"] = true
	}
	return data, nil
}

func Users(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		return nil, errors.New("wrong method")
	}

	first, step := getLimits(r)
	return orm.GeneralGet(
		orm.SQLSelectParams{
			Table:   "Users AS u",
			What:    "u.*",
			Options: orm.DoSQLOption("", "", "?,?", first, step),
		},
		nil,
		orm.User{},
	), nil
}

func ChangeProfile(w http.ResponseWriter, r *http.Request) error {
	userID := GetUserIDfromReq(w, r)
	if userID == -1 {
		return errors.New("не зарегистрированы в сети")
	}

	nickname, phone, pass := r.PostFormValue("nickname"), r.PostFormValue("phone"), r.PostFormValue("password")
	if CheckAllXSS(nickname, phone) != nil {
		return errors.New("не корректное содержимое")
	}
	u := &orm.User{
		ID: userID, Nickname: nickname, PhoneNumber: phone,
	}

	if pass != "" {
		if hashPass, e := bcrypt.GenerateFromPassword([]byte(pass), 4); e == nil {
			u.Password = string(hashPass)
		}
	}

	return u.Change()
}
