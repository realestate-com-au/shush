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
	"github.com/codegangsta/cli"
	"github.com/realestate-com-au/shush/awsmeta"
)

func main() {

	app := cli.NewApp()
	app.Name = "shush"
	app.Version = "1.1.0"
	app.Usage = "KMS encryption and decryption"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "context, C",
			Usage: "encryption context",
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
					abort(1, "no key specified")
				}
				key := c.Args().First()
				handle, err := makeKmsHandle(
					c.GlobalString("region"),
					c.GlobalString("context"),
				)
				if err != nil {
					abort(1, err)
				}
				plaintext, err := getPayload(c.Args()[1:])
				if err != nil {
					abort(1, err)
				}
				ciphertext, err := handle.Encrypt(plaintext, key)
				if err != nil {
					abort(2, err)
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
					c.GlobalString("context"),
				)
				if err != nil {
					abort(1, err)
				}
				ciphertext, err := getPayload(c.Args())
				if err != nil {
					abort(1, err)
				}
				plaintext, err := handle.Decrypt(ciphertext)
				if err != nil {
					abort(2, err)
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
				handle, err := makeKmsHandle(
					c.GlobalString("region"),
					c.GlobalString("context"),
				)
				if err != nil {
					abort(1, err)
				}
				for _, e := range os.Environ() {
					pair := strings.Split(e, "=")
					key := pair[0]
					if strings.HasPrefix(key, encryptedVarPrefix) {
						ciphertext := pair[1]
						plaintextKey := key[len(encryptedVarPrefix):len(key)]
						plaintext, err := handle.Decrypt(ciphertext)
						if err != nil {
							abort(2, fmt.Sprintf("cannot decrypt $%s\n", key))
						}
						os.Setenv(plaintextKey, plaintext)
					}
				}
				if len(c.Args()) == 0 {
					abort(1, "no command specified")
				}
				commandName := c.Args().First()
				commandPath, err := exec.LookPath(commandName)
				if err != nil {
					abort(3, fmt.Sprintf("cannot find $%s\n", commandName))
				}
				err = syscall.Exec(commandPath, c.Args(), os.Environ())
				if err != nil {
					abort(3, err)
				}
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

func makeKmsHandle(region string, contextString string) (ops *kmsHandle, err error) {
	encryptionContext, err := parseEncryptionContext(contextString)
	if err != nil {
		return
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

func parseEncryptionContext(contextString string) (kmsEncryptionContext, error) {
	if contextString == "" {
		return kmsEncryptionContext{}, nil
	}
	parts := strings.Split(contextString, "=")
	if len(parts) < 2 {
		return kmsEncryptionContext{}, errors.New("context must be provided in KEY=VALUE format")
	}
	var context = kmsEncryptionContext{
		parts[0]: &parts[1],
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
	fmt.Fprintf(os.Stderr, "ERROR: %s", message)
	os.Exit(status)
}
