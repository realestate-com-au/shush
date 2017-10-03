package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/realestate-com-au/shush/kms"
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
				handle, err := kms.NewHandle(
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
				handle, err := kms.NewHandle(
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
					handle, err := kms.NewHandle(
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
