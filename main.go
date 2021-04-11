package main

import (
	"fmt"
	"net/http"

	"github.com/ThuAbLKA/bushranger/model"
	"github.com/ThuAbLKA/bushranger/util"
)

func mainHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "BushRanger on duty!")
}

func main() {
	serviceHandler := model.ServiceHandler{}
	http.HandleFunc("/service", serviceHandler.Controller)

	nodehandler := model.NodeHandler{}
	http.HandleFunc("/node", nodehandler.Controller)

	fmt.Println("BushRanger running on port 3001...")

	http.HandleFunc("/", mainHandler)

	util.CheckError(http.ListenAndServe(":3001", nil))

}
