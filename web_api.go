package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strconv"
)

func setupWebApi(cnf *config) {
	if cnf.Logging.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	registerRoutes(cnf, r)

	if cnf.Server.UseManualCert {

		var x509 tls.Certificate
		x509, err := tls.LoadX509KeyPair(cnf.Server.ManualFullChain, cnf.Server.ManualPrivate)
		if err != nil {
			errorLogger.Fatalln(err)
		}

		server := &http.Server{
			Addr:      ":" + strconv.Itoa(cnf.Server.Port),
			Handler:   r,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{ x509 },
			},
		}
		go errorLogger.Fatalln(server.ListenAndServe())
	} else if cnf.Server.UseAutoCert {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cnf.Server.Url),
			Cache:      autocert.DirCache(cnf.Server.CertPath),
		}

		server := &http.Server{
			Addr:      ":" + strconv.Itoa(cnf.Server.Port),
			Handler:   r,
			TLSConfig: certManager.TLSConfig(),
		}
		go errorLogger.Fatalln(server.ListenAndServeTLS("", ""))
		// go errorLogger.Fatalln(http.l(":"+strconv.Itoa(cnf.Server.Port), certManager.HTTPHandler(server.Handler)))
	} else {
		go errorLogger.Fatalln(r.Run(":" + strconv.Itoa(cnf.Server.Port)))
	}
}
