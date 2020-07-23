package main

import (
	"fmt"
	"log"
	"os"

	"github.com/silva-willian/iac-aws-api-gateway-routes/api"
	"github.com/silva-willian/iac-aws-api-gateway-routes/aws"
	"github.com/silva-willian/iac-aws-api-gateway-routes/utils"
)

var apiHost string = os.Getenv("API_HOST")
var swaggerJSONPath string = os.Getenv("API_SWAGGER_JSON_PATH")
var gatewayName string = os.Getenv("GATEWAY_NAME")
var gatewayStage string = os.Getenv("GATEWAY_STAGE")
var gatewayPath string = os.Getenv("GATEWAY_PATH")

func main() {
	log.Println("Starting execution")

	err := utils.ValidateEnvs()

	if err != nil {
		log.Println("An error occurred while validating environment variables")
		utils.ReturnError(err)
	}

	resources, err := api.GetResources(fmt.Sprintf("%s/%s", apiHost, swaggerJSONPath))

	if err != nil {
		log.Println("An error occurred while performing resource discovery in application SwaggerJSON")
		utils.ReturnError(err)
	}

	err = aws.CreateAllResources(resources, apiHost)

	if err != nil {
		log.Println("An error occurred while creating resources in the Gateway API")
		utils.ReturnError(err)
	}

	log.Println("Process completed successfully")
}
