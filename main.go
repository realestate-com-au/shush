package main

import (
	"os"

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
			Name:   "encrypt",
			Usage:  "Encrypt with a KMS key",
			Action: KMSEncrytAction,
		},
		{
			Name:   "decrypt",
			Usage:  "Decrypt KMS ciphertext",
			Action: KMSDecryptAction,
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
			Action: SSMEncryptAction,
		},
		{
			Name:      "decryptssm",
			Usage:     "Decrypt SSM cipherkey",
			UsageText: "shush decryptssm <Parameter name>",
			Action:    SSMDecryptAction,
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
					customPrefix: c.String("prefix"),
				}).drive(&ProviderImpl{
					region:   c.GlobalString("region"),
					contexts: c.GlobalStringSlice("context"),
				})
				sys.ExecCommand(c.Args())
			},
		},
	}

	app.Run(os.Args)

}
