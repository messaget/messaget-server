package main

import (
	"github.com/rs/xid"
	"gopkg.in/olahol/melody.v1"
	"sync"
)

var (
	sessionMap  = make(map[string]*Session)
	sessionLock sync.RWMutex

	adminSessionMap  = make(map[string]*Session)
	adminSessionLock sync.RWMutex
)

type Session struct {
	Namespace    string          `json:"namespace"`
	Id           string          `json:"id"`
	Ip           string          `json:"ip"`
	IsController bool            `json:"isController"`
	Ws           *melody.Session `json:"-"`
}

func NewSession(namespace string, ip string, m *melody.Session, isController bool) (string, *Session) {
	var id = xid.New().String()
	var s = &Session{
		Namespace:    namespace,
		Id:           id,
		Ip:           ip,
		Ws:           m,
		IsController: isController,
	}

	if isController {
		adminSessionLock.Lock()
		defer adminSessionLock.Unlock()
		adminSessionMap[id] = s
	} else {
		sessionLock.Lock()
		defer sessionLock.Unlock()
		sessionMap[id] = s
	}
	return id, s
}

func UnmapSession(id string, isController bool) {
	if isController {
		adminSessionLock.Lock()
		defer adminSessionLock.Unlock()
		delete(adminSessionMap, id)
	} else {
		sessionLock.Lock()
		defer sessionLock.Unlock()
		delete(sessionMap, id)
	}
}
