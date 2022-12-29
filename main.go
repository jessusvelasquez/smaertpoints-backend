package main

import (
	"smartpoints/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main()  {
	lambda.Start(handler.StartHandler)
}