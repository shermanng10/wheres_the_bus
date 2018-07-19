package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	mux := AlexaMuxFactory()
	lambda.Start(mux.Handle)
}
