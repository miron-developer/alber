package app

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"alber/pkg/api"
	"alber/pkg/orm"
)

// SecureHeaderMiddleware set secure header option
func (app *Application) SecureHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		next.ServeHTTP(w, r)
	})
}

// AccessLogMiddleware logging request
func (app *Application) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		max, _ := strconv.Atoi(app.Config.MAX_REQUEST_COUNT)
		if app.CurrentRequestCount < max {
			app.CurrentRequestCount++
			// loging
			// app.ILog.Printf(logingReq(r))
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
		t, e := template.ParseFiles(wd + "/website/index.html")
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
		Path        string              `json:"route"`
		Description string              `json:"description"`
		Children    []Route             `json:"children"`
		Params      []map[string]string `json:"params"`
		Methods     []string            `json:"methods"`
	}

	data := api.API_RESPONSE{
		Err:  "",
		Code: 200,
		Data: []Route{
			{Path: "/user", Description: "получение одного пользователя(полное)", Methods: []string{"GET"}, Params: []map[string]string{{"k": "id", "v": "id пользователя"}}},
			{
				Path: "/users", Description: "получение id пользователей", Methods: []string{"GET"},
				Params: []map[string]string{{"k": "from", "v": "с какого индекса"}, {"k": "step", "v": "сколько взять"}},
			},
			{
				Path: "/parsels", Description: "посылки", Methods: []string{"GET"},
				Params: []map[string]string{
					{"k": "type", "v": "user/empty"}, {"k": "weight", "v": "масса"},
					{"k": "price", "v": "цена"}, {"k": "fromID", "v": "откуда"}, {"k": "toID", "v": "куда"},
					{"k": "from", "v": "с какого индекса"}, {"k": "step", "v": "сколько взять"},
				},
			},
			{
				Path: "/travelers", Description: "попутчики", Methods: []string{"GET"},
				Params: []map[string]string{
					{"k": "type", "v": "user(если твое)/empty(общий)"}, {"k": "fromID", "v": "откуда"}, {"k": "toID", "v": "куда"},
					{"k": "from", "v": "с какого индекса"}, {"k": "step", "v": "сколько взять"},
				},
			},
			{Path: "/images", Description: "прикрепленные фото", Methods: []string{"GET"}, Params: []map[string]string{{"k": "id", "v": "parsel id"}}},
			{Path: "/search", Description: "поиск(пока ищет города с таким названием)", Methods: []string{"GET"}, Params: []map[string]string{{"k": "q", "v": "поисковый запрос"}}},
			{Path: "/toptypes", Description: "поднятия", Methods: []string{"GET"}},
			{Path: "/travelTypes", Description: "типы путешествия", Methods: []string{"GET"}},
			{Path: "/countryCodes", Description: "коды стран", Methods: []string{"GET"}},
		},
	}

	api.DoJS(w, data)
}

// HUser for handle '/api/user/'
func (app *Application) HUser(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.User)
}

// HUser for handle '/api/users/'
func (app *Application) HUsers(w http.ResponseWriter, r *http.Request) {
	api.HApi(w, r, api.Users)
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
			api.SendErrorJSON(w, data, "нет зарегистрированы в сети")
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
		if e := api.TestPhone(phone, false); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}
		code := RandomStringFromCharsetAndLength("0123456789", 6)
		countryCode := r.PostFormValue("countryCode")
		scheme := "http"
		if r.TLS != nil {
			scheme += "s"
		}
		msg := scheme + "://" + r.Host + "/sign. Ал-Бер. Код: " + code

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
		if e := api.TestPhone(phone, false); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		code := RandomStringFromCharsetAndLength("0123456789", 6)
		countryCode := r.PostFormValue("countryCode")
		scheme := "http"
		if r.TLS != nil {
			scheme += "s"
		}
		msg := scheme + "://" + r.Host + "/sign. Ал-Бер. Код: " + code

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
			api.SendErrorJSON(w, data, "не зарегистрированы в сети")
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
			api.SendErrorJSON(w, data, "не зарегистрированы в сети")
			return
		}

		phoneDB, e := orm.GetOneFrom(orm.SQLSelectParams{
			Table:   "Users",
			What:    "phoneNumber",
			Options: orm.DoSQLOption("id=?", "", "", userID),
		})
		if e != nil {
			api.SendErrorJSON(w, data, "вас не существует в базе. вы кто такой?)")
			return
		}

		phone, nickname := getPhoneNumber(r.PostFormValue("phone")), r.PostFormValue("phone")
		if e := api.TestPhone(phone, true); e != nil {
			api.SendErrorJSON(w, data, e.Error())
			return
		}

		if countryCode := r.PostFormValue("countryCode"); phone != "" {
			phone = countryCode + phone
			if e := api.CheckPhoneAndNick(false, phone, nickname); e != nil {
				api.SendErrorJSON(w, data, e.Error())
				return
			}
		}
		data.Data = map[string]string{"login": phoneDB[0].(string), "newPhone": phone}

		code := RandomStringFromCharsetAndLength("0123456789", 6)
		cd := &Code{Value: phoneDB[0].(string), ExpireMin: app.CurrentMin + 60*1}
		scheme := "http"
		if r.TLS != nil {
			scheme += "s"
		}
		msg := scheme + "://" + r.Host + "/sign. Ал-Бер. Код: " + code

		app.m.Lock()
		app.UsersCode[code] = cd
		app.m.Unlock()

		// here sending sms to abonent
		if e := app.SendSMS(phoneDB[0].(string), "", msg); e != nil {
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

		if _, ok := app.UsersCode[r.PostFormValue("code")]; !ok {
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

// HItemUp travel or parsel up
func (app *Application) HItemUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := api.API_RESPONSE{
			Err:  "ok",
			Data: "",
			Code: 200,
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
