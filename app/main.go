package main

import (
	"fmt"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) (string, error) {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	return fmt.Sprintf("OK"), nil
}

func main() {
	lambda.Start(handleRequest)
}
