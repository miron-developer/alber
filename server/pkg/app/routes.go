package app

import "net/http"

func (app *Application) SetRoutes() http.Handler {
	appMux := http.NewServeMux()
	appMux.HandleFunc("/", app.HIndex)

	// sign
	signMux := http.NewServeMux()
	signMux.HandleFunc("/", app.HIndex)
	signMux.HandleFunc("/sms/up", app.HPreSignUpSMS)
	signMux.HandleFunc("/up", app.HSignUp)
	signMux.HandleFunc("/in", app.HSignIn)
	signMux.HandleFunc("/sms/ch", app.HPreChangePasswordSMS)
	signMux.HandleFunc("/re", app.HResetPassword)
	signMux.HandleFunc("/out", app.HLogout)
	signMux.HandleFunc("/status", app.HCheckUserLogged)
	appMux.Handle("/sign/", http.StripPrefix("/sign", signMux))

	// api routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", app.HApiIndex)
	apiMux.HandleFunc("/user", app.HUser)
	apiMux.HandleFunc("/users", app.HUsers)
	apiMux.HandleFunc("/parsels", app.HParsels)
	apiMux.HandleFunc("/travelers", app.HTravelers)
	apiMux.HandleFunc("/toptypes", app.HTopTypes)
	apiMux.HandleFunc("/travelTypes", app.HTravelTypes)
	apiMux.HandleFunc("/countryCodes", app.HCountryCodes)
	apiMux.HandleFunc("/search", app.HSearch)
	apiMux.HandleFunc("/images", app.HClippedImages)
	appMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// save
	saveMux := http.NewServeMux()
	saveMux.HandleFunc("/", app.HApiIndex)
	saveMux.HandleFunc("/parsel", app.HSaveParsel)
	saveMux.HandleFunc("/travel", app.HSaveTravel)
	saveMux.HandleFunc("/image", app.HSaveImage)
	appMux.Handle("/s/", http.StripPrefix("/s", saveMux))

	// edit
	editMux := http.NewServeMux()
	editMux.HandleFunc("/", app.HApiIndex)
	editMux.HandleFunc("/user/confirm", app.HPreChangeProfileSMS)
	editMux.HandleFunc("/user", app.HChangeProfile)
	editMux.HandleFunc("/parsel", app.HChangeParsel)
	editMux.HandleFunc("/travel", app.HChangeTravel)
	editMux.HandleFunc("/toptype", app.HChangeTop)
	editMux.HandleFunc("/up", app.HItemUp)
	appMux.Handle("/e/", http.StripPrefix("/e", editMux))

	// remove
	removeMux := http.NewServeMux()
	removeMux.HandleFunc("/", app.HApiIndex)
	removeMux.HandleFunc("/parsel", app.HRemoveParsel)
	removeMux.HandleFunc("/travel", app.HRemoveTravel)
	removeMux.HandleFunc("/image", app.HRemoveImage)
	appMux.Handle("/r/", http.StripPrefix("/r", removeMux))

	// static react get
	static := http.FileServer(http.Dir("website/static"))
	appMux.Handle("/static/", http.StripPrefix("/static/", static))

	// assets get
	assets := http.FileServer(http.Dir("assets"))
	appMux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// middlewares
	muxHanlder := app.AccessLogMiddleware(appMux)
	muxHanlder = app.SecureHeaderMiddleware(muxHanlder)
	return muxHanlder
}
