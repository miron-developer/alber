package app

import (
	"errors"
	"net/http"
	"regexp"
	"zhibek/pkg/api"
	"zhibek/pkg/orm"

	"golang.org/x/crypto/bcrypt"
)

// checkPhoneAndNick check if phone & nickname is empty or not
//	exist = true - user exist in db
func checkPhoneAndNick(isExist bool, phone, nickname string) error {
	results, _ := orm.GetFrom(orm.SQLSelectParams{
		What:    "phone, nickname",
		Table:   "Users",
		Options: orm.DoSQLOption("phone=? OR nickname=?", "", "", phone, nickname),
	})

	if !isExist && len(results) > 0 {
		if results[0][0].(string) == phone {
			return errors.New("this phone is not empty")
		}
		return errors.New("this nickname is not empty")
	}
	if isExist && len(results) == 0 {
		return errors.New("wrong login")
	}
	return nil
}

// checkPassword check is password is valid(up) or correct password(in)
//	exist = true - user exist in db
func checkPassword(isExist bool, pass, login string) error {
	if !isExist {
		if !regexp.MustCompile(`[A-Z]`).MatchString(pass) {
			return errors.New("password must have A-Z")
		}
		if !regexp.MustCompile(`[a-z]`).MatchString(pass) {
			return errors.New("password must have a-z(small)")
		}
		if !regexp.MustCompile(`[0-9]`).MatchString(pass) {
			return errors.New("password must have 0-9")
		}
		if len(pass) < 8 {
			return errors.New("password must have at least 8 character")
		}
	} else {
		dbPass, e := orm.GetOneFrom(orm.SQLSelectParams{
			What:    "password",
			Table:   "Users",
			Options: orm.DoSQLOption("email = ? OR nickname = ?", "", "", login, login),
		})
		if e != nil {
			return errors.New("wrong login")
		}
		return bcrypt.CompareHashAndPassword([]byte(dbPass[0].(string)), []byte(pass))
	}
	return nil
}

// SignUp check validate, start session
func (app *Application) SignUp(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	phone := getPhoneNumber(r.PostFormValue("phone"))
	nickname := r.PostFormValue("nickname")
	code := r.PostFormValue("code")
	pass := ""

	// XSS
	if api.CheckAllXSS(nickname) != nil {
		return nil, errors.New("danger nickname")
	}

	// checking code from sms
	if validPhone, exist := app.UsersCode[code]; !exist || validPhone.Value != phone {
		return nil, errors.New("wrong code")
	}

	// check phone and nick
	if e := checkPhoneAndNick(false, phone, nickname); e != nil {
		return nil, e
	}

	// generating password
	for {
		tempPass := StringWithCharset(12)
		if e := checkPassword(false, tempPass, ""); e == nil {
			pass = tempPass
			break
		}
	}

	hashPass, e := bcrypt.GenerateFromPassword([]byte(pass), 4)
	if e != nil {
		return nil, errors.New("internal server error: password")
	}

	user := &orm.User{
		Nickname: nickname, PhoneNumber: phone, Password: string(hashPass),
	}
	userID, e := user.Create()
	if e != nil {
		return nil, e
	}

	// start session
	if e := api.SessionStart(w, r, userID); e != nil {
		return nil, e
	}

	// send SMS with temp_password & login
	// or mb make notify on front
	return map[string]interface{}{"login": phone, "password": pass}, e
}

// SignIn check password and login from db and request + oauth2
func (app *Application) SignIn(w http.ResponseWriter, r *http.Request) (int, error) {
	phone := getPhoneNumber(r.PostFormValue("phone"))
	pass := r.PostFormValue("password")

	// checkings
	if e := checkPhoneAndNick(true, phone, ""); e != nil {
		return -1, e
	}
	if e := checkPassword(true, pass, ""); e != nil {
		return -1, errors.New("password is not correct")
	}

	res, e := orm.GetOneFrom(orm.SQLSelectParams{
		What:    "id",
		Table:   "Users",
		Options: orm.DoSQLOption("phone = ?", "", "", phone),
		Joins:   nil,
	})
	if e != nil {
		return -1, errors.New("wrong login")
	}

	ID := orm.FromINT64ToINT(res[0])
	return ID, api.SessionStart(w, r, ID)
}

// Logout user
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) error {
	id := api.GetUserIDfromReq(w, r)
	if id == -1 {
		return errors.New("not logged")
	}

	if e := orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Sessions",
		Options: orm.DoSQLOption("userID = ?", "", "", id),
	}); e != nil {
		return errors.New("not logouted")
	}

	api.SetCookie(w, "", -1)
	return nil
}

// ResetPassword send on email message code to reset password
func (app *Application) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	phone, ok := app.UsersCode[r.PostFormValue("code")]
	if !ok {
		return errors.New("wrong code")
	}

	newPass := r.PostFormValue("password")
	if e := checkPassword(false, newPass, ""); e != nil {
		return e
	}

	res, e := orm.GetOneFrom(orm.SQLSelectParams{
		What:    "id",
		Table:   "Users",
		Options: orm.DoSQLOption("phone = ?", "", "", phone.Value),
	})
	if e != nil {
		return errors.New("password do not changed")
	}

	password, e := bcrypt.GenerateFromPassword([]byte(newPass), 4)
	if e != nil {
		return errors.New("the new password do not created")
	}

	user := &orm.User{ID: orm.FromINT64ToINT(res[0]), Password: string(password)}
	return user.Change()
}
