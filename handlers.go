package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"github.com/franela/goreq"
	"fmt"
	"strconv"
	"errors"
)

type Object struct {
	Id   string `json:"id"`
	Data UserGraphNode`json:"data"`
	Version int `json:"version"`
}

type UserGraphNode struct {
	User User `json:"user"`
	Edges []string `json:"edges"`
}

func (u UserGraphResource) getConnectedUsers(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	obj, err := getObject(id, response)
	if err != nil {
		return
	}
	users := make([]User, 0)
	for _, edgeUserId := range obj.Data.Edges {
		relatedObj, err := getObject(edgeUserId, response)
		if err != nil {
			return
		}
		relatedUser :=  relatedObj.Data.User
		relatedUser.Id = relatedObj.Id
		users = append(users, relatedUser)

	}
	response.WriteEntity(users)
}

// PUT http://localhost:8080/users/{id}/connectedUsers/{id2}
//
func (u UserGraphResource) addConnectedUser(request *restful.Request, response *restful.Response) {
	user1 := request.PathParameter("user-id")
	user2 := request.PathParameter("dest-id")
	if user1==user2 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Source and destination ID is indentical")
		return
	}
	obj, err := getObject(user1, response)
	obj2, err2 := getObject(user2, response)
	if err != nil || err2 != nil {
		return
	}
	err = u.createNewConnection(obj, obj2)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	err2 = u.createNewConnection(obj2, obj)
	if err2 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteHeader(http.StatusCreated)
}

func (u UserGraphResource) createNewConnection(user1, user2 *Object) error {
	edges := user1.Data.Edges
	relatedId := user2.Id
	for _, value := range edges {
		if value == relatedId {
			return nil
		}
	}
	user1.Data.Edges = append(edges, relatedId)
	return updateObject(user1)
}

func updateObject(obj *Object) (error) {
	res, err := goreq.Request{
		Method: "PUT",
		Body: obj,
		Uri: "http://localhost:8090/tables/usergraph/objects/"+obj.Id,
		Accept: "application/json",
		ContentType: "application/json",
	}.Do()

	if err != nil {
		return err
	}
	res.Body.FromJsonTo(&obj)
	res.Body.Close()
	if res.StatusCode == 409 {
		return errors.New("409")
	}
	if res.StatusCode != 200 {
		return  errors.New(strconv.Itoa(res.StatusCode))
	}
	return nil
}


// GET http://localhost:8080/users/1
//
func (u UserGraphResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	obj, _ := getObject(id, response)
	if obj != nil {
		user :=  obj.Data.User
		user.Id = obj.Id
		response.WriteEntity(user)
	}
}

func getObject(id string, response *restful.Response) (*Object, error) {
	res, err := goreq.Request{
		Uri: "http://localhost:8090/tables/usergraph/objects/"+id,
		Accept: "application/json",
		ContentType: "application/json",
	}.Do()

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return nil, errors.New("500")
	}
	if res.StatusCode == 404 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Requested item could not be found.")
		return nil, errors.New("404")
	}
	if res.StatusCode != 200 {
		response.AddHeader("Content-Type", "text/plain")
		value, _ := res.Body.ToString()
		response.WriteErrorString(http.StatusInternalServerError,"Object store returned: "+
		strconv.Itoa(res.StatusCode) +" Response body:\n"+value)
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
	obj := new(Object)
	res.Body.FromJsonTo(&obj)
	res.Body.Close()
	return obj, nil
}

// POST http://localhost:8080/users
// <User><Name>Melissa</Name></User>
//
func (u *UserGraphResource) createUser(request *restful.Request, response *restful.Response) {
	obj := new(Object)
	err := request.ReadEntity(&obj.Data.User)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	res, err := goreq.Request{
		Method: "POST",
		Body: obj,
		Uri: "http://localhost:8090/tables/usergraph/objects/",
		Accept: "application/json",
		ContentType: "application/json",
	}.Do()

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	if res.StatusCode != 201 {
		response.AddHeader("Content-Type", "text/plain")
		value, _ := res.Body.ToString()
		response.WriteErrorString(http.StatusInternalServerError,"Object store returned: "+
				strconv.Itoa(res.StatusCode)+" Response body:\n"+value)
		return
	}
	res.Body.FromJsonTo(&obj)
	res.Body.Close()
	fmt.Printf("\n%+v",obj)
	user :=  obj.Data.User
	user.Id = obj.Id
	response.WriteEntity(user)
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserGraphResource) updateUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	obj, err := getObject(id, response)
	if err != nil {
		return
	}
	err = request.ReadEntity(&obj.Data.User)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	obj.Data.User.Id = obj.Id
	err = updateObject(obj)
	response.WriteEntity(obj.Data.User)
}