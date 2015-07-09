package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
)

func (u UserGraphResource) getConnectedUsers(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	_, exists := u.users[id]
	if !exists {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	users := make([]User, 0)
	for _, value := range u.edges[id] {
		users = append(users, u.users[value])
	}
	response.WriteEntity(users)
}

// PUT http://localhost:8080/users/{id}/connectedUsers/{id2}
//
func (u UserGraphResource) addConnectedUser(request *restful.Request, response *restful.Response) {
	user1 := request.PathParameter("user-id")
	user2 := request.PathParameter("dest-id")
	if len(u.users[user1].Id) == 0 || len(u.users[user2].Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	u.createNewConnection(user1, user2)
	u.createNewConnection(user2, user1)
	response.WriteHeader(http.StatusCreated)
}

func (u UserGraphResource) createNewConnection(user1, user2 string) {
	related := u.edges[user1]
	for _, value := range related {
		if value == user2 {
			return //connection exists
		}
	}
	u.edges[user1] = append(related, user2)
}

// GET http://localhost:8080/users/1
//
func (u UserGraphResource) findUser(request *restful.Request, response *restful.Response) {
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
func (u UserGraphResource) listUsers(request *restful.Request, response *restful.Response) {
	values := make([]User, len(u.users))
	i := 0
	for _, value := range u.users {
		values[i] = value
		i += 1
	}
	response.WriteEntity(values)
}

// POST http://localhost:8080/users
// <User><Name>Melissa</Name></User>
//
func (u *UserGraphResource) createUser(request *restful.Request, response *restful.Response) {
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
func (u *UserGraphResource) updateUser(request *restful.Request, response *restful.Response) {
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
func (u *UserGraphResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}
