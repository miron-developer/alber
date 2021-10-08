package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"alber/pkg/orm"
)

type API_RESPONSE struct {
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

// XSS check
func xss(data string) error {
	if data == "" {
		return nil
	}

	rg := regexp.MustCompile(`<+[\w\s/]+>+`)
	if rg.MatchString(data) {
		return errors.New("xss data")
	}
	return nil
}

// TestPhone is phone number
func TestPhone(phone string, omitEmpty bool) error {
	if omitEmpty && phone == "" {
		return nil
	}

	if phone == "" && !omitEmpty {
		return errors.New("пустой номер")
	}

	rg := regexp.MustCompile(`^[\d+]+$`)
	if !rg.MatchString(phone) {
		return errors.New("не корректный номер")
	}
	return nil
}

// CheckAllXSS check xss
func CheckAllXSS(testers ...string) error {
	for _, v := range testers {
		if e := xss(v); e != nil {
			return e
		}
	}
	return nil
}

// Gets int value from string with setted default value
func getIntFromString(src string, def int) int {
	val := def
	if v, e := strconv.Atoi(src); e == nil {
		val = v
	}
	return val
}

// returm from and step
func getLimits(r *http.Request) (int, int) {
	return getIntFromString(r.FormValue("from"), 0), getIntFromString(r.FormValue("step"), 10)
}

// DoJS do json and write it
func DoJS(w http.ResponseWriter, data interface{}) {
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "Application/json")
	w.Write(js)
}

// getAuthCookie check if user is logged
func getAuthCookie(r *http.Request) (string, error) {
	cookie, e := r.Cookie(cookieName)
	if e != nil {
		return "", errors.New("не зарегистрированы в сети")
	}
	return url.QueryUnescape(cookie.Value)
}

// TimeExpire time.Now().Add(some duration) and return it by string
func TimeExpire(add time.Duration) string {
	return time.Now().Add(add).Format("2006-01-02 15:04:05")
}

// GetUserIDfromReq gets users id from requst
func GetUserIDfromReq(w http.ResponseWriter, r *http.Request) int {
	sesID, e := getAuthCookie(r)
	if sesID == "" || e != nil {
		return -1
	}

	userID, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   "Sessions",
		What:    "userID",
		Options: orm.DoSQLOption("id = ?", "", "", sesID),
	})
	if e != nil {
		return -1
	}

	// update cooks & sess
	ses := &orm.Session{ID: sesID, Expire: TimeExpire(sessionExpire)}
	ses.Change()
	SetCookie(w, sesID, int(sessionExpire/timeSecond))
	return orm.FromINT64ToINT(userID[0])
}

// GetUserID get userID from rq or from get rq
func GetUserID(w http.ResponseWriter, r *http.Request, reqID string) (int, error) {
	if reqID != "" {
		return strconv.Atoi(reqID)
	}
	if userID := GetUserIDfromReq(w, r); userID != -1 {
		return userID, nil
	}
	return -1, errors.New("не зарегистрированы в сети")
}

// SendErrorJSON send to front error
func SendErrorJSON(w http.ResponseWriter, data API_RESPONSE, err string) {
	data.Err = err
	data.Code = 401
	DoJS(w, data)
}

// HApi general handler from api
func HApi(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request) (interface{}, error)) {
	data := API_RESPONSE{
		Err:  "ok",
		Data: "",
		Code: 200,
	}

	datas, e := f(w, r)
	if e != nil {
		SendErrorJSON(w, data, e.Error())
		return
	}
	data.Data = datas
	DoJS(w, data)
}
