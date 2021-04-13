package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ThuAbLKA/bushranger/util"
	"github.com/go-redis/redis/v8"
)

type Service struct {
	ID          string `json: "Id"`
	Description string `json: "description"`
	Nodes       []Node `json: "nodes"`
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
	dbClient   *redis.Client
	repository ServiceRespository
}

func NewServiceHandler(db *redis.Client) *ServiceHandler {
	handler := new(ServiceHandler)
	handler.repository.store = make(map[string]Service)
	handler.dbClient = db

	return handler
}

func (c *Service) AddNode(node Node) []Node {
	c.Nodes = append(c.Nodes, node)
	return c.Nodes
}

// Controller
func (c *ServiceHandler) Controller(res http.ResponseWriter, req *http.Request) {
	var ctx = context.Background()

	switch req.Method {
	case http.MethodGet:
		//var services []Service

		iter := c.dbClient.Scan(ctx, 0, "SER*", 0).Iterator()
		for iter.Next(ctx) {
			fmt.Println(iter.Val())
			var tempService Service
			result, _ := c.dbClient.Get(ctx, iter.Val()).Result()
			err := json.Unmarshal([]byte(result), &tempService)
			//append(services, tempService)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(""))
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
		var extservice Service

		// check whether the service is already there
		p, err := c.dbClient.Get(ctx, "SER-"+newrec.ID).Result()

		if err != redis.Nil {
			err = json.Unmarshal([]byte(p), &extservice)

			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

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

		if extservice.ID == "" {
			// add new service
			service := Service{
				ID:          "SER-" + newrec.ServiceName,
				Description: newrec.Description,
				Nodes:       []Node{node},
			}
			c.repository.store[newrec.ID] = service
			// store in the RDB
			obj, err := json.Marshal(service)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			c.dbClient.Set(ctx, "SER-"+newrec.ID, string(obj), 0)

		} else {
			extservice.Nodes = extservice.AddNode(node)
			// save to the DB
			newobj, _ := json.Marshal(extservice)
			c.dbClient.Set(ctx, "SER-"+newrec.ID, string(newobj), 0)

		}

	}
}
