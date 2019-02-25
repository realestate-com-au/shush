package ssm

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/realestate-com-au/shush/awsmeta"
)

// Handle structure the client for the Secret Manager
type Handle struct {
	Client    *ssm.SSM
	Prefix    string
	CipherKey string
}

// New returns the reference to Handle object
func New(region string, prefix string, cipherkey string) (ops *Handle, err error) {

	if region == "" {
		region = awsmeta.GetRegion()
		if region == "" {
			err = errors.New("please specify region (--region or $AWS_DEFAULT_REGION)")
			return
		}
	}
	// Create a SSM client with additional configuration
	client := ssm.New(session.New(), aws.NewConfig().WithRegion(region))
	ops = &Handle{
		Client:    client,
		Prefix:    prefix,
		CipherKey: cipherkey,
	}
	return
}

// Decrypt reveal the value of the SSM key
func (h *Handle) Decrypt() (string, error) {
	var (
		cipherkey      = h.CipherKey
		withDecryption = true
	)

	output, err := h.Client.GetParameter(&ssm.GetParameterInput{
		Name:           &cipherkey,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		log.Println("Error Decrypt SSM key")
		return "", err
	}

	return *output.Parameter.Value, nil
}

// DecryptEnv
func (h *Handle) DecryptEnv() {

}
