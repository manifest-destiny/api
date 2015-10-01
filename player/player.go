package player

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

type Player struct {
	Id, GoogleId, Name, Email, Picture string
}

type PlayerResource struct {
	// Store *store.Store
}

func (p PlayerResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/p").
		Doc("Manage plalyers").
		Consumes(restful.MIME_JSON, "application/x-www-form-urlencoded").
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/tokensignin").
		To(p.tokenSignIn).
		Doc("sign in").
		Consumes("application/x-www-form-urlencoded").
		Produces(restful.MIME_JSON).
		Operation("tokenSignIn").
		Param(ws.BodyParameter("signin-token", "User's sign in token").DataType("string")).
		Writes(struct{ Message string }{}))

	container.Add(ws)
}

func (p PlayerResource) tokenSignIn(request *restful.Request, response *restful.Response) {
	token, err := request.BodyParameter("signin-token")
	if err != nil {
		writeMessage(response, http.StatusBadRequest, "Missing signin-token")
		return
	}
	e := struct{ Message string }{Message: token}
	response.WriteEntity(e)
}

// func (p PlayerResource) createGame(request *restful.Request, response *restful.Response) {
// 	p, err := uthPlayer(request)
// }

// writeMessage Convenience function for sending http error (and other) messages
// in JSON format
func writeMessage(res *restful.Response, httpStatus int, msg string) {
	e := struct{ Message string }{Message: msg}
	res.WriteHeader(httpStatus)
	res.WriteEntity(e)
}
