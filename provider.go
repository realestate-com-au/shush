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
	KMSDecryptEnv(string, string)
	SSMDecryptEnv(string, string)
}

// ProviderImpl structure data for provider
type ProviderImpl struct {
	region   string
	contexts []string
}

// KMSDecryptEnv implement decrypt env
func (kmsImpl *ProviderImpl) KMSDecryptEnv(value string, plaintextKey string) {
	(&kms.Handler{
		Service:      kms.Client(kmsImpl.region),
		Context:      kmsImpl.contexts,
		Prefix:       KMSPrefix,
		CipherKey:    value,
		PlaintextKey: plaintextKey,
	}).DecryptEnv()
}

// SSMDecryptEnv implement decrypt env
func (ssmImpl *ProviderImpl) SSMDecryptEnv(value string, plaintextKey string) {
	(&ssm.Handler{
		Service:          ssm.Client(ssmImpl.region),
		Prefix:           SSMPrefix,
		ParameterKeyName: value,
		PlaintextKey:     plaintextKey,
	}).DecryptEnv()
}

// envDrive structure the data to for environment variable decryptions
type envDriver struct {
	variables    []string
	customPrefix string // Support KMS custom prefix to be backward compatible
}

// driver scans env variables prefix with customPrefix for decryption
func (e *envDriver) drive(p Provider) {

	for _, secret := range e.variables {
		keyValuePair := strings.SplitN(secret, "=", 2)
		key := keyValuePair[0]
		value := keyValuePair[1] //KMS encrypted value or SSM parameter key name

		switch {
		case isKMSHandler(secret, e.customPrefix):
			// Update per KMS environment variable
			plaintextKey := key[len(KMSPrefix):len(key)]
			p.KMSDecryptEnv(value, plaintextKey)
		case isSSMHander(secret):
			// Update per SSM environment variable
			plaintextKey := key[len(SSMPrefix):len(key)]
			p.SSMDecryptEnv(value, plaintextKey)
		}
	}

}
