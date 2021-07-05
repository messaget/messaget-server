package main

import (
	"github.com/rs/xid"
	"gopkg.in/olahol/melody.v1"
	"sync"
)

var (
	guid = xid.New()
	sessionMap = make(map[string]*session)
	sessionLock sync.RWMutex
)

type session struct {
	namespace string
	id        string
	ip        string
	ws        *melody.Session
}

func NewSession(namespace string, ip string, m *melody.Session) (string, *session)  {
	var id = guid.String();
	var s = &session{
		namespace: namespace,
		id: id,
		ip: ip,
		ws: m,
	}
	sessionLock.Lock()
	defer sessionLock.Unlock()
	sessionMap[id] = s
	return id, s
}

func UnmapSession(id string)  {
	sessionLock.Lock()
	defer sessionLock.Unlock()
	delete(sessionMap, id)
}