package main

import (
	"github.com/emicklei/go-restful"
	"github.com/franela/goreq"
	"net/http"
	"strconv"
)

// GET http://localhost:8085/users/{id}
//
func (u UserGraphResource) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	obj, err := u.getObject(id, response)
	if err != nil {
		return
	}
	user := obj.Data.User
	user.Id = obj.Id
	response.WriteEntity(user)
}

// POST http://localhost:8085/users
// {"name":"user name"}
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
		Method:      "POST",
		Body:        obj,
		Uri:         u.baseUrl,
		Accept:      "application/json",
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
		response.WriteErrorString(http.StatusInternalServerError, "Object store returned: "+
			strconv.Itoa(res.StatusCode)+" Response body:\n"+value)
		return
	}
	res.Body.FromJsonTo(&obj)
	res.Body.Close()
	user := obj.Data.User
	user.Id = obj.Id
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(user)
}

// PUT http://localhost:8085/users/{id}
// {"name":"user name"}
//
func (u *UserGraphResource) updateUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	obj, err := u.getObject(id, response)
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
	err = u.updateObject(obj)
	response.WriteEntity(obj.Data.User)
}

// GET http://localhost:8085/users/{id}/connectedUsers
//
func (u UserGraphResource) getConnectedUsers(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	obj, err := u.getObject(id, response)
	if err != nil {
		return
	}
	users := make([]User, 0)
	for _, edgeUserId := range obj.Data.Edges {
		relatedObj, err := u.getObject(edgeUserId, response)
		if err != nil {
			return
		}
		relatedUser := relatedObj.Data.User
		relatedUser.Id = relatedObj.Id
		users = append(users, relatedUser)

	}
	response.WriteEntity(users)
}

// PUT http://localhost:8085/users/{id}/connectedUsers/{id2}
//
func (u UserGraphResource) addConnectedUser(request *restful.Request, response *restful.Response) {
	user1 := request.PathParameter("user-id")
	user2 := request.PathParameter("dest-id")
	if user1 == user2 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Source and destination ID is indentical")
		return
	}
	obj, err := u.getObject(user1, response)
	obj2, err2 := u.getObject(user2, response)
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