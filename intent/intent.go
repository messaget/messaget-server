package intent

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context, intent Intent)

type Intent struct {
	Intent  string   `json:"intent"`
	Targets []string `json:"targets"`
	Message string   `json:"message"`
}

type IntentHandler struct {
	Name    string
	Handler *Handler
}

type IntentMap struct {
	intents map[string]*IntentHandler
}

func NewIntentMap() IntentMap {
	return IntentMap{
		intents: make(map[string]*IntentHandler),
	}
}

func (im *IntentMap) AddIntent(name string, h Handler) {
	im.intents[name] = &IntentHandler{
		Name:    name,
		Handler: &h,
	}
}

func (im *IntentMap) GetHandler(name string) (Handler, error) {
	h, found := im.intents[name]

	if !found {
		return nil, errors.New("Unknown intent")
	}

	return *h.Handler, nil
}