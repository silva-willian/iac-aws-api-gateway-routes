package aws

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
)

var awsRegion string = os.Getenv("AWS_REGION")

func getSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
}

func handleAWSError(err error) error {
	if err == nil {
		return nil
	}
	awsErr, ok := err.(awserr.Error)

	if ok {
		log.Println(awsErr)
	}
	return err
}
