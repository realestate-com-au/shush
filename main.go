package main

import (
	// "encoding/base64"
	// "github.com/awslabs/aws-sdk-go/aws"
	// "github.com/awslabs/aws-sdk-go/gen/kms"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "shush"
	app.Version = "1.0.0"
	app.Usage = "KMS encryption and decryption"

	app.Commands = []cli.Command{
		{
			Name:  "encrypt",
			Usage: "Encrypt with a KMS key",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					abort(1, "no key specified")
				}
				key := c.Args().First()
				plaintext, err := getPayload(c.Args()[1:])
				if err != nil {
					abort(1, err)
				}
				fmt.Println("encrypt " + string(plaintext) + " with key " + key)
			},
		},
		{
			Name:  "decrypt",
			Usage: "Decrypt KMS ciphertext",
			Action: func(c *cli.Context) {
				ciphertext, err := getPayload(c.Args())
				if err != nil {
					abort(1, err)
				}
				fmt.Println("decrypt " + string(ciphertext))
			},
		},
	}

	app.Run(os.Args)

}

func getPayload(args []string) ([]byte, error) {
	if len(args) >= 1 {
		return []byte(args[0]), nil
	} else {
		return ioutil.ReadAll(os.Stdin)
	}
}

func abort(status int, message interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s", message)
	os.Exit(status)
}
