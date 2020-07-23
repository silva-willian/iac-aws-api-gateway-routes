package api

// SwaggerResponse represents the http response structure of the route.
type SwaggerResponse struct {
	Parameters []SwaggerParameters    `json:"parameters"`
	Responses  map[string]interface{} `json:"responses"`
}

// SwaggerParameters represents the parameters of a route
type SwaggerParameters struct {
	Name     string `json:"name"`
	Location string `json:"in"`
	Required bool   `json:"required"`
}

// Swagger represents swagger jsons main object
type Swagger struct {
	Path map[string]interface{} `json:"paths"`
}

// Resource represents the routes themselves from an api
type Resource struct {
	Path    string
	Methods []ResourceMethod
}

// ResourceMethod represents the methods of a resources
type ResourceMethod struct {
	Verb       string
	Status     []string
	Parameters []ResourceParameters
}

// ResourceParameters represents the parameters of a resources
type ResourceParameters struct {
	Name     string
	Location string
	Required bool
}
