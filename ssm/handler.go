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
type Handler struct {
	Client       *ssm.SSM
	Prefix       string
	CipherKey    string
	PlaintextKey string
	KeyID        string
	Plaintext    string
}

// Client establish a session to AWS
func Client(region string) (client *ssm.SSM, err error) {
	if region == "" {
		region = awsmeta.GetRegion()
		if region == "" {
			err = errors.New("please specify region (--region or $AWS_DEFAULT_REGION)")
			return
		}
	}
	client = ssm.New(session.New(), aws.NewConfig().WithRegion(region))
	return
}

func (h *Handler) Encrypt() (string, error) {
	return "", nil
}

// Decrypt reveal the value of the SSM key
func (h *Handler) Decrypt() (string, error) {
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
func (h *Handler) DecryptEnv() {

	plaintext, err := h.Decrypt()
	sys.CheckError(err, sys.SsmError)
	os.Setenv(h.PlaintextKey, plaintext)
}
