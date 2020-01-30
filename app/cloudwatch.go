package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type AutoscalingNotificationMetadata struct {
	DiskSpaceUtilizationPeriod     int64   `json:"DiskSpaceUtilizationPeriod,string"`
	DiskSpaceUtilizationThreshold  float64 `json:"DiskSpaceUtilizationThreshold,string"`
	SNSNotificationTopicArn        string  `json:"SNSNotificationTopicArn"`
	DiskSpaceUtilizationFilesystem string  `json:"DiskSpaceUtilizationFilesystem"`
	DiskSpaceUtilizationMountPath  string  `json:"DiskSpaceUtilizationMountPath"`
}

// Creates CloudWatch Alarm for specified instance
func CWPutDiskSpaceUtilizationMetricAlarm(event AutoscalingEvent) error {
	metadata := AutoscalingNotificationMetadata{
		DiskSpaceUtilizationPeriod:     300,
		DiskSpaceUtilizationThreshold:  70.0,
		DiskSpaceUtilizationFilesystem: "/dev/nvme0n1p1",
		DiskSpaceUtilizationMountPath:  "/",
	}
	alarmActions := []*string{}
	actionsEnabled := false
	if err := json.Unmarshal([]byte(event.NotificationMetadata), &metadata); err != nil {
		return err
	}
	if metadata.SNSNotificationTopicArn != "" {
		alarmActions = append(alarmActions, aws.String(metadata.SNSNotificationTopicArn))
		actionsEnabled = true
	}
	_, err := Cloudwatch.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String("ASG/" + event.AutoScalingGroupName + "/" + event.EC2InstanceID + "/DiskSpaceUtilization"),
		ComparisonOperator: aws.String(cloudwatch.ComparisonOperatorGreaterThanThreshold),
		EvaluationPeriods:  aws.Int64(1),
		MetricName:         aws.String("DiskSpaceUtilization"),
		Namespace:          aws.String("System/Linux"),
		Period:             aws.Int64(metadata.DiskSpaceUtilizationPeriod),
		Statistic:          aws.String(cloudwatch.StatisticAverage),
		Threshold:          aws.Float64(metadata.DiskSpaceUtilizationThreshold),
		ActionsEnabled:     aws.Bool(actionsEnabled),
		AlarmDescription:   aws.String(fmt.Sprintf("Alarm when server Disk Space Utilization exceeds %.2f%", metadata.DiskSpaceUtilizationThreshold)),
		AlarmActions:       alarmActions,
		// This is apart of the default workflow actions. This one will reboot the instance, if the
		// alarm is triggered.
		// AlarmActions: []*string{
		// 	aws.String(fmt.Sprintf("arn:aws:swf:us-east-1:%s:action/actions/AWS_EC2.InstanceId.Reboot/1.0", instance)),
		// },
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("MountPath"),
				Value: aws.String(metadata.DiskSpaceUtilizationMountPath),
			},
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(event.EC2InstanceID),
			},
			{
				Name:  aws.String("Filesystem"),
				Value: aws.String(metadata.DiskSpaceUtilizationFilesystem),
			},
		},
	})
	return err
}

// Removes CloudWatch Alarm for specified instance
func CWDeleteDiskSpaceUtilizationMetricAlarm(event AutoscalingEvent) error {
	params := &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []*string{
			aws.String("ASG/" + event.AutoScalingGroupName + "/" + event.EC2InstanceID + "/DiskSpaceUtilization"),
		}}
	resp, err := Cloudwatch.DeleteAlarms(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	return nil
}

// Creates CloudWatch Alarm for specified instance
func CWPutStatusCheckFailedMetricAlarm(event AutoscalingEvent) error {
	metadata := AutoscalingNotificationMetadata{}
	alarmActions := []*string{
		// https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_PutMetricAlarm.html
		// Recover EC2 Instance
		aws.String(fmt.Sprintf("arn:aws:automate:%s:ec2:recover", AwsRegion)),
	}
	if err := json.Unmarshal([]byte(event.NotificationMetadata), &metadata); err != nil {
		return err
	}
	if metadata.SNSNotificationTopicArn != "" {
		alarmActions = append(alarmActions, aws.String(metadata.SNSNotificationTopicArn))
	}
	_, err := Cloudwatch.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String("ASG/" + event.AutoScalingGroupName + "/" + event.EC2InstanceID + "/StatusCheckFailed"),
		ComparisonOperator: aws.String(cloudwatch.ComparisonOperatorGreaterThanOrEqualToThreshold),
		EvaluationPeriods:  aws.Int64(1),
		Namespace:          aws.String("AWS/EC2"),
		MetricName:         aws.String("StatusCheckFailed"),
		Period:             aws.Int64(300.0),
		Statistic:          aws.String(cloudwatch.StatisticMaximum),
		Threshold:          aws.Float64(1.0),
		ActionsEnabled:     aws.Bool(true),
		AlarmDescription:   aws.String("Alarm when the Status Check Failed"),
		AlarmActions:       alarmActions,
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(event.EC2InstanceID),
			},
		},
	})
	return err
}

// Removes CloudWatch Alarm for specified instance
func CWDeleteCheckFailedMetricAlarm(event AutoscalingEvent) error {
	params := &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []*string{
			aws.String("ASG/" + event.AutoScalingGroupName + "/" + event.EC2InstanceID + "/StatusCheckFailed"),
		}}
	resp, err := Cloudwatch.DeleteAlarms(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	return nil
}
