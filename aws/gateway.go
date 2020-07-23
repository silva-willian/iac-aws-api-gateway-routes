package aws

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/silva-willian/iac-aws-api-gateway-routes/api"
)

var apiID string = ""
var gatewayName = os.Getenv("GATEWAY_NAME")
var basePath string = os.Getenv("GATEWAY_BASE_PATH")
var gatewayResources []GatewayResource = []GatewayResource{}
var baseHost string = ""

// CreateAllResources is the function that takes the feature list and creates it in AWS
func CreateAllResources(resources []api.Resource, host string) error {
	err := getRestAPIS()

	if err != nil {
		return err
	}

	if apiID == "" {
		return fmt.Errorf("Unidentified Gateway API")
	}

	baseHost = host
	itens, err := getAllResources()

	if err != nil {
		return err
	}

	gatewayResources = itens

	err = deleteBaseResourceIFExists()

	if err != nil {
		return err
	}

	err = createResource(basePath, "/")

	if err != nil {
		return err
	}

	for _, resource := range resources {
		err := createResourceAndDependencies(resource)

		if err != nil {
			return err
		}
	}

	err = createDeploy()

	if err != nil {
		return err
	}

	return nil
}

func getGatewayClient() (*apigateway.APIGateway, error) {
	sess, err := getSession()

	if err != nil {
		return nil, handleAWSError(err)
	}

	return apigateway.New(sess), nil
}

func createResourceAndDependencies(resource api.Resource) error {

	paths := strings.Split(resource.Path, "/")
	basePathItem := fmt.Sprintf("/%s", basePath)

	for _, path := range paths {

		if path == "" {
			continue
		}

		err := createResource(path, basePathItem)
		basePathItem = fmt.Sprintf("%s/%s", basePathItem, path)

		if err != nil {
			log.Printf("Error on created API Gateway [%s] resource", basePathItem)
			return err
		}
	}

	log.Printf("Successfully created API Gateway [%s] resource", basePathItem)

	resourceID, err := getResourceID(fmt.Sprintf("/%s%s", basePath, resource.Path))

	if err != nil {
		log.Printf("Error finding ResourceID from path [%s]", basePathItem)
		return err
	}

	authorizerID, err := getAuthorizer()

	if err != nil {
		log.Printf("Error retrieving gateway authorizer [%s]", basePathItem)
		return err
	}

	for _, method := range resource.Methods {
		err = createMethod(resourceID, method.Verb, authorizerID, method.Parameters)

		if err != nil {
			log.Printf("Error creating method [%s] for route [%s]", method.Verb, basePathItem)
			return err
		}

		err = createMethodIntegration(resourceID, method.Verb, fmt.Sprintf("%s%s", baseHost, resource.Path), method.Parameters)

		if err != nil {
			log.Printf("Error creating method [%s] integration for route [%s]", method.Verb, basePathItem)

			return err
		}

		for index, statusCode := range method.Status {

			partner := statusCode

			if index == 0 {
				partner = "-"
			}
			err = createMethodResponse(resourceID, method.Verb, statusCode, partner)

			if err != nil {
				log.Printf("Error creating method [%s] response for route [%s] with status [%s]", method.Verb, basePathItem, statusCode)
				return err
			}
		}

	}

	log.Printf("Success in creating all [%s] feature dependencies in API Gateway", basePathItem)

	return nil
}

func getRestAPIS() error {
	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	result, err := client.GetRestApis(nil)

	if err != nil { // resp is now filled
		return handleAWSError(err)
	}

	if len(result.Items) <= 0 {
		return fmt.Errorf("There are no APIS rest registered")
	}

	for _, item := range result.Items {

		if *item.Name == gatewayName {
			apiID = *item.Id
			return nil
		}
	}

	return nil
}
