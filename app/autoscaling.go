package main

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
	NotificationMetadata string `json:"NotificationMetadata"`
	LifecycleActionToken string `json:"LifecycleActionToken"`

	Event string `json:"Event,omitempty"`
}
