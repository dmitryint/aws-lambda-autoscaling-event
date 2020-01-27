package main

import (
	"fmt"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	fmt.Printf("Event: ", event)
	return fmt.Sprintf("OK"), nil
}

func main() {
	lambda.Start(handleRequest)
}
