package main

import (
	"github.com/gin-gonic/gin"
	"github.com/messaget/messaget/intent"
	"strconv"
	"strings"
)

func intentListClients(c *gin.Context, intent intent.Intent) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	var sil []*Session

	for s := range sessionMap {
		sil = append(sil, sessionMap[s])
	}

	if sil == nil {
		sil = make([]*Session, 0)
	}

	c.JSON(200, sil)
}

func intentFindByNamespaceExact(c *gin.Context, intent intent.Intent) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	// find targets
	var sil []*Session
	for s := range sessionMap {
		if sessionMap[s].Namespace == intent.Namespace {
			sil = append(sil, sessionMap[s])
		}
	}

	if sil == nil {
		sil = make([]*Session, 0)
	}

	c.JSON(200, sil)
}

func intentFindByNamespaces(c *gin.Context, intent intent.Intent) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	// find targets
	var sil []*Session
	for s := range sessionMap {
		if strings.Contains(sessionMap[s].Namespace, intent.Namespace) {
			sil = append(sil, sessionMap[s])
		}
	}

	if sil == nil {
		sil = make([]*Session, 0)
	}

	c.JSON(200, sil)
}

func intendFindByIds(c *gin.Context, intent intent.Intent) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	// find targets
	var sil []*Session
	for i := range intent.Targets {
		for s := range sessionMap {
			if sessionMap[s].Id == intent.Targets[i] {
				sil = append(sil, sessionMap[s])
			}
		}
	}

	if sil == nil {
		sil = make([]*Session, 0)
	}

	c.JSON(200, sil)
}

func intentKickClients(c *gin.Context, intent intent.Intent) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	// find targets
	var sil []*Session
	for i := range intent.Targets {
		for s := range sessionMap {
			if sessionMap[s].Id == intent.Targets[i] {
				sil = append(sil, sessionMap[s])
				sessionMap[s].Ws.Close()
			}
		}
	}

	if sil == nil {
		sil = make([]*Session, 0)
	}

	c.JSON(200, sil)
}

func intentSendToId(c *gin.Context, intent intent.Intent) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	// find targets
	var sil []*Session
	for i := range intent.Targets {
		for s := range sessionMap {
			if sessionMap[s].Id == intent.Targets[i] {
				sil = append(sil, sessionMap[s])
			}
		}
	}

	// send id's
	var send = 0
	for i := range sil {
		if !sil[i].Ws.IsClosed() {
			sil[i].Ws.Write([]byte(intent.Message))
			send++
		}
	}

	c.JSON(200, gin.H{
		"sent": strconv.Itoa(send),
	})
}
