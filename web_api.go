package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strconv"
)

func setupWebApi(cnf *config) {
	r := gin.Default()

	registerRoutes(cnf, r)

	if cnf.Server.UseAutoCert {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cnf.Server.Url),
			Cache:      autocert.DirCache(cnf.Server.CertPath),
		}

		server := &http.Server{
			Addr:    ":" + strconv.Itoa(cnf.Server.Port),
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go errorLogger.Fatalln(http.ListenAndServe(":"+strconv.Itoa(cnf.Server.Port), certManager.HTTPHandler(server.Handler)))
	} else {
		go errorLogger.Fatalln(r.Run(":" + strconv.Itoa(cnf.Server.Port)))
	}
}
