package main

import (
	"fmt"

	"gitlab.com/countableset/lambda-s3-cloudflare/awss3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// TODO
	// configure creds profile?
	// configure bucket name?
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "logtest")},
	)
	svc := s3.New(sess)
	bucket := "logs.countableset.com"
	policy := awss3.GetPolicy(svc, bucket)
	fmt.Println(policy)
	fmt.Println("Hello, world.")
}
