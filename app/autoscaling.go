package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// AutoscalingEvent
type AutoscalingEvent struct {
	LifecycleHookName    string `json:"LifecycleHookName"`
	AccountID            string `json:"AccountId"`
	RequestID            string `json:"RequestId"`
	LifecycleTransition  string `json:"LifecycleTransition"`
	AutoScalingGroupName string `json:"AutoScalingGroupName"`
	Service              string `json:"Service"`
	Time                 string `json:"Time"`
	EC2InstanceID        string `json:"EC2InstanceId"`
	NotificationMetadata string `json:"NotificationMetadata,omitempty"`
	LifecycleActionToken string `json:"LifecycleActionToken"`

	Event string `json:"Event,omitempty"`
}

// CompleteLifecycleAction
func CompleteLifecycleAction(event AutoscalingEvent, result string) error {
	input := &autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  aws.String(event.AutoScalingGroupName),
		LifecycleActionResult: aws.String(result),
		LifecycleActionToken:  aws.String(event.LifecycleActionToken),
		LifecycleHookName:     aws.String(event.LifecycleHookName),
	}
	_, err := Autoscaling.CompleteLifecycleAction(input)
	return err
}
