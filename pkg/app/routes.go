package app

import "net/http"

func (app *Application) SetRoutes() http.Handler {
	appMux := http.NewServeMux()
	appMux.HandleFunc("/", app.HIndex)

	// sign
	signMux := http.NewServeMux()
	signMux.HandleFunc("/", app.HIndex)
	signMux.HandleFunc("/sms/s", app.HPreSignUpSMS)
	signMux.HandleFunc("/up/check", app.HPreSignUpCheck)
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
	apiMux.HandleFunc("/parsels", app.HParsels)
	apiMux.HandleFunc("/travelers", app.HTravelers)
	apiMux.HandleFunc("/toptypes", app.HTopTypes)
	apiMux.HandleFunc("/search", app.HSearch)
	apiMux.HandleFunc("/images", app.HClippedImages)
	appMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// edit
	editMux := http.NewServeMux()
	editMux.HandleFunc("/", app.HApiIndex)
	editMux.HandleFunc("/user", app.HChangeProfile)
	editMux.HandleFunc("/user/confirm", app.HConfirmChangeProfile)
	editMux.HandleFunc("/parsel", app.HChangeParsel)
	editMux.HandleFunc("/travel", app.HChangeTravel)
	editMux.HandleFunc("/toptype", app.HChangeTop)
	appMux.Handle("/e/", http.StripPrefix("/e", editMux))

	// save
	saveMux := http.NewServeMux()
	saveMux.HandleFunc("/", app.HIndex)
	saveMux.HandleFunc("/parsel", app.HSaveParsel)
	saveMux.HandleFunc("/travel", app.HSaveTravel)
	saveMux.HandleFunc("/image", app.HSaveImage)
	appMux.Handle("/s/", http.StripPrefix("/s", saveMux))

	// assets get
	assets := http.FileServer(http.Dir("assets"))
	appMux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// middlewares
	muxHanlder := app.AccessLogMiddleware(appMux)
	muxHanlder = app.SecureHeaderMiddleware(muxHanlder)
	return muxHanlder
}
