package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func onEc2InstanceLaunching(event AutoscalingEvent) error {
	err := CWPutDiskSpaceUtilizationMetricAlarm(event)
	if err != nil {
		CompleteLifecycleAction(event, "ABANDON")
		return err
	}

	err = CWPutStatusCheckFailedSystemMetricAlarm(event)
	if err != nil {
		CompleteLifecycleAction(event, "ABANDON")
		return err
	}

	err = CWPutStatusCheckFailedInstanceMetricAlarm(event)
	if err != nil {
		CompleteLifecycleAction(event, "ABANDON")
		return err
	}

	CompleteLifecycleAction(event, "CONTINUE")
	return err
}

func onEc2InstanceTerminating(event AutoscalingEvent) error {
	errorResponse := error(nil)
	err := CWDeleteDiskSpaceUtilizationMetricAlarm(event)
	if err != nil {
		errorResponse = err
	}

	err = CWDeleteCheckFailedSystemMetricAlarm(event)
	if err != nil {
		errorResponse = err
	}

	err = CWDeleteCheckFailedInstanceMetricAlarm(event)
	if err != nil {
		errorResponse = err
	}

	CompleteLifecycleAction(event, "CONTINUE")
	return errorResponse
}

func makeEventHandler(record AutoscalingEvent) (func(event AutoscalingEvent) error, error) {
	eventName := record.Event
	if record.LifecycleTransition != "" {
		eventName = record.LifecycleTransition
	}
	switch eventName {
	case "autoscaling:EC2_INSTANCE_LAUNCHING":
		return onEc2InstanceLaunching, nil
	case "autoscaling:EC2_INSTANCE_TERMINATING":
		return onEc2InstanceTerminating, nil
	case "autoscaling:TEST_NOTIFICATION":
		return func(event AutoscalingEvent) error {
			return nil
		}, nil
	default:
		return func(event AutoscalingEvent) error {
			return fmt.Errorf("")
		}, fmt.Errorf("Unknown LifecycleTransition: %s", record.LifecycleTransition)
	}
}

func handleRequest(ctx context.Context, sqsEvent events.SNSEvent) (string, error) {
	for _, message := range sqsEvent.Records {
		entity := message.SNS
		fmt.Printf("The message %s for event source %s = %s \n", entity.MessageID, message.EventSource, entity.Message)
		var record AutoscalingEvent
		if err := json.Unmarshal([]byte(entity.Message), &record); err != nil {
			return "", fmt.Errorf("Unable to unmarshal Message")
		}
		processEvent, err := makeEventHandler(record)
		if err != nil {
			return "", err
		}
		err = processEvent(record)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("OK"), nil
}

func main() {
	lambda.Start(handleRequest)
}
