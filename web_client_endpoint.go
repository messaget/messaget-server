package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

var clientMelody = melody.New()

func handleClientEndpoint(c *gin.Context) {
	ns := c.Query("namespace")

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

	clientMelody.HandleConnect(func(s *melody.Session) {
		// build Session
		var query = s.Request.URL.Query()
		var namespace = query.Get("namespace")
		var remoteIp = readUserIp(s.Request)

		id, session := NewSession(namespace, remoteIp, s, false)
		s.Set("Session", session)

		infoLogger.Println("Accepted public client", id)
		broadcastClientJoin(session)
	})

	clientMelody.HandleDisconnect(func(s *melody.Session) {
		// get Session
		var r, _ = s.Get("Session")
		var session = r.(*Session)

		infoLogger.Println("Dropping public client", session.Id)

		UnmapSession(session.Id, false)
		broadcastClientLeave(session)
	})

	clientMelody.HandleMessage(func(s *melody.Session, msg []byte) {
		// TODO: Make this work, somehow
	})
}

func readUserIp(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
