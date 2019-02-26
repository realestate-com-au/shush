package sys

// UsageError indicate incorrect usage of "shush"
const UsageError = 64

// KmsError indicate KMS encrypt/decrypt issues
const KmsError = 69

// SsmError indicate SSM decrypt issues
const SsmError = 70 // SSM decrypt issues

// ExecError indicate error to execute the command
const ExecError = 126

// CommandNotFoundError indicate cannot find the specified command
const CommandNotFoundError = 127

// CheckError abort the service
func CheckError(err error, code int) {
	if err != nil {
		Abort(code, err)
	}
}
