package app

import (
	"errors"
	"net/http"
	"strings"
	"text/template"
	"zhibek/pkg/api"
)

// SecureHeaderMiddleware set secure header option
func (app *Application) SecureHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cross-origin-resource-policy", "cross-origin")
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

// AccessLogMiddleware logging request
func (app *Application) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.CurrentRequestCount < app.MaxRequestCount {
			app.CurrentRequestCount++
			app.ILog.Printf(logingReq(r))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "service is overloaded", 529)
			app.ELog.Println(errors.New("rate < curl"))
		}
	})
}

// Hindex for handle '/'
func (app *Application) HIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, e := template.ParseFiles("/assets/index.html")
		if e != nil {
			http.Error(w, "can't load this page", 500)
			app.ELog.Println(e)
			return
		}
		t.Execute(w, nil)
	}
}

/* ------------------------------------------- API ------------------------------------------------ */

// HUser for handle '/api/'
func (app *Application) HApiIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		Possible routes:
			- /user
			- /parsels
			- /travelers
			- /images
			- /search
			- /toptypes
	`))
}

// HUser for handle '/api/user/'
func (app *Application) HUser(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.User)
}

// HParsels for handle '/api/parsels'
func (app *Application) HParsels(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.Parsels)
}

// HTravelers for handle '/api/travelers'
func (app *Application) HTravelers(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.Travelers)
}

// HTopTypes for handle '/api/toptypes'
func (app *Application) HTopTypes(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.TopTypes)
}

// HSearch for handle '/api/search'
func (app *Application) HSearch(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.Search)
}

// HClippedFiles for handle '/api/images'
func (app *Application) HClippedImages(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.Images)
}

/* --------------------------------------------- Logical ---------------------------------- */
// ---------------------------------------------- Sign ---------------------------------------

// HcheckUserLogged for handle '/status'
func (app *Application) HCheckUserLogged(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		userID := api.GetUserIDfromReq(w, r)
		if userID == -1 {
			api.SendErrorJSON(w, data, "not logged")
			return
		}

		data.Data = map[string]int{"id": userID}
		api.DoJS(w, data)
	}
}

// HPreSignUpCheck for handle '/sign/up/check'
func (app *Application) HPreSignUpCheck(w http.ResponseWriter, r *http.Request) {
	data := api.API_RESPONSE{
		Err:  "ok",
		Data: "",
		Code: 200,
	}

	phone := strings.Trim(r.PostFormValue("phone"), " ")
	nickname := r.PostFormValue("nickname")
	// check is unique phone&nickname
	if e := checkPhoneAndNick(false, phone, nickname); e != nil {
		api.SendErrorJSON(w, data, e.Error())
		return
	}
	api.DoJS(w, data)
}

// HPreSignUpSMS for handle '/sign/sms/s'
func (app *Application) HPreSignUpSMS(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		phone := strings.Trim(r.PostFormValue("phone"), " ")
		code := StringWithCharset(8)
		msg := `
			Вы собираетесь зарегистрироваться на платформе Жибек.
			Введите этот код на сайте для подтверждения: ` + code + `
		`

		// send SMS
		if e := app.SendSMS(phone, msg); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		app.m.Lock()
		app.UsersCode[code] = phone
		app.m.Unlock()
		api.DoJS(w, data)
	}
}

// HSignUp for handle '/sign/up'
func (app *Application) HSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		successData, e := app.SignUp(w, r)
		if e != nil {
			data.Err = e.Error()
		}
		if successData != nil {
			data.Data = successData
		}

		// delete unnecessary code
		app.m.Lock()
		delete(app.UsersCode, r.PostFormValue("code"))
		app.m.Unlock()

		api.DoJS(w, data)
	}
}

// HSignIn for handle '/sign/in'
func (app *Application) HSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		ID, e := app.SignIn(w, r)
		if e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		data.Data = map[string]int{"id": ID}
		api.DoJS(w, data)
	}
}

// HSaveNewPassword for handle '/sign/sms/ch'
func (app *Application) HPreChangePasswordSMS(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		phone := strings.Trim(r.PostFormValue("phone"), " ")
		code := StringWithCharset(8)
		msg := `
			Вы собираетесь изменить пароль на платформе Жибек.
			Введите этот код на сайте для подтверждения: ` + code + `
		`

		// send SMS
		if e := app.SendSMS(phone, msg); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		app.m.Lock()
		app.UsersCode[code] = phone
		app.m.Unlock()
		api.DoJS(w, data)
	}
}

// HRestore for handle '/sign/re'
func (app *Application) HResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		if e := app.ResetPassword(w, r); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		// delete unnecessary code
		app.m.Lock()
		delete(app.UsersCode, r.PostFormValue("code"))
		app.m.Unlock()

		api.DoJS(w, data)
	}
}

// HLogout for handle '/sign/out'
func (app *Application) HLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		userID := api.GetUserIDfromReq(w, r)
		if userID == -1 {
			api.SendErrorJSON(w, data, "not logged")
			return
		}

		if e := app.Logout(w, r); e != nil {
			api.SendErrorJSON(w, data, e.Error())
		}
		api.DoJS(w, data)
	}
}

// ------------------------------------------- Change ------------------------------------------

// HConfirmChangeProfile save user settings
func (app *Application) HConfirmChangeProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		phone := r.PostFormValue("phone")
		code := StringWithCharset(8)
		msg := `
			Код подтверждения измения на платформе Жибек: ` + code + `
		`
		// here sending sms to abonent
		if e := app.SendSMS(phone, msg); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		api.DoJS(w, data)
	}
}

// HChangeProfile user data
func (app *Application) HChangeProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		_, ok := app.UsersCode[r.PostFormValue("code")]
		if !ok {
			api.SendErrorJSON(w, data, "wrong code")
			return
		}

		if e := api.ChangeProfile(w, r); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		// delete unnecessary code
		app.m.Lock()
		delete(app.UsersCode, r.PostFormValue("code"))
		app.m.Unlock()

		api.DoJS(w, data)
	}
}

// HChangeParsel parsel change data
func (app *Application) HChangeParsel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		if e := api.ChangeParsel(w, r); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		api.DoJS(w, data)
	}
}

// HChangeTravel travel change data
func (app *Application) HChangeTravel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		if e := api.ChangeTravel(w, r); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		api.DoJS(w, data)
	}
}

// here will be pay confirm

// HChangeTop travel or parsel change top
func (app *Application) HChangeTop(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		// check payed code
		_, ok := app.UsersCode[r.PostFormValue("code")]
		if !ok {
			api.SendErrorJSON(w, data, "not payed yet")
			return
		}

		if e := api.ChangeTop(w, r); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		// delete unnecessary code
		app.m.Lock()
		delete(app.UsersCode, r.PostFormValue("code"))
		app.m.Unlock()

		api.DoJS(w, data)
	}
}

// ------------------------------------------- Save ------------------------------------------

// HSaveParsel create parsel
func (app *Application) HSaveParsel(w http.ResponseWriter, r *http.Request) {
	api.HSaves(w, r, api.CreateParsel)
}

// HSaveTravel create parsel
func (app *Application) HSaveTravel(w http.ResponseWriter, r *http.Request) {
	api.HSaves(w, r, api.CreateTravel)
}

// HSaveImage save one image
func (app *Application) HSaveImage(w http.ResponseWriter, r *http.Request) {
	link, name, e := uploadFile("file", r)
	if e != nil {
		return
	}
	r.PostForm.Set("link", link)
	r.PostForm.Set("filename", name)
	api.HSaves(w, r, api.CreateImage)
}
