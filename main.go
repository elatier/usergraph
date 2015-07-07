package main

import (
    "log"
    "net/http"
    "strconv"

    "github.com/emicklei/go-restful"
    "github.com/emicklei/go-restful/swagger"
)

type User struct {
    Id string `json:"id"`
    Name string `json:"name"`
}

type UserResource struct {
    // normally one would use DAO (data access object)
    users map[string]User
}

func (u UserResource) Register(container *restful.Container) {
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
        Writes(User{})) // on the response

    ws.Route(ws.GET("").To(u.listUsers).
        // docs
        Doc("get a user").
        Operation("findUser").
        Writes([]User{})) // on the response

    ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
        // docs
        Doc("update a user").
        Operation("updateUser").
        Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
        ReturnsError(409, "duplicate user-id", nil).
        Reads(User{})) // from the request

    ws.Route(ws.POST("").To(u.createUser).
        // docs
        Doc("create a user").
        Operation("createUser").
        Reads(User{})) // from the request

    ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
        // docs
        Doc("delete a user").
        Operation("removeUser").
        Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

    container.Add(ws)
}

// GET http://localhost:8080/users/1
//
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("user-id")
    usr := u.users[id]
    if len(usr.Id) == 0 {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
        return
    }
    response.WriteEntity(usr)
}

// GET http://localhost:8080/users/
//
func (u UserResource) listUsers(request *restful.Request, response *restful.Response) {
    values := make([]User, len(u.users))
    i := 0
    for _,value := range u.users {
        values[i] = value
        i += 1
    }
    response.WriteEntity(values)
}
// POST http://localhost:8080/users
// <User><Name>Melissa</Name></User>
//
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
    usr := new(User)
    err := request.ReadEntity(usr)
    if err != nil {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    usr.Id = strconv.Itoa(len(u.users) + 1) // simple id generation
    u.users[usr.Id] = *usr
    response.WriteHeader(http.StatusCreated)
    response.WriteEntity(usr)
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("user-id")
    usr := new(User)
    err := request.ReadEntity(&usr)
    usr.Id = id
    if err != nil {
        response.AddHeader("Content-Type", "text/plain")
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    u.users[usr.Id] = *usr

    response.WriteEntity(usr)
}

// DELETE http://localhost:8080/users/1
//
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("user-id")
    delete(u.users, id)
}

func main() {
    // to see what happens in the package, uncomment the following
    //restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

    wsContainer := restful.NewContainer()
    userRes := UserResource{map[string]User{}}
    userRes.Register(wsContainer)

    // Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
    // You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
    // Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
    config := swagger.Config{
        WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
        WebServicesUrl: "http://localhost:8080",
        ApiPath:        "/apidocs.json",

        // Optionally, specifiy where the UI is located
        SwaggerPath:     "/apidocs/",
        SwaggerFilePath: "/Users/kriaval/developer/swagger-ui"}
    swagger.RegisterSwaggerService(config, wsContainer)

    log.Printf("start listening on localhost:8080")
    server := &http.Server{Addr: ":8080", Handler: wsContainer}
    log.Fatal(server.ListenAndServe())
}