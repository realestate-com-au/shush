package kms

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/realestate-com-au/shush/awsmeta"
	"github.com/realestate-com-au/shush/sys"
)

// EncryptionContext define the format required for kms encryption context
type EncryptionContext map[string]*string

// Handler Structure encapsulating stuff common to encrypt and decrypt.
type Handler struct {
	Client       *kms.KMS
	Context      EncryptionContext
	Prefix       string
	CipherKey    string
	PlaintextKey string
	KeyID        string
	Plaintext    string
}

// Client establish a session to AWS
func Client(region string) (client *kms.KMS, err error) {
	if region == "" {
		region = awsmeta.GetRegion()
		if region == "" {
			err = errors.New("please specify region (--region or $AWS_DEFAULT_REGION)")
			return
		}
	}
	client = kms.New(session.New(), aws.NewConfig().WithRegion(region))
	return
}

// ParseEncryptionContext encryption context is required to decrypt the data
func ParseEncryptionContext(contextStrings []string) (EncryptionContext, error) {
	context := make(EncryptionContext, len(contextStrings))
	for _, s := range contextStrings {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("context must be provided in NAME=VALUE format")
		}
		context[parts[0]] = &parts[1]
	}
	return context, nil
}

// Encrypt plaintext using specified key.
func (h *Handler) Encrypt() (string, error) {

	output, err := h.Client.Encrypt(&kms.EncryptInput{
		KeyId:             &h.KeyID,
		EncryptionContext: h.Context,
		Plaintext:         []byte(h.Plaintext),
	})
	if err != nil {
		return "", err
	}
	ciphertext := base64.StdEncoding.EncodeToString(output.CiphertextBlob)
	return ciphertext, nil
}

// Decrypt ciphertext.
func (h *Handler) Decrypt() (string, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(h.CipherKey)
	if err != nil {
		return "", err
	}
	output, err := h.Client.Decrypt(&kms.DecryptInput{
		EncryptionContext: h.Context,
		CiphertextBlob:    ciphertextBlob,
	})
	if err != nil {
		return "", err
	}
	return string(output.Plaintext), nil
}

// DecryptEnv update the local environment variable with decrypted keys
func (h *Handler) DecryptEnv() {

	plaintext, err := h.Decrypt()
	sys.CheckError(err, sys.KmsError)
	os.Setenv(h.PlaintextKey, plaintext)
}
