package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"time"
	"zhibek/pkg/api"
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

// DoBackup make backup every 30 min
func (app *Application) DoBackup() error {
	cmd := exec.Command("cp", `db/zhibek.db`, `db/zhibek_backup.db`)
	return cmd.Run()
}

// CheckPerMin call SessionGC per minute that delete expired sessions and do db backup
func (app *Application) CheckPerMin() {
	timer := time.NewTicker(1 * time.Minute)
	for {
		// manage timer
		<-timer.C
		timer.Reset(1 * time.Minute)

		// change conf app
		app.CurrentRequestCount = 0
		app.CurrentMin++

		// do general actions
		if app.CurrentMin == 60*24 {
			app.CurrentMin = 0
		}
		if app.CurrentMin == 30 {
			if e := app.DoBackup(); e == nil {
				app.ILog.Println("backup created!")
			} else {
				app.ELog.Println(e)
			}
		}
		if e := api.SessionGC(); e != nil {
			app.ELog.Println(e)
		}

		// remove expired codes
		go func() {
			for code, v := range app.UsersCode {
				if v.ExpireMin == app.CurrentMin {
					app.m.Lock()
					delete(app.UsersCode, code)
					app.m.Unlock()
				}
			}
		}()
	}
}

type MOBIZONE_API_RESP struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// SendSMS make sending sms
func (app *Application) SendSMS(phone, msg string) error {
	HOST := "https://api.mobizon.kz/service"
	SERVICE := "message"
	METHOD := "sendsmsmessage"

	params := url.Values{}
	params.Set("recipient", phone)
	params.Set("apiKey", app.Config.MOBIZON_API_KEY)
	params.Set("text", url.QueryEscape(msg))

	// send post rq
	resp, e := http.PostForm(HOST+SERVICE+METHOD, params)
	if e != nil {
		return errors.New("internal server error: api not response")
	}

	// get response data
	content, e := io.ReadAll(resp.Body)
	if e != nil {
		return errors.New("internal server error: content error")
	}

	// convert data to struct
	result := &MOBIZONE_API_RESP{}
	if e := json.Unmarshal(content, result); e != nil {
		return errors.New("internal server error: parse json")
	}

	// handle api errors
	if result.Message != "" || result.Code == 1 {
		return errors.New("wrong phone")
	}
	return nil
}
