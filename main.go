package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/sys"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "shush"
	app.Version = "1.5.4"
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "trim, t",
					Usage: "If set, remove leading and trailing whitespace from plaintext",
				},
				cli.BoolFlag{
					Name:  "no-warn-whitespace, w",
					Usage: "If set, suppress warnings about whitespace in plaintext",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					sys.Abort(sys.UsageError, "no key specified")
				}
				key := c.Args().First()

				if !isValidUUID(key) && !isArn(key) {
					if !isAlias(key) {
						key = "alias/" + key
					}
				}

				handle, err := kms.NewHandle(
					c.GlobalString("region"),
					c.GlobalStringSlice("context"),
				)
				if err != nil {
					sys.Abort(sys.UsageError, err)
				}
				plaintext, err := sys.GetPayload(c.Args()[1:])
				if err != nil {
					sys.Abort(sys.UsageError, err)
				}
				if c.Bool("trim") {
					plaintext = strings.TrimSpace(plaintext)
				}

				// Warn if input has suspicious whitespace.
				if plaintext != strings.TrimSpace(plaintext) &&
					!c.Bool("no-warn-whitespace") {
					fmt.Fprintf(os.Stderr, "shush: 🚨 WARNING: Plaintext contains suspicious whitespace, consider --trim, or silence with --no-warn-whitespace\n")
				}

				ciphertext, err := handle.Encrypt(plaintext, key)
				if err != nil {
					sys.Abort(sys.KmsError, err)
				}
				fmt.Println(ciphertext)
			},
		},
		{
			Name:  "decrypt",
			Usage: "Decrypt KMS ciphertext",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "print-key",
					Usage: "Print the key instead of the deciphered text",
				},
			},
			Action: func(c *cli.Context) {
				handle, err := kms.NewHandle(
					c.GlobalString("region"),
					c.GlobalStringSlice("context"),
				)
				if err != nil {
					sys.Abort(sys.UsageError, err)
				}
				ciphertext, err := sys.GetPayload(c.Args())
				if err != nil {
					sys.Abort(sys.UsageError, err)
				}
				plaintext, keyId, err := handle.Decrypt(ciphertext)
				if err != nil {
					sys.Abort(sys.KmsError, err)
				}
				if c.Bool("print-key") {
					fmt.Print(keyId)
				} else {
					fmt.Print(plaintext)
				}
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
			SkipArgReorder: true,
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
					handle, err := kms.NewHandle(
						c.GlobalString("region"),
						c.GlobalStringSlice("context"),
					)
					if err != nil {
						sys.Abort(sys.UsageError, err)
					}
					for _, e := range os.Environ() {
						keyValuePair := strings.SplitN(e, "=", 2)
						key := keyValuePair[0]
						if strings.HasPrefix(key, encryptedVarPrefix) {
							ciphertext := keyValuePair[1]
							plaintextKey := key[len(encryptedVarPrefix):len(key)]
							plaintext, _, err := handle.Decrypt(ciphertext)
							if err != nil {
								sys.Abort(sys.KmsError, fmt.Sprintf("cannot decrypt $%s; %s\n", key, err))
							}
							os.Setenv(plaintextKey, plaintext)
						}
					}
				}
				sys.ExecCommand(c.Args())
			},
		},
	}

	app.Run(os.Args)

}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func isArn(u string) bool {
	return strings.HasPrefix(u, "arn:aws:kms")
}

func isAlias(u string) bool {
	return strings.HasPrefix(u, "alias/")
}
