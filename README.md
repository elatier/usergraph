# usergraph
Build using "go build" and run using "./usergraph".
Swagger spec file is available at http://localhost:8085/apidocs.json by default.
Swagger UI needs to be downloaded separately from https://github.com/swagger-api/swagger-ui to view it.

Default API base URL: http://localhost:8085/users/

The service models users and their connections (graph). This service uses https://github.com/elatier/objectstore service for state storage.