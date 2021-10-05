package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"alber/pkg/orm"

	uuid "github.com/satori/go.uuid"
)

const sessionExpire = 24 * time.Hour
const timeSecond = time.Second
const cookieName = "ZBK_ID"

// get uuid for session
func sessionID() string {
	var e error
	u1 := uuid.Must(uuid.NewV4(), e)
	return fmt.Sprint(u1)
}

// SetCookie set cookie with expire
func SetCookie(w http.ResponseWriter, sid string, expire int) {
	sidCook := http.Cookie{
		Name:   cookieName,
		Value:  url.QueryEscape(sid),
		MaxAge: expire,

		Path: "/",
		// Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &sidCook)
}

// SessionStart start user session
func SessionStart(w http.ResponseWriter, r *http.Request, userID int) error {
	cookie, e := r.Cookie(cookieName)
	sidFromCookie := ""
	sidFromDB := ""
	isToCreate := false

	// get all sids
	if e == nil && cookie.Value != "" {
		sidFromCookie, _ = url.QueryUnescape(cookie.Value)
	}

	res, e := orm.GetOneFrom(orm.SQLSelectParams{
		Table:   "Sessions",
		What:    "id",
		Options: orm.DoSQLOption("userID = ?", "", "", userID),
		Joins:   nil,
	})
	if res != nil && e == nil {
		sidFromDB = res[0].(string)
	}

	// select one sid
	sid := sidFromDB
	if sid == "" {
		isToCreate = true
		sid = sidFromCookie
	}
	if sid == "" {
		sid = sessionID()
	}

	// create or change session
	s := &orm.Session{ID: sid, Expire: TimeExpire(sessionExpire), UserID: userID}
	if isToCreate {
		e = s.Create()
	} else {
		e = s.Change()
	}
	if e != nil {
		return errors.New("session error")
	}

	SetCookie(w, sid, int(sessionExpire/timeSecond))
	return nil
}

// SessionGC delete expired session
func SessionGC() error {
	return orm.DeleteByParams(orm.SQLDeleteParams{
		Table:   "Sessions",
		Options: orm.DoSQLOption("datetime(expireDatetime) < datetime('"+TimeExpire(time.Nanosecond)+"')", "", ""),
	})
}
