package api

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"alber/pkg/orm"

	"golang.org/x/crypto/bcrypt"
)

// CheckPhoneAndNick check if phone & nickname is empty or not
//	exist = true - user exist in db
func CheckPhoneAndNick(isExist bool, phone, nickname string) error {
	results, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   "Users",
		What:    "phoneNumber, nickname",
		Options: orm.DoSQLOption("phoneNumber=? OR nickname=?", "", "", phone, nickname),
	})

	if e != nil && isExist {
		return errors.New("не корретный логин")
	}
	if e != nil && !isExist {
		return nil
	}
	if !isExist {
		if results[0].(string) == phone {
			return errors.New("такой телефон существует")
		}
		return errors.New("такой никнейм существует")
	}
	return nil
}

// CheckPassword check is password is valid(up) or correct password(in)
//	exist = true - user exist in db
func CheckPassword(isExist bool, pass, login string) error {
	if !isExist {
		if !regexp.MustCompile(`[A-Z]`).MatchString(pass) {
			return errors.New("пароль должени иметь латинские буквы A-Z")
		}
		if !regexp.MustCompile(`[a-z]`).MatchString(pass) {
			return errors.New("пароль должени иметь латинские буквы a-z(маленькие)")
		}
		if !regexp.MustCompile(`[0-9]`).MatchString(pass) {
			return errors.New("пароль должени иметь цифры 0-9")
		}
		if len(pass) < 8 {
			return errors.New("пароль должени иметь как минимум 8 символов")
		}
	} else {
		dbPass, e := orm.GetOneFrom(orm.SQLSelectParams{
			Table:   "Users",
			What:    "password",
			Options: orm.DoSQLOption("phoneNumber = ?", "", "", login),
		})
		if e != nil {
			return errors.New("не корретный логин")
		}

		if e := bcrypt.CompareHashAndPassword([]byte(dbPass[0].(string)), []byte(pass)); e != nil {
			return errors.New("не корретный пароль")
		}
		return nil
	}
	return nil
}

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

	// admin check
	if data[0]["phoneNumber"] == "+77787833831" {
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

	if nickname != "" || phone != "" {
		if e := CheckPhoneAndNick(false, phone, nickname); e != nil {
			return e
		}
	}

	u := &orm.User{
		ID: userID, Nickname: nickname, PhoneNumber: phone,
	}
	if pass != "" {
		if e := CheckPassword(false, pass, ""); e != nil {
			return e
		}
		if hashPass, e := bcrypt.GenerateFromPassword([]byte(pass), 4); e == nil {
			u.Password = string(hashPass)
		}
	}

	return u.Change()
}
