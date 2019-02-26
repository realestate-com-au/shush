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

type kmsEncryptionContext map[string]*string

// Structure encapsulating stuff common to encrypt and decrypt.
//
type Handle struct {
	Client       *kms.KMS
	Context      kmsEncryptionContext
	Prefix       string
	CipherKey    string
	PlaintextKey string
}

func New(region string, context []string, prefix string, cipherkey string, plaintextKey string) (ops *Handle, err error) {
	encryptionContext, err := parseEncryptionContext(context)
	if err != nil {
		return nil, fmt.Errorf("could not parse encryption context: %v", err)
	}
	if region == "" {
		region = awsmeta.GetRegion()
		if region == "" {
			err = errors.New("please specify region (--region or $AWS_DEFAULT_REGION)")
			return
		}
	}
	client := kms.New(session.New(), aws.NewConfig().WithRegion(region))
	ops = &Handle{
		Client:       client,
		Context:      encryptionContext,
		Prefix:       prefix,
		CipherKey:    cipherkey,
		PlaintextKey: plaintextKey,
	}
	return
}

func parseEncryptionContext(contextStrings []string) (kmsEncryptionContext, error) {
	context := make(kmsEncryptionContext, len(contextStrings))
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
func (h *Handle) Encrypt(plaintext string, keyID string) (string, error) {
	output, err := h.Client.Encrypt(&kms.EncryptInput{
		KeyId:             &keyID,
		EncryptionContext: h.Context,
		Plaintext:         []byte(plaintext),
	})
	if err != nil {
		return "", err
	}
	ciphertext := base64.StdEncoding.EncodeToString(output.CiphertextBlob)
	return ciphertext, nil
}

// Decrypt ciphertext.
func (h *Handle) Decrypt() (string, error) {
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
func (h *Handle) DecryptEnv() {

	plaintext, err := h.Decrypt()
	sys.CheckError(err, sys.KmsError)
	os.Setenv(h.PlaintextKey, plaintext)
}
