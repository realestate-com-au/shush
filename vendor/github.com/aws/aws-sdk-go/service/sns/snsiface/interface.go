// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package snsiface provides an interface for the Amazon Simple Notification Service.
package snsiface

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNSAPI is the interface type for sns.SNS.
type SNSAPI interface {
	AddPermissionRequest(*sns.AddPermissionInput) (*request.Request, *sns.AddPermissionOutput)

	AddPermission(*sns.AddPermissionInput) (*sns.AddPermissionOutput, error)

	ConfirmSubscriptionRequest(*sns.ConfirmSubscriptionInput) (*request.Request, *sns.ConfirmSubscriptionOutput)

	ConfirmSubscription(*sns.ConfirmSubscriptionInput) (*sns.ConfirmSubscriptionOutput, error)

	CreatePlatformApplicationRequest(*sns.CreatePlatformApplicationInput) (*request.Request, *sns.CreatePlatformApplicationOutput)

	CreatePlatformApplication(*sns.CreatePlatformApplicationInput) (*sns.CreatePlatformApplicationOutput, error)

	CreatePlatformEndpointRequest(*sns.CreatePlatformEndpointInput) (*request.Request, *sns.CreatePlatformEndpointOutput)

	CreatePlatformEndpoint(*sns.CreatePlatformEndpointInput) (*sns.CreatePlatformEndpointOutput, error)

	CreateTopicRequest(*sns.CreateTopicInput) (*request.Request, *sns.CreateTopicOutput)

	CreateTopic(*sns.CreateTopicInput) (*sns.CreateTopicOutput, error)

	DeleteEndpointRequest(*sns.DeleteEndpointInput) (*request.Request, *sns.DeleteEndpointOutput)

	DeleteEndpoint(*sns.DeleteEndpointInput) (*sns.DeleteEndpointOutput, error)

	DeletePlatformApplicationRequest(*sns.DeletePlatformApplicationInput) (*request.Request, *sns.DeletePlatformApplicationOutput)

	DeletePlatformApplication(*sns.DeletePlatformApplicationInput) (*sns.DeletePlatformApplicationOutput, error)

	DeleteTopicRequest(*sns.DeleteTopicInput) (*request.Request, *sns.DeleteTopicOutput)

	DeleteTopic(*sns.DeleteTopicInput) (*sns.DeleteTopicOutput, error)

	GetEndpointAttributesRequest(*sns.GetEndpointAttributesInput) (*request.Request, *sns.GetEndpointAttributesOutput)

	GetEndpointAttributes(*sns.GetEndpointAttributesInput) (*sns.GetEndpointAttributesOutput, error)

	GetPlatformApplicationAttributesRequest(*sns.GetPlatformApplicationAttributesInput) (*request.Request, *sns.GetPlatformApplicationAttributesOutput)

	GetPlatformApplicationAttributes(*sns.GetPlatformApplicationAttributesInput) (*sns.GetPlatformApplicationAttributesOutput, error)

	GetSubscriptionAttributesRequest(*sns.GetSubscriptionAttributesInput) (*request.Request, *sns.GetSubscriptionAttributesOutput)

	GetSubscriptionAttributes(*sns.GetSubscriptionAttributesInput) (*sns.GetSubscriptionAttributesOutput, error)

	GetTopicAttributesRequest(*sns.GetTopicAttributesInput) (*request.Request, *sns.GetTopicAttributesOutput)

	GetTopicAttributes(*sns.GetTopicAttributesInput) (*sns.GetTopicAttributesOutput, error)

	ListEndpointsByPlatformApplicationRequest(*sns.ListEndpointsByPlatformApplicationInput) (*request.Request, *sns.ListEndpointsByPlatformApplicationOutput)

	ListEndpointsByPlatformApplication(*sns.ListEndpointsByPlatformApplicationInput) (*sns.ListEndpointsByPlatformApplicationOutput, error)

	ListEndpointsByPlatformApplicationPages(*sns.ListEndpointsByPlatformApplicationInput, func(*sns.ListEndpointsByPlatformApplicationOutput, bool) bool) error

	ListPlatformApplicationsRequest(*sns.ListPlatformApplicationsInput) (*request.Request, *sns.ListPlatformApplicationsOutput)

	ListPlatformApplications(*sns.ListPlatformApplicationsInput) (*sns.ListPlatformApplicationsOutput, error)

	ListPlatformApplicationsPages(*sns.ListPlatformApplicationsInput, func(*sns.ListPlatformApplicationsOutput, bool) bool) error

	ListSubscriptionsRequest(*sns.ListSubscriptionsInput) (*request.Request, *sns.ListSubscriptionsOutput)

	ListSubscriptions(*sns.ListSubscriptionsInput) (*sns.ListSubscriptionsOutput, error)

	ListSubscriptionsPages(*sns.ListSubscriptionsInput, func(*sns.ListSubscriptionsOutput, bool) bool) error

	ListSubscriptionsByTopicRequest(*sns.ListSubscriptionsByTopicInput) (*request.Request, *sns.ListSubscriptionsByTopicOutput)

	ListSubscriptionsByTopic(*sns.ListSubscriptionsByTopicInput) (*sns.ListSubscriptionsByTopicOutput, error)

	ListSubscriptionsByTopicPages(*sns.ListSubscriptionsByTopicInput, func(*sns.ListSubscriptionsByTopicOutput, bool) bool) error

	ListTopicsRequest(*sns.ListTopicsInput) (*request.Request, *sns.ListTopicsOutput)

	ListTopics(*sns.ListTopicsInput) (*sns.ListTopicsOutput, error)

	ListTopicsPages(*sns.ListTopicsInput, func(*sns.ListTopicsOutput, bool) bool) error

	PublishRequest(*sns.PublishInput) (*request.Request, *sns.PublishOutput)

	Publish(*sns.PublishInput) (*sns.PublishOutput, error)

	RemovePermissionRequest(*sns.RemovePermissionInput) (*request.Request, *sns.RemovePermissionOutput)

	RemovePermission(*sns.RemovePermissionInput) (*sns.RemovePermissionOutput, error)

	SetEndpointAttributesRequest(*sns.SetEndpointAttributesInput) (*request.Request, *sns.SetEndpointAttributesOutput)

	SetEndpointAttributes(*sns.SetEndpointAttributesInput) (*sns.SetEndpointAttributesOutput, error)

	SetPlatformApplicationAttributesRequest(*sns.SetPlatformApplicationAttributesInput) (*request.Request, *sns.SetPlatformApplicationAttributesOutput)

	SetPlatformApplicationAttributes(*sns.SetPlatformApplicationAttributesInput) (*sns.SetPlatformApplicationAttributesOutput, error)

	SetSubscriptionAttributesRequest(*sns.SetSubscriptionAttributesInput) (*request.Request, *sns.SetSubscriptionAttributesOutput)

	SetSubscriptionAttributes(*sns.SetSubscriptionAttributesInput) (*sns.SetSubscriptionAttributesOutput, error)

	SetTopicAttributesRequest(*sns.SetTopicAttributesInput) (*request.Request, *sns.SetTopicAttributesOutput)

	SetTopicAttributes(*sns.SetTopicAttributesInput) (*sns.SetTopicAttributesOutput, error)

	SubscribeRequest(*sns.SubscribeInput) (*request.Request, *sns.SubscribeOutput)

	Subscribe(*sns.SubscribeInput) (*sns.SubscribeOutput, error)

	UnsubscribeRequest(*sns.UnsubscribeInput) (*request.Request, *sns.UnsubscribeOutput)

	Unsubscribe(*sns.UnsubscribeInput) (*sns.UnsubscribeOutput, error)
}

var _ SNSAPI = (*sns.SNS)(nil)
