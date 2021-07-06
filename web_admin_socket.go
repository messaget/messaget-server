package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/messaget/messaget/intent"
	"gopkg.in/olahol/melody.v1"
)

var adminMelody = melody.New()

func handleControllerSocketEndpoint(c *gin.Context) {
	err := adminMelody.HandleRequest(c.Writer, c.Request)
	if err != nil {
		err := c.AbortWithError(400, err)
		if err != nil {
			errorLogger.Println(err)
			c.Abort()
			return
		}
	}
}

type socketTransactionResponse struct {
	Failed   bool        `json:"failed"`
	Id       string      `json:"transaction_id"`
	Status   int         `json:"status"`
	Response interface{} `json:"response"`
}

func setupAdminMelody() {
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
		// parse
		var i intent.Intent
		if err := json.Unmarshal(msg, &i); err != nil {
			errorLogger.Println(err)
			// nothing? I suppose
			return
		}

		// handle
		handler, err := intentHandler.GetHandler(i.Intent)
		if err != nil {
			errorLogger.Println(err)
			s.Write(serializeToBytes(socketTransactionResponse{
				Failed:   true,
				Id:       i.TransactionId,
				Status:   400,
				Response: "Unknown intent " + i.Intent,
			}))
			return
		}

		response, statusCode, err := handler(i)

		if err != nil {
			// fail transaction ID
			s.Write(serializeToBytes(socketTransactionResponse{
				Failed:   true,
				Id:       i.TransactionId,
				Status:   statusCode,
				Response: err,
			}))
		} else {
			s.Write(serializeToBytes(socketTransactionResponse{
				Failed:   false,
				Id:       i.TransactionId,
				Status:   statusCode,
				Response: response,
			}))
		}
	})
}

func serializeToBytes(i interface{}) []byte {
	bolB, _ := json.Marshal(i)
	return bolB
}

type UpdatePacket struct {
	Event  string  `json:"event"`
	Client Session `json:"client"`
}

func broadcastClientJoin(session *Session) {
	adminSessionLock.RLock()
	defer adminSessionLock.RUnlock()
	for s := range adminSessionMap {
		err := adminSessionMap[s].Ws.Write(serializeToBytes(UpdatePacket{
			Event:  "CLIENT_ADD",
			Client: *session,
		}))
		if err != nil {
			errorLogger.Println(err)
		}
	}
}

func broadcastClientLeave(session *Session) {
	adminSessionLock.RLock()
	defer adminSessionLock.RUnlock()
	for s := range adminSessionMap {
		err := adminSessionMap[s].Ws.Write(serializeToBytes(UpdatePacket{
			Event:  "CLIENT_LEAVE",
			Client: *session,
		}))
		if err != nil {
			errorLogger.Println(err)
		}
	}
}
