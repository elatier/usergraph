package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

type UserGraphResource struct {
	baseUrl string
}

func (u UserGraphResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Doc("Manage Users").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		ReturnsError(404, "User could not be found", nil).
		Writes(User{})) // on the response

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		ReturnsError(404, "User could not be found", nil).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").To(u.createUser).
		// docs
		Doc("create a user").
		Operation("createUser").
		Returns(201, "User creted", User{}).
		Reads(User{})) // from the request

	ws.Route(ws.GET("/{user-id}/connectedUsers").To(u.getConnectedUsers).
		// docs
		Doc("get the of list connected users").
		Operation("getConnectedUsers").
		Param(ws.PathParameter("user-id", "identifier of the source user").DataType("string")).
		ReturnsError(404, "User could not be found", nil).
		Writes([]User{})) // on the response

	ws.Route(ws.PUT("/{user-id}/connectedUsers/{dest-id}").To(u.addConnectedUser).
		// docs
		Doc("add a connected user relation").
		Operation("addConnectedUser").
		Param(ws.PathParameter("user-id", "identifier of the source user").DataType("string")).
		Param(ws.PathParameter("dest-id", "identifier of the destination user").DataType("string")).
		ReturnsError(400, "duplicate user-id", nil).
		ReturnsError(404, "User could not be found", nil))

	container.Add(ws)
}

func main() {
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))
	ug := UserGraphResource{"http://localhost:8090/tables/usergraph/objects/"}
	wsContainer := restful.NewContainer()

	ug.Register(wsContainer)

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8085",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/kriaval/developer/swagger-ui"}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:8085")
	server := &http.Server{Addr: ":8085", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
