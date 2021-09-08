package app

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
	"wnet/pkg/orm"
)

func StringWithCharset(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// write in log each request
func logingReq(r *http.Request) string {
	return fmt.Sprintf("%v %v: '%v'\n", r.RemoteAddr, r.Method, r.URL)
}

// XCSSOther check
func XCSSOther(data string) error {
	if data == "" {
		return nil
	}

	rg := regexp.MustCompile(`<+[\w\s/]+>+`)
	if rg.MatchString(data) {
		return errors.New("xss data")
	}
	return nil
}

// TimeExpire time.Now().Add(some duration) and return it by string
func TimeExpire(add time.Duration) string {
	return time.Now().Add(add).Format("2006-01-02 15:04:05")
}

// DoBackup make backup every 30 min
func (app *Application) DoBackup() error {
	return orm.UploadDB()
}

// getAuthCookie check if user is logged
func getAuthCookie(r *http.Request) (string, error) {
	cookie, e := r.Cookie(cookieName)
	if e != nil {
		return "", errors.New("cookie not founded")
	}
	return url.QueryUnescape(cookie.Value)
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
	setCookie(w, sesID, int(sessionExpire/timeSecond))
	return orm.FromINT64ToINT(userID[0])
}

func CarmaCountQ(column string, eq interface{}) orm.SQLSelectParams {
	return orm.SQLSelectParams{
		Table:   "Likes",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption(column+" = ?", "", "", eq),
	}
}

func EventAnswerQ(answer, eventID interface{}) orm.SQLSelectParams {
	return orm.SQLSelectParams{
		Table:   "EventAnswers",
		What:    "COUNT(id)",
		Options: orm.DoSQLOption("answer = ? AND eventID = ?", "", "", answer, eventID),
	}
}

func GetUserID(w http.ResponseWriter, r *http.Request, reqID string) (int, error) {
	if reqID != "" {
		return strconv.Atoi(reqID)
	}
	if userID := GetUserIDfromReq(w, r); userID != -1 {
		return userID, nil
	}
	return -1, errors.New("not logged")
}
