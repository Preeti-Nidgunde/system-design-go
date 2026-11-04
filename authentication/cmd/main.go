package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/Preeti-Nidgunde/system-design-go/authentication"
)

func main() {
	authsvc := authentication.NewAuthService()
	handler := authentication.NewLoginHandler(authsvc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("POST /user", handler.CreateUser)

	srv := &http.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
