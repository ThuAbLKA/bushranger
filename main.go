package main

import (
	"fmt"
	"net/http"

	"github.com/ThuAbLKA/bushranger/model"
	"github.com/ThuAbLKA/bushranger/util"
	"github.com/go-redis/redis/v8"
)

func mainHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "BushRanger on duty!")
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	serviceHandler := model.NewServiceHandler(rdb)
	http.HandleFunc("/service", serviceHandler.Controller)

	nodehandler := model.NodeHandler{}
	http.HandleFunc("/node", nodehandler.Controller)

	fmt.Println("BushRanger running on port 3001...")

	http.HandleFunc("/", mainHandler)

	util.CheckError(http.ListenAndServe(":3001", nil))

}
