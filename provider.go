package main

import (
	"strings"

	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/ssm"
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

// Provider defines interfaces for all service providers
type Provider interface {
	Encrypt() (string, error)
	Decrypt() (string, error)
	DecryptEnv()
}

// // execEnv implement update env variable as per service provider
// func execEnv(s Provider) {
// 	s.DecryptEnv()
// }

// // decrypt implement decrypt secret as per service provider
// func decrypt(s Provider) (string, error) {
// 	return s.Decrypt()
// }

// // encrypt implement encrypt secret as per service provider
// func encrypt(s Provider) (string, error) {
// 	return s.Encrypt()
// }

// envDrive structure the data to for environment variable decryptions
type envDriver struct {
	variables, contexts []string
	customPrefix        string // Support KMS custom prefix to be backward compatible
	region              string
}

// driver scans env variables prefix with customPrefix for decryption
func (e *envDriver) drive() {

	for _, secret := range e.variables {
		keyValuePair := strings.SplitN(secret, "=", 2)
		key := keyValuePair[0]
		value := keyValuePair[1] //KMS encrypted value or SSM parameter key name

		switch {
		case isKMSHandler(secret, e.customPrefix):
			// Update per KMS environment variable
			plaintextKey := key[len(KMSPrefix):len(key)]
			(&kms.Handler{
				Service:      kms.Client(e.region),
				Context:      e.contexts,
				Prefix:       KMSPrefix,
				CipherKey:    value,
				PlaintextKey: plaintextKey,
			}).DecryptEnv()
		case isSSMHander(secret):
			// Update per SSM environment variable
			plaintextKey := key[len(SSMPrefix):len(key)]
			(&ssm.Handler{
				Service:          ssm.Client(e.region),
				Prefix:           SSMPrefix,
				ParameterKeyName: value,
				PlaintextKey:     plaintextKey,
			}).DecryptEnv()
		}
	}

}
