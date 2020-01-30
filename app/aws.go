package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// AWS Clients that can be mocked for testing
var (
	Autoscaling = NewAutoscaling()
	Cloudwatch  = NewCloudwatch()
	AwsRegion   = GetAWSRegion()

	sess = session.Must(session.NewSession())
)

// NewAutoscaling is a Autoscaling client
func NewAutoscaling() *autoscaling.AutoScaling {
	c := autoscaling.New(sess)
	return c
}

// NewCloudwatch is a Cloudwatch client
func NewCloudwatch() *cloudwatch.CloudWatch {
	c := cloudwatch.New(sess)
	return c
}

// Get current AWS Region
func GetAWSRegion() string {
	return fmt.Sprintf("%v", sess.Config.Region)
}
