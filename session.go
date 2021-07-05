package main

import (
	"github.com/rs/xid"
	"gopkg.in/olahol/melody.v1"
	"sync"
)

var (
	sessionMap  = make(map[string]*Session)
	sessionLock sync.RWMutex
)

type Session struct {
	Namespace string          `json:"namespace"`
	Id        string          `json:"id"`
	Ip        string          `json:"ip"`
	Ws        *melody.Session `json:"-"`
}

func NewSession(namespace string, ip string, m *melody.Session) (string, *Session) {
	var id = xid.New().String()
	var s = &Session{
		Namespace: namespace,
		Id:        id,
		Ip:        ip,
		Ws:        m,
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessionMap[id] = s
	return id, s
}

func UnmapSession(id string) {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	delete(sessionMap, id)
}
