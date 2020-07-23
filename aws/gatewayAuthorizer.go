package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/apigateway"
)

var authorizerName string = os.Getenv("GATEWAY_AUTHORIZER_NAME")
var authorizerEnable string = os.Getenv("GATEWAY_AUTHORIZER_ENABLE")

func getAuthorizer() (string, error) {

	if authorizerEnable != "true" {
		return "", nil
	}

	client, err := getGatewayClient()

	if err != nil {
		return "", err
	}

	result, err := client.GetAuthorizers(&apigateway.GetAuthorizersInput{
		RestApiId: &apiID,
	})

	if err != nil {
		return "", handleAWSError(err)
	}

	if len(result.Items) <= 0 {
		return "", fmt.Errorf("API Gateway has no authorizers %v", result)
	}

	for _, item := range result.Items {
		if *item.Name == authorizerName {
			return *item.Id, nil
		}
	}

	return "", fmt.Errorf("Lambda Authorizer Not Found %v", result)
}
