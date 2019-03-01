package main

import (
	"strings"

	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/ssm"
)

const (
	// DefaultSSMPrefix used in environment variables
	DefaultSSMPrefix = "SSM_PS_"
	// DefaultKMSPrefix used in environment variables
	DefaultKMSPrefix = "KMS_ENCRYPTED_"
)

// ProviderIface defines interfaces for all service providers
type ProviderIface interface {
	KMSDecryptEnv(string, string)
	SSMDecryptEnv(string, string)
}

// ProviderImpl structure data for provider
type ProviderImpl struct {
	region   string
	contexts []string
}

// KMSDecryptEnv implement decrypt env
func (pi *ProviderImpl) KMSDecryptEnv(value string, plaintextKey string) {
	(&kms.Handler{
		Service:      kms.Client(pi.region),
		Context:      pi.contexts,
		CipherKey:    value,
		PlaintextKey: plaintextKey,
	}).DecryptEnv()
}

// SSMDecryptEnv implement decrypt env
func (pi *ProviderImpl) SSMDecryptEnv(value string, plaintextKey string) {
	(&ssm.Handler{
		Service:          ssm.Client(pi.region),
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
func (e *envDriver) drive(p ProviderIface) {

	// Ensure CustomPrefix specified
	if e.customPrefix == "" {
		e.customPrefix = DefaultKMSPrefix
	}

	for _, secret := range e.variables {
		keyValuePair := strings.SplitN(secret, "=", 2)
		key := keyValuePair[0]
		value := keyValuePair[1] //KMS encrypted value or SSM parameter key name

		switch {
		case isKMSHandler(secret, e.customPrefix):
			// Update per KMS environment variable
			plaintextKey := key[len(e.customPrefix):len(key)]
			p.KMSDecryptEnv(value, plaintextKey)
		case isSSMHander(secret):
			// Update per SSM environment variable
			plaintextKey := key[len(DefaultSSMPrefix):len(key)]
			p.SSMDecryptEnv(value, plaintextKey)
		}
	}

}

// isSSMHander return true if prefix with SSM_PS_
func isSSMHander(key string) bool {
	return strings.HasPrefix(key, DefaultSSMPrefix)
}

// isKMSHandler return true if prefix with KMS_ENCRYPTED_
func isKMSHandler(key string, customPrefix string) bool {
	return strings.HasPrefix(key, customPrefix)
}
