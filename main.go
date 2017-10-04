package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/sys"
	"github.com/urfave/cli"
)

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
					sys.Abort(sys.UsageError, "no key specified")
				}
				key := c.Args().First()
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
				plaintext, err := handle.Decrypt(ciphertext)
				if err != nil {
					sys.Abort(sys.KmsError, err)
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
							plaintext, err := handle.Decrypt(ciphertext)
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
