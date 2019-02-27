package main

import (
	"fmt"
	"os"

	"github.com/realestate-com-au/shush/kms"
	"github.com/realestate-com-au/shush/ssm"
	"github.com/realestate-com-au/shush/sys"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "shush"
	app.Version = "1.3.5"
	app.Usage = "KMS & SSM Parameter Store encryption and decryption"

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
				plaintext, err := sys.GetPayload(c.Args()[1:])
				sys.CheckError(err, sys.UsageError)
				key := c.Args().First()
				kc, err := kms.Client(c.GlobalString("region"))
				sys.CheckError(err, sys.KmsError)
				encryptionContext, err := kms.ParseEncryptionContext(c.GlobalStringSlice("context"))
				sys.CheckError(err, sys.KmsError)
				ciphertext, err := encrypt(&kms.Handler{
					Client:    kc,
					Context:   encryptionContext,
					CipherKey: plaintext,
					KeyID:     key,
					Plaintext: plaintext,
				})
				sys.CheckError(err, sys.KmsError)
				fmt.Println(ciphertext)
			},
		},
		{
			Name:  "decrypt",
			Usage: "Decrypt KMS ciphertext",
			Action: func(c *cli.Context) {
				ciphertext, err := sys.GetPayload(c.Args())
				sys.CheckError(err, sys.UsageError)
				kc, err := kms.Client(c.GlobalString("region"))
				sys.CheckError(err, sys.KmsError)
				encryptionContext, err := kms.ParseEncryptionContext(c.GlobalStringSlice("context"))
				sys.CheckError(err, sys.KmsError)
				plaintext, err := decrypt(&kms.Handler{
					Client:    kc,
					Context:   encryptionContext,
					CipherKey: ciphertext,
				})
				sys.CheckError(err, sys.KmsError)
				fmt.Print(plaintext)
			},
		},
		{
			Name:      "encryptssm",
			Usage:     "Encrypt SSM Parameter (kms encryption is optional)",
			UsageText: "shush encryptssm --kms <kms key> <Parameter key name> <Parameter Value>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "kms",
					Usage: "Use KMS to encrypt the parameter",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					sys.Abort(sys.UsageError, "Much specify a parameter key and a value")
				}
				sc, err := ssm.Client(c.GlobalString("region"))
				sys.CheckError(err, sys.SsmError)
				paramVal, err := sys.GetPayload(c.Args()[1:])
				sys.CheckError(err, sys.UsageError)
				output, err := encrypt(&ssm.Handler{
					Client:           sc,
					ParameterKeyName: c.Args().First(),
					ParameterValue:   paramVal,
					KMSKeyID:         c.String("kms"),
				})
				sys.CheckError(err, sys.SsmError)
				fmt.Println(output)
			},
		},
		{
			Name:      "decryptssm",
			Usage:     "Decrypt SSM cipherkey",
			UsageText: "shush decryptssm <Parameter key name>",
			Action: func(c *cli.Context) {
				ssmkey, err := sys.GetPayload(c.Args())
				sys.CheckError(err, sys.UsageError)
				sc, err := ssm.Client(c.GlobalString("region"))
				sys.CheckError(err, sys.SsmError)
				plaintext, err := decrypt(&ssm.Handler{
					Client:           sc,
					ParameterKeyName: ssmkey,
				})
				sys.CheckError(err, sys.SsmError)
				fmt.Print(plaintext)
			},
		},
		{
			Name:  "exec",
			Usage: "Execute a command",
			Flags: []cli.Flag{
				cli.StringFlag{
					// Support KMS custom prefix to be backward compatible
					Name:  "prefix",
					Usage: "additional environment variable prefix",
					Value: KMSPrefix,
				},
			},
			SkipArgReorder: true,
			Action: func(c *cli.Context) {
				(&envDriver{
					variables:    os.Environ(),
					contexts:     c.GlobalStringSlice("context"),
					region:       c.GlobalString("region"),
					customPrefix: c.String("prefix"),
				}).drive()
				sys.ExecCommand(c.Args())
			},
		},
	}

	app.Run(os.Args)

}
