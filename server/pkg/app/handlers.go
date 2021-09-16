package app

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"zhibek/pkg/api"
	"zhibek/pkg/orm"
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
		max, _ := strconv.Atoi(app.Config.MAX_REQUEST_COUNT)
		if app.CurrentRequestCount < max {
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
		wd, _ := os.Getwd()
		t, e := template.ParseFiles(wd + "/dist/index.html")
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
	type Route struct {
		Path     string  `json:"route"`
		Children []Route `json:"children"`
	}

	data := api.API_RESPONSE{
		Err:  "",
		Code: 200,
		Data: Route{
			Path: "/",
			Children: []Route{
				{Path: "/user"},
				{Path: "/parsels"},
				{Path: "/travelers"},
				{Path: "/images"},
				{Path: "/search"},
				{Path: "/toptypes"},
			},
		},
	}

	api.DoJS(w, data)
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

// HTravelTypes for handle '/api/travelTypes'
func (app *Application) HTravelTypes(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.TravelTypes)
}

// HCountryCodes for handle '/api/countryCodes'
func (app *Application) HCountryCodes(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.CountryCodes)
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

// HPreSignUpSMS for handle '/sign/sms/s'
func (app *Application) HPreSignUpSMS(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
		}

		phone := getPhoneNumber(r.PostFormValue("phone"))
		code := StringWithCharset(8)
		countryCode := r.PostFormValue("countryCode")
		msg := "Регистрирация на платформе Жибек. Код подтверждения: " + code

		// send SMS
		if e := app.SendSMS(phone, countryCode, msg); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		app.m.Lock()
		app.UsersCode[code] = &Code{Value: countryCode + phone, ExpireMin: app.CurrentMin + 60*1}
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
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		data.Data = successData

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

		phone := getPhoneNumber(r.PostFormValue("phone"))
		code := StringWithCharset(8)
		countryCode := r.PostFormValue("countryCode")
		msg := "Изменение пароля на платформе Жибек. Код подтверждения: " + code

		// send SMS
		if e := app.SendSMS(phone, countryCode, msg); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		app.m.Lock()
		app.UsersCode[code] = &Code{Value: countryCode + phone, ExpireMin: app.CurrentMin + 60*1}
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
func (app *Application) HPreChangeProfileSMS(w http.ResponseWriter, r *http.Request) {
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

		phone := getPhoneNumber(r.PostFormValue("phone"))
		countryCode := r.PostFormValue("countryCode")
		code := StringWithCharset(8)
		cd := &Code{Value: countryCode + phone, ExpireMin: app.CurrentMin + 60*1}
		msg := "Изменение аккаунта на платформе Жибек. Код подтверждения: " + code

		if phone == "" {
			phoneDB, e := orm.GetOneFrom(orm.SQLSelectParams{
				Table:   "Users",
				What:    "phoneNumber",
				Options: orm.DoSQLOption("id=?", "", "", userID),
			})
			if e != nil {
				api.SendErrorJSON(w, data, "not logged")
				return
			}
			phone = phoneDB[0].(string)
			countryCode = ""
			cd.Value = ""
		}

		data.Data = map[string]string{"login": countryCode + phone}

		app.m.Lock()
		app.UsersCode[code] = cd
		app.m.Unlock()

		// here sending sms to abonent
		if e := app.SendSMS(phone, countryCode, msg); e != nil {
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

		phone, ok := app.UsersCode[r.PostFormValue("code")]
		if !ok {
			api.SendErrorJSON(w, data, "wrong code")
			return
		}

		// set correct phone
		r.PostForm.Set("phone", phone.Value.(string))

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

// HItemUp travel or parsel up
func (app *Application) HItemUp(w http.ResponseWriter, r *http.Request) {
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

		if e := api.ItemUp(w, r); e != nil {
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
	api.HApi(w, r, api.CreateParsel)
}

// HSaveTravel create parsel
func (app *Application) HSaveTravel(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.CreateTravel)
}

// HSaveImage save one image
func (app *Application) HSaveImage(w http.ResponseWriter, r *http.Request) {
	link, name, e := uploadFile("file", r)
	if e != nil {
		return
	}
	r.PostForm.Set("link", link)
	r.PostForm.Set("filename", name)
	api.HApi(w, r, api.CreateImage)
}

// ------------------------------------------- Remove ----------------------------------------
// HRemoveParsel create parsel
func (app *Application) HRemoveParsel(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.RemoveParsel)
}

// HRemoveTravel create parsel
func (app *Application) HRemoveTravel(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.RemoveTraveler)
}

// HRemoveImage save one image
func (app *Application) HRemoveImage(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.RemoveImage)
}
