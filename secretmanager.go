package main

import (
	"log"
	"os"
	"strings"

	"github.com/realestate-com-au/shush/ssm"
	"github.com/realestate-com-au/shush/sys"
)

const (
	defaultSSMPrefix = "SSM_PS_"
	defaultKMSPrefix = "KMS_ENCRYPTED_"
)

func isSSMHander(key string) bool {

	return strings.HasPrefix(key, defaultSSMPrefix)
}

func isKMSHandler(key string, customPrefix string) bool {
	if customPrefix != "" {
		return strings.HasPrefix(key, customPrefix)
	}

	return strings.HasPrefix(key, defaultKMSPrefix)
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
			log.Printf("Found kms: %v", secret)
		case isSSMHander(secret):
			plaintextKey := key[len(defaultSSMPrefix):len(key)]
			handle, err := ssm.New(
				region,
				defaultSSMPrefix,
				ciphertext,
				plaintextKey,
			)
			sys.CheckError(err, sys.SsmError)
			execEnv(handle)

		}
	}
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
