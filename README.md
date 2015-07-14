# usergraph
Build using *go build* and run using *./usergraph*.
Swagger spec file is available at http://localhost:8085/apidocs.json by default.
Swagger UI needs to be downloaded separately from https://github.com/swagger-api/swagger-ui to view it.

Default API base URL: http://localhost:8085/users/

The service models users and their connections to each other (graph). 
This service uses https://github.com/elatier/objectstore service for state storage. The data structures used:

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

- Edges - stores use connections as an array of related User Ids
- UserGraphNode - service specific data structure, making use of customizable ObjectStore Data field.