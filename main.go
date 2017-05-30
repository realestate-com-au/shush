package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/realestate-com-au/shush/awsmeta"
	"github.com/urfave/cli"
)

const usageError = 64            // incorrect usage of "shush"
const kmsError = 69              // KMS encrypt/decrypt issues
const execError = 126            // cannot execute the specified command
const commandNotFoundError = 127 // cannot find the specified command

func main() {

	app := cli.NewApp()
	app.Name = "shush"
	app.Version = "1.3.2"
	app.Usage = "KMS encryption and decryption"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "context, C",
			Usage:  "encryption context",
			EnvVar: "KMS_ENCRYPTION_CONTEXT",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "AWS region",
			EnvVar: "AWS_DEFAULT_REGION",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "encrypt",
			Usage: "Encrypt with a KMS key",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					abort(usageError, "no key specified")
				}
				key := c.Args().First()
				handle, err := makeKmsHandle(
					c.GlobalString("region"),
					c.GlobalStringSlice("context"),
				)
				if err != nil {
					abort(usageError, err)
				}
				plaintext, err := getPayload(c.Args()[1:])
				if err != nil {
					abort(usageError, err)
				}
				ciphertext, err := handle.Encrypt(plaintext, key)
				if err != nil {
					abort(kmsError, err)
				}
				fmt.Println(ciphertext)
			},
		},
		{
			Name:  "decrypt",
			Usage: "Decrypt KMS ciphertext",
			Action: func(c *cli.Context) {
				handle, err := makeKmsHandle(
					c.GlobalString("region"),
					c.GlobalStringSlice("context"),
				)
				if err != nil {
					abort(usageError, err)
				}
				ciphertext, err := getPayload(c.Args())
				if err != nil {
					abort(usageError, err)
				}
				plaintext, err := handle.Decrypt(ciphertext)
				if err != nil {
					abort(kmsError, err)
				}
				fmt.Print(plaintext)
			},
		},
		{
			Name:  "exec",
			Usage: "Execute a command",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix",
					Usage: "environment variable prefix",
					Value: "KMS_ENCRYPTED_",
				},
			},
			Action: func(c *cli.Context) {
				encryptedVarPrefix := c.String("prefix")
				foundEncrypted := false
				for _, e := range os.Environ() {
					if strings.HasPrefix(e, encryptedVarPrefix) {
						foundEncrypted = true
						break
					}
				}
				if foundEncrypted {
					handle, err := makeKmsHandle(
						c.GlobalString("region"),
						c.GlobalStringSlice("context"),
					)
					if err != nil {
						abort(usageError, err)
					}
					for _, e := range os.Environ() {
						keyValuePair := strings.SplitN(e, "=", 2)
						key := keyValuePair[0]
						if strings.HasPrefix(key, encryptedVarPrefix) {
							ciphertext := keyValuePair[1]
							plaintextKey := key[len(encryptedVarPrefix):len(key)]
							plaintext, err := handle.Decrypt(ciphertext)
							if err != nil {
								abort(kmsError, fmt.Sprintf("cannot decrypt $%s; %s\n", key, err))
							}
							os.Setenv(plaintextKey, plaintext)
						}
					}
				}
				execCommand(c.Args())
			},
		},
	}

	app.Run(os.Args)

}

type kmsEncryptionContext map[string]*string

// Structure encapsulating stuff common to encrypt and decrypt.
//
type kmsHandle struct {
	Client  *kms.KMS
	Context kmsEncryptionContext
}

func makeKmsHandle(region string, context []string) (ops *kmsHandle, err error) {
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
	ops = &kmsHandle{
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
		context[parts[0]] = &parts[1]
	}
	return context, nil
}

// Encrypt plaintext using specified key.
func (h *kmsHandle) Encrypt(plaintext string, keyID string) (string, error) {
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
func (h *kmsHandle) Decrypt(ciphertext string) (string, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(ciphertext)
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

// Get input, from command-line (if present) or STDIN.
func getPayload(args []string) (string, error) {
	if len(args) >= 1 {
		return args[0], nil
	}
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(input), nil
}

func abort(status int, message interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	os.Exit(status)
}

func execCommand(args []string) {
	if len(args) == 0 {
		abort(usageError, "no command specified")
	}
	commandName := args[0]
	commandPath, err := exec.LookPath(commandName)
	if err != nil {
		abort(commandNotFoundError, fmt.Sprintf("cannot find '%s'\n", commandName))
	}
	err = syscall.Exec(commandPath, args, os.Environ())
	if err != nil {
		abort(execError, err)
	}
}
