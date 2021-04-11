package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ThuAbLKA/bushranger/util"
)

type Node struct {
	ID           string    `json:"id"`
	ServiceName  string    `json: "serviceName"`
	Host         string    `json: "host"`
	Port         int16     `json: "port"`
	Secure       bool      `json: "secure"`
	HealthCheck  string    `json: "healthCheck"`
	LastCheck    time.Time `json: "lastCheck"`
	FirstContact time.Time `json: "firstContact"`
}

// NodeRepository
type NodeRepository struct {
	store map[string]Node
}

// NodeHandler
type NodeHandler struct {
	repository NodeRepository
}

// Controller
func (c *NodeHandler) Controller(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		nodes, err := json.Marshal(c.repository.store)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(nodes)
	case http.MethodPost:
		var newrec Node
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &newrec)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		newrec.ID = util.GenerateId()
		c.repository.store[newrec.ID] = newrec

	}
}
