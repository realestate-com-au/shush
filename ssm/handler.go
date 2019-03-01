package ssm

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/realestate-com-au/shush/awsmeta"
	"github.com/realestate-com-au/shush/sys"
)

// Handler structure the client for the Secret Manager
type Handler struct {
	// Service provide interfaces for AWS SDK features
	Service AWSIface

	// PlaintextKey used in environment variables as the key without the prefix
	PlaintextKey string

	// KMSKeyID is the ID for KMS key, used in parameter store encryption
	KMSKeyID string

	// Plaintext is the decrypted secret
	Plaintext string

	// ParameterKeyName is the name of the parameter. For information about valid values for parameter names, see https://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-parameter-name-constraints.html
	ParameterKeyName string

	// ParameterValue must not nest another parameter. Do not use {{}} in the value.
	ParameterValue string

	// The type of parameter. Valid values include the following: String or SecureString.
	// NOTE: AWS CloudFormation doesn't support the SecureString parameter type.
	// For more information see https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-parameter.html
	ParameterType string
}

// Client prepare AWS config
func Client(region string) *ssm.SSM {
	if region == "" {
		region = awsmeta.GetRegion()
		if region == "" {
			log.Fatalln("please specify region (--region or $AWS_DEFAULT_REGION)")
		}
	}
	return ssm.New(session.New(), aws.NewConfig().WithRegion(region))
}

// AWSIface abstract AWS SDK required method
type AWSIface interface {
	PutParameter(*ssm.PutParameterInput) (*ssm.PutParameterOutput, error)
	GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
}

// Encrypt the SSM parameter value
func (h *Handler) Encrypt() (string, error) {

	// Only support the `String` and `SecureString`
	var (
		overwrite = true
		err       error
	)

	// Initial PutParameterInput
	switch {
	case len(h.KMSKeyID) != 0:
		// Encrypted Parameter keys
		h.setParameterType("SecureString")
		_, err = h.Service.PutParameter(&ssm.PutParameterInput{
			Name:      &(h.ParameterKeyName),
			KeyId:     &(h.KMSKeyID),
			Type:      &(h.ParameterType),
			Value:     &(h.ParameterValue),
			Overwrite: &overwrite,
		})
	default:
		// Unencrypted Parameter keys
		h.setParameterType("String")
		_, err = h.Service.PutParameter(&ssm.PutParameterInput{
			Name:      &(h.ParameterKeyName),
			Type:      &(h.ParameterType),
			Value:     &(h.ParameterValue),
			Overwrite: &overwrite,
		})
	}

	if err != nil {
		log.Printf("Error Encrypt the SSM value: %v", err)
		return "", err
	}
	return h.ParameterKeyName, nil
}

// setParameterType for encrypted parameter keys
func (h *Handler) setParameterType(paramType string) {
	h.ParameterType = paramType
}

// Decrypt reveal the value of the SSM parameter key
func (h *Handler) Decrypt() (string, error) {

	withDecryption := true

	output, err := h.Service.GetParameter(&ssm.GetParameterInput{
		Name:           &(h.ParameterKeyName),
		WithDecryption: &withDecryption,
	})
	if err != nil {
		log.Printf("Error Decrypt SSM key: %v", err)
		return "", err
	}

	return *output.Parameter.Value, nil
}

// DecryptEnv update the local environment variable with decrypted keys
func (h *Handler) DecryptEnv() {

	plaintext, err := h.Decrypt()
	sys.CheckError(err, sys.SsmError)
	os.Setenv(h.PlaintextKey, plaintext)
}
