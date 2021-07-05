package main

import (
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

	if cnf.Server.UseAutoCert {
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
