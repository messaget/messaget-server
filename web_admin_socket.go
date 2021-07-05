package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

var adminMelody = melody.New()

func handleControllerSocketEndpoint(c *gin.Context) {
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

func setupAdminMelody()  {
	adminMelody.HandleError(func(session *melody.Session, err error) {
		errorLogger.Println(err)
	})

	adminMelody.HandleConnect(func(s *melody.Session) {
		// build Session
		var remoteIp = readUserIp(s.Request)

		id, session := NewSession("controller", remoteIp, s, true)
		s.Set("Session", session)

		infoLogger.Println("Accepted controller client", id)
	})

	adminMelody.HandleDisconnect(func(s *melody.Session) {
		// get Session
		var r, _ = s.Get("Session")
		var session = r.(*Session)

		infoLogger.Println("Dropping controller client", session.Id)
		UnmapSession(session.Id, true)
	})

	adminMelody.HandleMessage(func(s *melody.Session, msg []byte) {
		// TODO: Make this work, somehow
	})
}

func serializeToBytes(i interface{}) []byte  {
	bolB, _ := json.Marshal(i)
	return bolB
}

func broadcastClientJoin(session *Session)  {
	adminSessionLock.RLock()
	defer adminSessionLock.RUnlock()
	for s := range adminSessionMap {
		adminSessionMap[s].Ws.Write(serializeToBytes(gin.H{
			"event": "CLIENT_ADD",
			"client": sessionMap,
		}))
	}
}

func broadcastClientLeave(session *Session)  {
	adminSessionLock.RLock()
	defer adminSessionLock.RUnlock()
	for s := range adminSessionMap {
		adminSessionMap[s].Ws.Write(serializeToBytes(gin.H{
			"event": "CLIENT_LEAVE",
			"client": sessionMap,
		}))
	}
}
