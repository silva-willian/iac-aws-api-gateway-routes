package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/silva-willian/iac-aws-api-gateway-routes/api"
)

func createMethod(resourceID, method, authorizer string, parameters []api.ResourceParameters) error {

	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	parametersRequest, err := formatParametersForMethodRequest(parameters)

	if err != nil {
		return err
	}
	_, err = client.PutMethod(&apigateway.PutMethodInput{
		ApiKeyRequired:    aws.Bool(false),
		AuthorizationType: aws.String("NONE"),
		HttpMethod:        &method,
		ResourceId:        &resourceID,
		RestApiId:         &apiID,
		RequestParameters: parametersRequest,
	})

	if err != nil {
		return handleAWSError(err)
	}

	return nil
}

func createMethodIntegration(resourceID, method, host string, parameters []api.ResourceParameters) error {

	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	parametersRequest, err := formatParametersForIntegrationRequest(parameters)

	if err != nil {
		return err
	}

	if method == "OPTIONS" {
		var defaultTemplate = map[string]*string{
			"application/json": aws.String("{'statusCode': 200}"),
		}

		_, err = client.PutIntegration(&apigateway.PutIntegrationInput{
			HttpMethod:            &method,
			IntegrationHttpMethod: &method,
			ResourceId:            &resourceID,
			RestApiId:             &apiID,
			Type:                  aws.String("MOCK"),
			RequestParameters:     parametersRequest,
			RequestTemplates:      defaultTemplate,
		})

	} else {
		_, err = client.PutIntegration(&apigateway.PutIntegrationInput{
			HttpMethod:            &method,
			IntegrationHttpMethod: &method,
			ResourceId:            &resourceID,
			RestApiId:             &apiID,
			Type:                  aws.String("HTTP"),
			Uri:                   &host,
			RequestParameters:     parametersRequest,
		})

	}

	if err != nil {
		return handleAWSError(err)
	}

	return nil
}

func createMethodResponse(resourceID, method, statusCode, pattern string) error {

	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	if method == "OPTIONS" {

		var defaultCors = map[string]*bool{
			"method.response.header.Access-Control-Allow-Headers": aws.Bool(true),
			"method.response.header.Access-Control-Allow-Methods": aws.Bool(true),
			"method.response.header.Access-Control-Allow-Origin":  aws.Bool(true),
		}

		req, _ := client.PutMethodResponseRequest(&apigateway.PutMethodResponseInput{
			HttpMethod:         &method,
			ResourceId:         &resourceID,
			RestApiId:          &apiID,
			StatusCode:         &statusCode,
			ResponseParameters: defaultCors,
		})

		err = req.Send()

	} else {
		var defaultCors = map[string]*bool{

			"method.response.header.Access-Control-Allow-Origin": aws.Bool(true),
		}

		req, _ := client.PutMethodResponseRequest(&apigateway.PutMethodResponseInput{
			HttpMethod:         &method,
			ResourceId:         &resourceID,
			RestApiId:          &apiID,
			StatusCode:         &statusCode,
			ResponseParameters: defaultCors,
		})

		err = req.Send()

	}

	if err != nil {
		return handleAWSError(err)
	}

	if method == "OPTIONS" {

		var defaultCors = map[string]*string{

			"method.response.header.Access-Control-Allow-Headers": aws.String("'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,userInformation,idToken'"),
			"method.response.header.Access-Control-Allow-Methods": aws.String("'GET,OPTIONS,POST,PUT,PATCH,DELETE'"),
			"method.response.header.Access-Control-Allow-Origin":  aws.String("'*'"),
		}

		_, err = client.PutIntegrationResponse(&apigateway.PutIntegrationResponseInput{
			HttpMethod:         &method,
			ResourceId:         &resourceID,
			RestApiId:          &apiID,
			StatusCode:         &statusCode,
			SelectionPattern:   &pattern,
			ResponseParameters: defaultCors,
		})

	} else {
		var defaultCors = map[string]*string{
			"method.response.header.Access-Control-Allow-Origin": aws.String("'*'"),
		}

		_, err = client.PutIntegrationResponse(&apigateway.PutIntegrationResponseInput{
			HttpMethod:         &method,
			ResourceId:         &resourceID,
			RestApiId:          &apiID,
			StatusCode:         &statusCode,
			SelectionPattern:   &pattern,
			ResponseParameters: defaultCors,
		})

	}

	if err != nil {
		return handleAWSError(err)
	}

	return nil
}

func formatParametersForMethodRequest(parameters []api.ResourceParameters) (map[string]*bool, error) {

	if parameters == nil || len(parameters) <= 0 {
		return nil, nil
	}

	list := map[string]*bool{}
	list["method.request.header.Accept"] = aws.Bool(false)
	list["method.request.header.Content-type"] = aws.Bool(false)

	for _, parameter := range parameters {
		list[fmt.Sprintf("method.request.%s.%s", parameter.Location, parameter.Name)] = &parameter.Required
	}

	return list, nil
}

func formatParametersForIntegrationRequest(parameters []api.ResourceParameters) (map[string]*string, error) {

	if parameters == nil || len(parameters) <= 0 {
		return nil, nil
	}

	list := map[string]*string{}
	key := fmt.Sprintf("integration.request.header.Accept")
	value := fmt.Sprintf("method.request.header.Accept")
	list[key] = &value
	key2 := fmt.Sprintf("integration.request.header.Content-type")
	value2 := fmt.Sprintf("method.request.header.Content-type")
	list[key2] = &value2

	for _, parameter := range parameters {
		key := fmt.Sprintf("integration.request.%s.%s", parameter.Location, parameter.Name)
		value := fmt.Sprintf("method.request.%s.%s", parameter.Location, parameter.Name)
		list[key] = &value
	}

	return list, nil
}
