package cmd

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful"
)

type SSHTuple struct {
	IP, Machine string
}

func CreateHTTPAPIHandler() http.Handler {
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)
	wsContainer.Add(apiV1Ws)
	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiV1Ws.Route(apiV1Ws.POST("/add").To(handleAddIP).Reads(SSHTuple{}))
	apiV1Ws.Route(apiV1Ws.DELETE("/").To(handleDelIP))
	return wsContainer
}

// TODO  - Broke validation
func validadeEntity(ipTuple *SSHTuple, request *restful.Request, response *restful.Response) *restful.Response {
	err := request.ReadEntity(&ipTuple)
	if err != nil {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
	return response
}

func isValid(response *restful.Response) bool {
	return response.StatusCode() == 200
}

func handleAddIP(request *restful.Request, response *restful.Response) {
	ipTuple := new(SSHTuple)
	response = validadeEntity(ipTuple, request, response)
	if !isValid(response) {
		return
	}

	ip, machine := ipTuple.IP, ipTuple.Machine

	err := client.Set(ip+"|"+machine, "true", time.Month).Err()
	if err != nil {
		panic(err)
	}
	inserted := InsertIPOnSG(ip, machine)
	if inserted {
		response.WriteHeader(http.StatusCreated)
	} else {
		response.WriteHeader(http.StatusConflict)
	}
}

func handleDelIP(request *restful.Request, response *restful.Response) {
	ipTuple := new(SSHTuple)
	response = validadeEntity(ipTuple, request, response)
	if !isValid(response) {
		return
	}

	deleted := deleteIPFromSG(ipTuple.IP, ipTuple.Machine)
	if !deleted {
		response.WriteHeader(http.StatusNotFound)
	}
}
