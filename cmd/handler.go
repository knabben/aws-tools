package cmd

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
)

func CreateHTTPAPIHandler() http.Handler {
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)
	wsContainer.Add(apiV1Ws)
	apiV1Ws.Path("/api/v1")

	apiV1Ws.Route(
		apiV1Ws.POST("/app").To(handleDeploy),
	)
	return wsContainer

}

func handleDeploy(request *restful.Request, response *restful.Response) {
	IP := "teste"
	Machine := "teste"
	startAWSSession(IP, Machine)
}
