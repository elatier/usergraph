package main

import (
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/franela/goreq"
	"net/http"
	"strconv"
)

type Object struct {
	Id      string        `json:"id"`
	Data    UserGraphNode `json:"data"`
	Version int           `json:"version"`
}

type UserGraphNode struct {
	User  User     `json:"user"`
	Edges []string `json:"edges"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//retrieve object from the object store
func (u *UserGraphResource) getObject(id string, response *restful.Response) (*Object, error) {
	res, err := goreq.Request{
		Uri:         u.baseUrl + id,
		Accept:      "application/json",
		ContentType: "application/json",
	}.Do()

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return nil, errors.New("500")
	}
	if res.StatusCode == 404 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Requested item could not be found.")
		return nil, errors.New("404")
	}
	if res.StatusCode != 200 {
		response.AddHeader("Content-Type", "text/plain")
		value, _ := res.Body.ToString()
		response.WriteErrorString(http.StatusInternalServerError, "Object store returned: "+
			strconv.Itoa(res.StatusCode)+" Response body:\n"+value)
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}
	obj := new(Object)
	res.Body.FromJsonTo(&obj)
	res.Body.Close()
	return obj, nil
}

//store argument object and store
func (u *UserGraphResource) updateObject(obj *Object) error {
	res, err := goreq.Request{
		Method:      "PUT",
		Body:        obj,
		Uri:         u.baseUrl + obj.Id,
		Accept:      "application/json",
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
		return errors.New(strconv.Itoa(res.StatusCode))
	}
	return nil
}

//Connect user1 to user2 by adding user2 id reference to user1 edge array and store the changes
func (u UserGraphResource) createNewConnection(user1, user2 *Object) error {
	edges := user1.Data.Edges
	relatedId := user2.Id
	for _, value := range edges {
		if value == relatedId {
			return nil
		}
	}
	user1.Data.Edges = append(edges, relatedId)
	return u.updateObject(user1)
}
