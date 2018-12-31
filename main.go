package main

import (
	"github.com/lepinkainen/titleparser/lambda"

	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	awslambda.Start(lambda.HandleRequest)
}
