package app

import (
	"errors"
	"net/http"
	"regexp"

	"alber/pkg/api"
	"alber/pkg/orm"

	"golang.org/x/crypto/bcrypt"
)

// checkPhoneAndNick check if phone & nickname is empty or not
//	exist = true - user exist in db
func checkPhoneAndNick(isExist bool, phone, nickname string) error {
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

// checkPassword check is password is valid(up) or correct password(in)
//	exist = true - user exist in db
func checkPassword(isExist bool, pass, login string) error {
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

// SignUp check validate, start session
func (app *Application) SignUp(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	nickname := r.PostFormValue("nickname")
	code := r.PostFormValue("code")
	pass := ""

	// XSS
	if api.CheckAllXSS(nickname) != nil {
		return nil, errors.New("ошибка имени")
	}

	// checking code from sms
	validPhone, ok := app.UsersCode[code]
	if !ok {
		return nil, errors.New("не корретный код")
	}

	// check phone and nick
	if e := checkPhoneAndNick(false, validPhone.Value.(string), nickname); e != nil {
		return nil, e
	}

	// generating password
	for {
		tempPass := RandomStringFromCharsetAndLength("0123456789ABCDEFJGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 12)
		if e := checkPassword(false, tempPass, ""); e == nil {
			pass = tempPass
			break
		}
	}

	hashPass, e := bcrypt.GenerateFromPassword([]byte(pass), 4)
	if e != nil {
		return nil, errors.New("ошибка сервера: пароль")
	}

	user := &orm.User{
		Nickname: nickname, PhoneNumber: validPhone.Value.(string), Password: string(hashPass),
	}
	userID, e := user.Create()
	if e != nil {
		return nil, errors.New("ошибка сервера: не удалось создать пользователя")
	}

	// start session
	if e := api.SessionStart(w, r, userID); e != nil {
		return nil, e
	}

	// send SMS with temp_password & login
	// or mb make notify on front
	return map[string]interface{}{"login": validPhone.Value.(string), "password": pass}, e
}

// SignIn check password and login from db and request + oauth2
func (app *Application) SignIn(w http.ResponseWriter, r *http.Request) (int, error) {
	phone := getPhoneNumber(r.PostFormValue("phone"))
	pass := r.PostFormValue("password")

	// checkings
	if e := api.TestPhone(phone); e != nil {
		return -1, e
	}
	if e := checkPhoneAndNick(true, phone, phone); e != nil {
		return -1, e
	}
	if e := checkPassword(true, pass, phone); e != nil {
		return -1, e
	}

	res, e := orm.GetOneFrom(orm.SQLSelectParams{
		What:    "id",
		Table:   "Users",
		Options: orm.DoSQLOption("phoneNumber = ?", "", "", phone),
		Joins:   nil,
	})
	if e != nil {
		return -1, errors.New("неправильный логин")
	}

	ID := orm.FromINT64ToINT(res[0])
	return ID, api.SessionStart(w, r, ID)
}

// Logout user
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) error {
	id := api.GetUserIDfromReq(w, r)
	if id == -1 {
		return errors.New("не зарегистрированы в сети")
	}

	if e := orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Sessions",
		Options: orm.DoSQLOption("userID = ?", "", "", id),
	}); e != nil {
		return errors.New("не произведен выход")
	}

	api.SetCookie(w, "", -1)
	return nil
}

// ResetPassword send on email message code to reset password
func (app *Application) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	phone, ok := app.UsersCode[r.PostFormValue("code")]
	if !ok {
		return errors.New("не корректный код")
	}

	newPass := r.PostFormValue("password")
	if e := checkPassword(false, newPass, ""); e != nil {
		return e
	}

	res, e := orm.GetOneFrom(orm.SQLSelectParams{
		What:    "id",
		Table:   "Users",
		Options: orm.DoSQLOption("phoneNumber = ?", "", "", phone.Value),
	})
	if e != nil {
		return errors.New("не корректный телефон")
	}

	password, e := bcrypt.GenerateFromPassword([]byte(newPass), 4)
	if e != nil {
		return errors.New("ошибка сервера: новый пароль не создан")
	}

	user := &orm.User{ID: orm.FromINT64ToINT(res[0]), Password: string(password)}
	return user.Change()
}
