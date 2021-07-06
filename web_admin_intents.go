package main

import (
	"github.com/gin-gonic/gin"
	"github.com/messaget/messaget/intent"
	"strconv"
	"strings"
)

func intentListClients(intent intent.Intent) (interface{}, int, error) {
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	var sil []*Session

	for s := range sessionMap {
		sil = append(sil, sessionMap[s])
	}

	if sil == nil {
		sil = make([]*Session, 0)
	}

	return sil, 200, nil
}

func intentFindByNamespaceExact(intent intent.Intent) (interface{}, int, error) {
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

	return sil, 200, nil
}

func intentFindByNamespaces(intent intent.Intent) (interface{}, int, error) {
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

	return sil, 200, nil
}

func intendFindByIds(intent intent.Intent) (interface{}, int, error) {
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

	return sil, 200, nil
}

func intentKickClients(intent intent.Intent) (interface{}, int, error) {
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

	return sil, 200, nil
}

func intentSendToId(intent intent.Intent) (interface{}, int, error) {
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

	return gin.H{
		"sent": strconv.Itoa(send),
	}, 200, nil
}
