package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	mux := InitAlexaMux()
	lambda.Start(mux.Handle)
}
