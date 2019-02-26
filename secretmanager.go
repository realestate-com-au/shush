package main

import (
	"strings"

	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/ssm"
	"github.com/realestate-com-au/shush/sys"
)

var (
	// SSMPrefix the default
	SSMPrefix = "SSM_PS_"
	// KMSPrefix the default
	KMSPrefix = "KMS_ENCRYPTED_"
)

func isSSMHander(key string) bool {
	return strings.HasPrefix(key, SSMPrefix)
}

func isKMSHandler(key string, customPrefix string) bool {

	if customPrefix != KMSPrefix {
		KMSPrefix = customPrefix
		return strings.HasPrefix(key, customPrefix)
	}

	return strings.HasPrefix(key, KMSPrefix)
}

// SecretManager defines interfaces for all secret providers
type SecretManager interface {
	Encrypt() (string, error)
	Decrypt() (string, error)
	DecryptEnv()
}

// execEnv implement update env variable as per secret provider
func execEnv(s SecretManager) {
	s.DecryptEnv()
}

// decrypt implement decrypt secret as per secret provider
func decrypt(s SecretManager) (string, error) {
	return s.Decrypt()
}

// encrypt implement encrypt secret as per secret provider
func encrypt(s SecretManager) (string, error) {
	return s.Encrypt()
}

// envDrive structure the data to for environment variable decryptions
type envDrive struct {
	variables, contexts []string
	encryptedVarPrefix  string
	region              string
}

// driver scans all env variables to decrypt
func (e *envDrive) driver() {

	for _, secret := range e.variables {
		keyValuePair := strings.SplitN(secret, "=", 2)
		key := keyValuePair[0]
		ciphertext := keyValuePair[1]

		switch {
		case isKMSHandler(secret, e.encryptedVarPrefix):
			// Update per KMS environment variable
			plaintextKey := key[len(KMSPrefix):len(key)]
			c, err := kms.Client(e.region)
			sys.CheckError(err, sys.KmsError)
			encryptionContext, err := kms.ParseEncryptionContext(e.contexts)
			sys.CheckError(err, sys.KmsError)
			execEnv(&kms.Handle{
				Client:       c,
				Context:      encryptionContext,
				Prefix:       KMSPrefix,
				CipherKey:    ciphertext,
				PlaintextKey: plaintextKey,
			})
		case isSSMHander(secret):
			// Update per SSM environment variable
			plaintextKey := key[len(SSMPrefix):len(key)]
			c, err := ssm.Client(e.region)
			sys.CheckError(err, sys.SsmError)
			execEnv(&ssm.Handle{
				Client:       c,
				Prefix:       SSMPrefix,
				CipherKey:    ciphertext,
				PlaintextKey: plaintextKey,
			})
		}
	}
}
