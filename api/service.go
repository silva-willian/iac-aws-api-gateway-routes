package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// GetResources is the function that recovers api resources from swagger json
func GetResources(host string) ([]Resource, error) {
	log.Printf("Querying for host [%s]", host)
	swagger, err := getSwaggerJSON(host)

	if err != nil {
		log.Println("Error retrieving Swagger JSON")
		return nil, err
	}

	resources, err := parseSwaggerJSON(swagger)

	if err != nil {
		log.Println("Error converting Swagger JSON")
		return nil, err
	}

	log.Println("Swagger Json successfully converted")
	return resources, nil
}

func getSwaggerJSON(host string) (Swagger, error) {

	client := http.Client{}
	request, err := http.NewRequest("GET", host, nil)

	res, err := client.Do(request)

	if err != nil {
		log.Panicf("Error in GET %s %v", host, err)
		return Swagger{}, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panicf("ReadError in GET %s %v", host, err)
		return Swagger{}, nil
	}

	var result Swagger
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Panicf("JsonError in GET %s %v", host, err)
		return Swagger{}, err
	}

	return result, nil
}

func parseSwaggerJSON(swagger Swagger) ([]Resource, error) {
	resources := []Resource{}

	for path, content := range swagger.Path {
		resource := Resource{}
		resource.Path = path
		pathContent, err := content.(map[string]interface{})

		if err == false {
			return nil, fmt.Errorf("Error converting SwaggerJson path structure")
		}
		resourceMethodOptions := ResourceMethod{
			Verb:   "OPTIONS",
			Status: []string{"200"},
		}
		resource.Methods = append(resource.Methods, resourceMethodOptions)
		for method, methodContent := range pathContent {
			resourceMethod := ResourceMethod{}

			resourceMethod.Verb = strings.ToUpper(method)

			result, parameters, err := parseMethodSwaggerJSON(methodContent)

			if err != nil {
				return nil, err
			}
			resourceMethod.Parameters = parameters
			resourceMethod.Status = result
			resource.Methods = append(resource.Methods, resourceMethod)

		}
		resources = append(resources, resource)

	}

	return resources, nil
}

func parseMethodSwaggerJSON(content interface{}) ([]string, []ResourceParameters, error) {

	jsonbody, err := json.Marshal(content)

	if err != nil {
		return nil, nil, err
	}

	swaggerResponse := SwaggerResponse{}
	err = json.Unmarshal(jsonbody, &swaggerResponse)

	if err != nil {
		return nil, nil, err
	}

	if swaggerResponse.Responses == nil || len(swaggerResponse.Responses) <= 0 {
		return nil, nil, fmt.Errorf("Error converting swagger json response http codes")
	}

	statusCode := []string{}

	for status := range swaggerResponse.Responses {
		statusCode = append(statusCode, status)
	}

	resourceParameters := []ResourceParameters{}

	for _, parameter := range swaggerResponse.Parameters {

		if parameter.Location == "body" || parameter.Location == "formData" {
			continue
		}

		if parameter.Location == "query" {
			parameter.Location = "querystring"
		}

		resourceParameters = append(resourceParameters, ResourceParameters{
			Name:     parameter.Name,
			Location: parameter.Location,
			Required: parameter.Required,
		})
	}

	return statusCode, resourceParameters, nil
}
