package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

var clientMelody = melody.New()

func handleClientEndpoint(c *gin.Context) {
	ns, found := c.Params.Get("namespace")

	if !found {
		err := c.AbortWithError(400, NoNamespaceError)
		if err != nil {
			errorLogger.Println(err)
			c.Abort()
			return
		}
	}

	if len(ns) > cnf.Auth.MaxNamespaceLength {
		err := c.AbortWithError(400, NamespaceTooLong)
		if err != nil {
			errorLogger.Println(err)
			c.Abort()
			return
		}
	}

	err := clientMelody.HandleRequest(c.Writer, c.Request)
	if err != nil {
		err := c.AbortWithError(400, err)
		if err != nil {
			errorLogger.Println(err)
			c.Abort()
			return
		}
	}
}

func setupMelody()  {
	clientMelody.HandleError(func(session *melody.Session, err error) {
		errorLogger.Println(err)
	})

	clientMelody.HandleDisconnect(func(session *melody.Session) {

	})

	clientMelody.HandleMessage(func(s *melody.Session, msg []byte) {

	})
}
