package main

import (
	"log"
	"os"
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
	Decrypt() (string, error)
	DecryptEnv()
}

// driver scans all env variables to decrypt
func driver(variables []string, region string, encryptedVarPrefix string, contexts []string) {

	for _, secret := range variables {

		keyValuePair := strings.SplitN(secret, "=", 2)
		key := keyValuePair[0]
		ciphertext := keyValuePair[1]

		switch {
		case isKMSHandler(secret, encryptedVarPrefix):

			plaintextKey := key[len(KMSPrefix):len(key)]
			handle, err := kms.New(
				region,
				contexts,
				KMSPrefix,
				ciphertext,
				plaintextKey,
			)
			sys.CheckError(err, sys.KmsError)
			execEnv(handle)
		case isSSMHander(secret):

			plaintextKey := key[len(SSMPrefix):len(key)]
			handle, err := ssm.New(
				region,
				SSMPrefix,
				ciphertext,
				plaintextKey,
			)
			sys.CheckError(err, sys.SsmError)
			execEnv(handle)

		}
	}
	log.Printf("Updated env KMS: %v", os.Getenv("ABCD"))
	log.Printf("Updated env DUY: %v", os.Getenv("DUY"))
}

// execEnv implement update env variable as per secret provider
func execEnv(s SecretManager) {
	s.DecryptEnv()
}

// decrypt implement decrypt secret as per secret provider
func decrypt(s SecretManager) (string, error) {
	return s.Decrypt()
}
