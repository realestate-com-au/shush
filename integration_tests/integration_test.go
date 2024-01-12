package integration_tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

const expectedValue = "superdupersecret"

func TestMain(m *testing.M) {
	// Builds from the latest sources to run our integration test
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}
	make := exec.Command("go", "build", "-o", "shush")
	var stderr bytes.Buffer
	make.Stderr = &stderr
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for shush %v", err)
		fmt.Println(stderr.String())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestItEncryptsAndDecrypts(t *testing.T) {
	key := os.Getenv("SHUSH_KEY")

	if key == "" {
		t.Error("SHUSH_KEY was not found in environment. Please set this to a usable KMS key and re-run the test with appropriate AWS credentials")
		t.FailNow()
	}
	alias := os.Getenv("SHUSH_ALIAS")

	if alias == "" {
		t.Error("SHUSH_ALIAS was not found in environment. Please set this to the alias of the key you specified as SHUSH_KEY and re-run the test")
		t.FailNow()
	}

	if _, err := os.Stat("shush"); os.IsNotExist(err) {
		t.Error("Could not find shush binary")
		t.FailNow()
	}

	keysToCheck := []string{key, alias, fmt.Sprintf("alias/%s", alias)}

	for k := range keysToCheck {
		encryptedValue := encryptValue(t, keysToCheck[k], expectedValue)
		decryptValue(t, encryptedValue)
		validateKey(t, encryptedValue)
	}
}

func validateKey(t *testing.T, secret string) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	decryption := exec.Command(path.Join(dir, "shush"), "decrypt", "--print-key", secret)
	var stdoutDecryption bytes.Buffer
	var stderrDecryption bytes.Buffer
	decryption.Stdout = &stdoutDecryption
	decryption.Stderr = &stderrDecryption
	err = decryption.Run()
	if err != nil {
		t.Errorf("Failed to print the key: %v\n%s", err, stderrDecryption.String())
	}

	keyArn := stdoutDecryption.String()
	keyId := os.Getenv("SHUSH_KEY")

	if !strings.HasSuffix(keyArn, keyId) {
		t.Errorf("Expected '%s' to end with %s", keyArn, keyId)
	}
}

func decryptValue(t *testing.T, secret string) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	decryption := exec.Command(path.Join(dir, "shush"), "decrypt", secret)
	var stdoutDecryption bytes.Buffer
	var stderrDecryption bytes.Buffer
	decryption.Stdout = &stdoutDecryption
	decryption.Stderr = &stderrDecryption
	err = decryption.Run()
	if err != nil {
		t.Errorf("Failed to decrypt: %v\n%s", err, stderrDecryption.String())
	}

	if stdoutDecryption.String() != expectedValue {
		t.Errorf("Expected '%s' but got '%s'", expectedValue, stdoutDecryption.String())
	}

	return stdoutDecryption.String()
}

func encryptValue(t *testing.T, key string, secret string) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	encryption := exec.Command(path.Join(dir, "shush"), "encrypt", key, secret)
	t.Logf("Running command 'shush encrypt %s %s'", key, secret)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	encryption.Stdout = &stdout
	encryption.Stderr = &stderr
	err = encryption.Run()
	if err != nil {
		t.Log(stdout.String())
		t.Errorf("Failed to encrypt with key '%s': %v\n%s", key, err, stderr.String())
		t.FailNow()
	}

	return stdout.String()
}
