package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/apigateway"
)

func deleteBaseResourceIFExists() error {

	if !resourceExists(fmt.Sprintf("/%s", basePath)) {
		return nil
	}

	resourceID, err := getResourceID(fmt.Sprintf("/%s", basePath))

	if err != nil {
		return err
	}

	err = deleteResource(resourceID)

	if err != nil {
		return err
	}

	itens, err := getAllResources()

	if err != nil {
		return err
	}

	gatewayResources = itens
	return nil
}

func getAllResources() ([]GatewayResource, error) {
	client, err := getGatewayClient()

	if err != nil {
		return nil, err
	}

	pageNum := 0
	itens := []GatewayResource{}

	err = client.GetResourcesPages(&apigateway.GetResourcesInput{
		RestApiId: &apiID,
	}, func(page *apigateway.GetResourcesOutput, lastPage bool) bool {
		pageNum++

		for _, resource := range page.Items {
			item := GatewayResource{
				Path: *resource.Path,
				ID:   *resource.Id,
			}

			itens = append(itens, item)
		}

		return !lastPage
	})
	log.Println(pageNum)

	if err != nil {
		return nil, handleAWSError(err)
	}

	return itens, nil
}

func createResource(path, pathParent string) error {
	formatPath := formatCompleteResource(pathParent, path)

	if resourceExists(formatPath) {
		return nil
	}

	parentID, err := getResourceID(pathParent)

	if err != nil {
		return err
	}

	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	resource, err := client.CreateResource(&apigateway.CreateResourceInput{
		RestApiId: &apiID,
		PathPart:  &path,
		ParentId:  &parentID,
	})

	if err != nil {
		return handleAWSError(err)
	}

	gatewayResources = append(gatewayResources, GatewayResource{
		Path: formatPath,
		ID:   *resource.Id,
	})

	log.Printf("Path [%s] successfully created on Gateway", formatPath)

	return nil
}

func deleteResource(resourceID string) error {

	client, err := getGatewayClient()

	if err != nil {
		return err
	}

	_, err = client.DeleteResource(&apigateway.DeleteResourceInput{
		RestApiId:  &apiID,
		ResourceId: &resourceID,
	})

	if err != nil {
		return err
	}

	return nil
}

func formatCompleteResource(pathParent, path string) string {

	if pathParent == "/" {
		return fmt.Sprintf("/%s", path)
	}

	return fmt.Sprintf("%s/%s", pathParent, path)
}

func resourceExists(path string) bool {

	for _, resource := range gatewayResources {

		if resource.Path == path {
			return true
		}
	}

	return false
}

func getResourceID(path string) (string, error) {

	path = standardizingPath(path)

	for _, resource := range gatewayResources {

		if resource.Path == path {
			return resource.ID, nil
		}
	}

	return "", fmt.Errorf("ResourceID of route [%s] does not exist", path)
}

func standardizingPath(path string) string {

	if path == "/" {
		return path
	}

	lastChar := path[len(path)-1 : len(path)]

	if lastChar == "/" {
		return path[0 : len(path)-1]
	}

	return path
}
