package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/apigateway"
)

var stage string = os.Getenv("GATEWAY_STAGE")
var description string = fmt.Sprintf("Project Publication [%s], base path [%s]", os.Getenv("PROJECT"), os.Getenv("API_GATEWAY_BASE_PATH"))

func createDeploy() error {

	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	_, err = client.CreateDeployment(&apigateway.CreateDeploymentInput{
		RestApiId:        &apiID,
		StageName:        &stage,
		StageDescription: &stage,
		Description:      &description,
	})

	if err != nil {
		return handleAWSError(err)
	}

	return nil
}
