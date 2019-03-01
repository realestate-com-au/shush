package main

import (
	"fmt"
	"strings"

	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/ssm"
	"github.com/realestate-com-au/shush/sys"
	"github.com/urfave/cli"
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
func (pi *ProviderImpl) KMSDecryptEnv(value string, plaintextKey string) {
	(&kms.Handler{
		Service:      kms.Client(pi.region),
		Context:      pi.contexts,
		Prefix:       KMSPrefix,
		CipherKey:    value,
		PlaintextKey: plaintextKey,
	}).DecryptEnv()
}

// SSMDecryptEnv implement decrypt env
func (pi *ProviderImpl) SSMDecryptEnv(value string, plaintextKey string) {
	(&ssm.Handler{
		Service:          ssm.Client(pi.region),
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

func KMSEncrytAction(c *cli.Context) {
	if len(c.Args()) == 0 {
		sys.Abort(sys.UsageError, "no key specified")
	}
	plaintext, err := sys.GetPayload(c.Args()[1:])
	sys.CheckError(err, sys.UsageError)
	key := c.Args().First()
	ciphertext, err := (&kms.Handler{
		Service:   kms.Client(c.GlobalString("region")),
		Context:   c.GlobalStringSlice("context"),
		CipherKey: plaintext,
		KeyID:     key,
		Plaintext: plaintext,
	}).Encrypt()
	sys.CheckError(err, sys.KmsError)
	fmt.Println(ciphertext)
}

func KMSDecryptAction(c *cli.Context) {
	ciphertext, err := sys.GetPayload(c.Args())
	sys.CheckError(err, sys.UsageError)
	plaintext, err := (&kms.Handler{
		Service:   kms.Client(c.GlobalString("region")),
		Context:   c.GlobalStringSlice("context"),
		CipherKey: ciphertext,
	}).Decrypt()
	sys.CheckError(err, sys.KmsError)
	fmt.Print(plaintext)
}

func SSMEncryptAction(c *cli.Context) {
	if len(c.Args()) == 0 {
		sys.Abort(sys.UsageError, "Much specify a parameter key and a value")
	}
	paramVal, err := sys.GetPayload(c.Args()[1:])
	sys.CheckError(err, sys.UsageError)
	output, err := (&ssm.Handler{
		Service:          ssm.Client(c.GlobalString("region")),
		ParameterKeyName: c.Args().First(),
		ParameterValue:   paramVal,
		KMSKeyID:         c.String("kms"),
	}).Encrypt()
	sys.CheckError(err, sys.SsmError)
	fmt.Println(output)
}

func SSMDecryptAction(c *cli.Context) {
	ssmkey, err := sys.GetPayload(c.Args())
	sys.CheckError(err, sys.UsageError)
	plaintext, err := (&ssm.Handler{
		Service:          ssm.Client(c.GlobalString("region")),
		ParameterKeyName: ssmkey,
	}).Decrypt()
	sys.CheckError(err, sys.SsmError)
	fmt.Print(plaintext)
}
