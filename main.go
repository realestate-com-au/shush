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
	app.Version = "1.4.0"
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
				ciphertext, err := (&kms.Handler{
					Service:   kms.Client(c.GlobalString("region")),
					Context:   c.GlobalStringSlice("context"),
					CipherKey: plaintext,
					KeyID:     key,
					Plaintext: plaintext,
				}).Encrypt()
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
				plaintext, err := (&kms.Handler{
					Service:   kms.Client(c.GlobalString("region")),
					Context:   c.GlobalStringSlice("context"),
					CipherKey: ciphertext,
				}).Decrypt()
				sys.CheckError(err, sys.KmsError)
				fmt.Print(plaintext)
			},
		},
		{
			Name:      "encryptssm",
			Usage:     "Encrypt SSM Parameter (kms encryption is optional)",
			UsageText: "shush encryptssm --kms <kms key> <Parameter name> <Parameter Value>",
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
			},
		},
		{
			Name:      "decryptssm",
			Usage:     "Decrypt SSM cipherkey",
			UsageText: "shush decryptssm <Parameter name>",
			Action: func(c *cli.Context) {
				ssmkey, err := sys.GetPayload(c.Args())
				sys.CheckError(err, sys.UsageError)
				plaintext, err := (&ssm.Handler{
					Service:          ssm.Client(c.GlobalString("region")),
					ParameterKeyName: ssmkey,
				}).Decrypt()
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
