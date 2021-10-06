package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"alber/pkg/app"
)

func main() {
	app := app.InitProg()

	app.ILog.Println("initialization completed!")

	// check sessions expire per minute
	go app.CheckPerMin()

	// server
	srv := http.Server{
		Addr:         ":" + app.Config.PORT,
		ErrorLog:     app.ELog,
		Handler:      app.SetRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305},
		},
	}

	fmt.Printf("server listening on port %v\n", app.Config.PORT)
	app.ILog.Printf("server listening on port %v\n", app.Config.PORT)

	// HTTP
	// app.ELog.Fatal(srv.ListenAndServe())

	// HTTPS
	app.ELog.Fatal(srv.ListenAndServeTLS("./tls/al-ber_kz.crt", "./tls/11029176.key"))
}
