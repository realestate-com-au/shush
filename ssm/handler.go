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
	Service          AWSIface
	Prefix           string
	PlaintextKey     string
	KMSKeyID         string
	Plaintext        string
	ParameterKeyName string
	ParameterValue   string
	ParameterType    string
}

// Client establish a session to AWS
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

// AWSImpl indicate SSM client
type AWSImpl struct {
	*ssm.SSM
}

// PutParameter implement AWS SDK SSM service
func (impl *AWSImpl) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	return impl.SSM.PutParameter(input)
}

// GetParameter implement AWS SDK SSM service
func (impl *AWSImpl) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return impl.SSM.GetParameter(input)
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
