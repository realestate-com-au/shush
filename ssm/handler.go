package ssm

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/realestate-com-au/shush/awsmeta"
	"github.com/realestate-com-au/shush/sys"
)

// Handle structure the client for the Secret Manager
type Handle struct {
	Client       *ssm.SSM
	Prefix       string
	CipherKey    string
	PlaintextKey string
}

// New returns the reference to Handle object
func New(region string, prefix string, cipherkey string, plaintextKey string) (ops *Handle, err error) {

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
		Client:       client,
		Prefix:       prefix,
		CipherKey:    cipherkey,
		PlaintextKey: plaintextKey,
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
		log.Printf("Error Decrypt SSM key: %v", err)
		return "", err
	}

	return *output.Parameter.Value, nil
}

// DecryptEnv update the local environment variable with decrypted keys
func (h *Handle) DecryptEnv() {

	plaintext, err := h.Decrypt()
	sys.CheckError(err, sys.SsmError)
	os.Setenv(h.PlaintextKey, plaintext)
}
