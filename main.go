package main

import (
	// fake import for handlers to run their init() functions
	_ "github.com/lepinkainen/titleparser/handler"
	"github.com/lepinkainen/titleparser/lambda"

	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	awslambda.Start(lambda.HandleRequest)
}
