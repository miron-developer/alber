package app

import (
	"fmt"
	"math/rand"
	"net/http"
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

// TODO: send SMS
// SendSMS make sending sms
func (app *Application) SendSMS(phone, msg string) error {
	return nil
}
