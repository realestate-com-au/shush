package kms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/realestate-com-au/shush/awsmeta"
)

type kmsEncryptionContext map[string]string

// Structure encapsulating stuff common to encrypt and decrypt.
//
type KmsHandle struct {
	Client  *kms.Client
	Context kmsEncryptionContext
}

func NewHandle(region string, encryptionContextData []string) (ops *KmsHandle, err error) {
	encryptionContext, err := parseEncryptionContext(encryptionContextData)
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

	conf, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		return nil, fmt.Errorf("could not create AWS session: %v", err)
	}

	client := kms.NewFromConfig(conf)
	ops = &KmsHandle{
		Client:  client,
		Context: encryptionContext,
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
		context[parts[0]] = parts[1]
	}
	return context, nil
}

// Encrypt plaintext using specified key.
func (h *KmsHandle) Encrypt(plaintext string, keyID string) (string, error) {
	output, err := h.Client.Encrypt(context.TODO(), &kms.EncryptInput{
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
func (h *KmsHandle) Decrypt(ciphertext string) (string, string, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", "", err
	}
	output, err := h.Client.Decrypt(context.TODO(), &kms.DecryptInput{
		EncryptionContext: h.Context,
		CiphertextBlob:    ciphertextBlob,
	})
	if err != nil {
		return "", "", err
	}

	keyId := ""
	if output.KeyId != nil {
		keyId = *output.KeyId
	}
	return string(output.Plaintext), keyId, nil
}
