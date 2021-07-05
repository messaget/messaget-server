package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/messaget/messaget/intent"
	"io/ioutil"
)

var intentHandler = intent.NewIntentMap()

func registerAdminIntents() {
	intentHandler.AddIntent("LIST_CLIENTS", intentListClients)
	intentHandler.AddIntent("SEND_TO_IDS", intentSendToId)
	intentHandler.AddIntent("FIND_BY_IDS", intendFindByIds)
	intentHandler.AddIntent("FIND_BY_NAMESPACE", intentFindByNamespaces)
	intentHandler.AddIntent("FIND_BY_NAMESPACE_EXACT", intentFindByNamespaceExact)
}

func handleIntentEndpoint(c *gin.Context) {
	// attempt to parse the intent
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, FailedIntent)
		return
	}

	// parse
	var i intent.Intent
	if err := json.Unmarshal(jsonData, &i); err != nil {
		c.JSON(400, FailedIntent)
		return
	}

	// handle
	handler, err := intentHandler.GetHandler(i.Intent)
	if err != nil {
		c.JSON(400, FailedIntent)
		return
	}

	handler(c, i)
}
