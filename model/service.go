package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ThuAbLKA/bushranger/util"
)

type Service struct {
	ID          string `json: "Id"`
	Description string `json: "description"`
	Nodes       []Node
}

type ServiceDao struct {
	ID          string `json: "Id"`
	Description string `json: "description"`
	ServiceName string `json: "serviceName"`
	Host        string `json: "host"`
	Port        int16  `json: "port"`
	Secure      bool   `json: "secure"`
	HealthCheck string `json: "healthCheck"`
}

// ServiceRespository
type ServiceRespository struct {
	store map[string]Service
}

// ServiceHandler
type ServiceHandler struct {
	repository ServiceRespository
}

// Controller
func (c *ServiceHandler) Controller(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		services, err := json.Marshal(c.repository.store)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(services)
	case http.MethodPost:

		var newrec ServiceDao
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
		// check whether the service is already there
		extservice := c.repository.store[newrec.ID]

		if &extservice != nil {
			// add a new node to the service
			node := Node{
				ID:           util.GenerateId(),
				ServiceName:  newrec.ServiceName,
				Host:         newrec.Host,
				Port:         newrec.Port,
				Secure:       newrec.Secure,
				HealthCheck:  newrec.HealthCheck,
				LastCheck:    time.Now(),
				FirstContact: time.Now(),
			}
			extservice.Nodes = append(extservice.Nodes, node)
		} else {
			// add new service
			service := Service{
				ID:          newrec.ServiceName,
				Description: newrec.Description,
				Nodes:       []Node{},
			}
			c.repository.store[newrec.ID] = service
		}

	}
}