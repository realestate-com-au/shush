package sys

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// GetPayload - get input, from command-line (if present) or STDIN.
//
func GetPayload(args []string) (string, error) {
	if len(args) >= 1 {
		return args[0], nil
	}
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(input)), nil
}

// Abort - exit with status and message
//
func Abort(status int, message interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	os.Exit(status)
}

// ExecCommand - exec an external command
//
func ExecCommand(args []string) {
	if len(args) == 0 {
		Abort(UsageError, "no command specified")
	}
	commandName := args[0]
	commandPath, err := exec.LookPath(commandName)
	if err != nil {
		Abort(CommandNotFoundError, fmt.Sprintf("cannot find '%s'\n", commandName))
	}
	err = syscall.Exec(commandPath, args, os.Environ())
	if err != nil {
		Abort(ExecError, err)
	}
}
