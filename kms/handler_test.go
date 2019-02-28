package kms

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/realestate-com-au/shush/kms/mock_kms"
)

func TestClient(t *testing.T) {
	output := Client("ap-southeast-2")
	var k *kms.KMS
	assert.IsType(t, k, output, "Client should be kms")
	assert.Equal(t, "https://kms.ap-southeast-2.amazonaws.com", output.Endpoint, "Should be kms API endpoint")
}

func TestParseEncryptionContext(t *testing.T) {
	ctrl := gomock.NewController(t)

	context := []string{"a=b"}

	defer ctrl.Finish()

	m := mock_kms.NewMockAWSIface(ctrl)
	output, _ := (&Handler{
		Service: m,
		Context: context,
	}).ParseEncryptionContext()

	assert.Equal(t, "b", *output["a"], "Encryption Context should be set")
}

func TestEncrypt(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		keyID     = ""
		plainText = ""
		context   = []string{"a=b"}
	)
	defer ctrl.Finish()
	m := mock_kms.NewMockAWSIface(ctrl)
	h := &Handler{
		Service:   m,
		Context:   context,
		Plaintext: plainText,
	}
	ec, _ := h.ParseEncryptionContext()
	m.EXPECT().Encrypt(gomock.Eq(&kms.EncryptInput{
		KeyId:             &(keyID),
		EncryptionContext: ec,
		Plaintext:         []byte(plainText),
	}))

	expect := base64.StdEncoding.EncodeToString([]byte("This is a kms secret"))
	output, _ := h.Encrypt()

	assert.Equal(t, expect, output, "ciphertext should be equal")
}

func TestDecrypt(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		plainText = ""
		context   = []string{"a=b"}
		cipherKey = ""
	)
	defer ctrl.Finish()
	m := mock_kms.NewMockAWSIface(ctrl)
	h := &Handler{
		Service:   m,
		Context:   context,
		Plaintext: plainText,
		CipherKey: cipherKey,
	}
	ciphertextBlob, _ := base64.StdEncoding.DecodeString(h.CipherKey)
	ec, _ := h.ParseEncryptionContext()
	m.EXPECT().Decrypt(gomock.Eq(&kms.DecryptInput{
		EncryptionContext: ec,
		CiphertextBlob:    ciphertextBlob,
	}))

	output, _ := h.Decrypt()
	assert.Equal(t, "This is a plain text", output, "plaintext should be equal after decryption")
}

func TestDecryptEnv(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		plainText    = ""
		context      = []string{"a=b"}
		cipherKey    = ""
		plainTextKey = "KEY-NAME"
	)
	defer ctrl.Finish()
	m := mock_kms.NewMockAWSIface(ctrl)
	h := &Handler{
		Service:      m,
		Context:      context,
		Plaintext:    plainText,
		CipherKey:    cipherKey,
		PlaintextKey: plainTextKey,
	}
	ciphertextBlob, _ := base64.StdEncoding.DecodeString(h.CipherKey)
	ec, _ := h.ParseEncryptionContext()
	m.EXPECT().Decrypt(gomock.Eq(&kms.DecryptInput{
		EncryptionContext: ec,
		CiphertextBlob:    ciphertextBlob,
	}))

	h.DecryptEnv()
	assert.Equal(t, "This is a plain text", os.Getenv(plainTextKey),
		"The environment variable should be set")
}
