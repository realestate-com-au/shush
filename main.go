package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"git.realestate.com.au/mwilliams/shush/awsmeta"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "shush"
	app.Version = "1.0.0"
	app.Usage = "KMS encryption and decryption"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "context, C",
			Usage: "encryption context",
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
				encryptionContext, err := parseEncryptionContext(c.GlobalString("context"))
				if err != nil {
					abort(1, err)
				}
				plaintext, err := getPayload(c.Args()[1:])
				if err != nil {
					abort(1, err)
				}
				ciphertext, err := encrypt(plaintext, key, encryptionContext)
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
				encryptionContext, err := parseEncryptionContext(c.GlobalString("context"))
				if err != nil {
					abort(1, err)
				}
				ciphertext, err := getPayload(c.Args())
				if err != nil {
					abort(1, err)
				}
				plaintext, err := decrypt(ciphertext, encryptionContext)
				if err != nil {
					abort(2, err)
				}
				fmt.Print(plaintext)
			},
		},
	}

	app.Run(os.Args)

}

func kmsClient() *kms.KMS {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = awsmeta.GetRegion()
	}
	return kms.New(&aws.Config{Region: region})
}

func encrypt(plaintext string, key string, encryptionContext map[string]*string) (string, error) {
	output, err := kmsClient().Encrypt(&kms.EncryptInput{
		KeyID:             &key,
		EncryptionContext: encryptionContext,
		Plaintext:         []byte(plaintext),
	})
	if err != nil {
		return "", err
	}
	ciphertext := base64.StdEncoding.EncodeToString(output.CiphertextBlob)
	return ciphertext, nil
}

func decrypt(ciphertext string, encryptionContext map[string]*string) (string, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	output, err := kmsClient().Decrypt(&kms.DecryptInput{
		CiphertextBlob: ciphertextBlob,
	})
	if err != nil {
		return "", err
	}
	return string(output.Plaintext), nil
}

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

func parseEncryptionContext(contextString string) (map[string]*string, error) {
	if contextString == "" {
		return map[string]*string{}, nil
	}
	parts := strings.Split(contextString, "=")
	if len(parts) < 2 {
		return map[string]*string{}, errors.New("context must be provided in KEY=VALUE format")
	}
	var context = map[string]*string{
		parts[0]: &parts[1],
	}
	return context, nil
}

func abort(status int, message interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s", message)
	os.Exit(status)
}
