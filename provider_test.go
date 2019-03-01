package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/realestate-com-au/shush/mock_provider"
	"github.com/stretchr/testify/assert"
)

func TestIsSSMHander(t *testing.T) {

	var keys = []struct {
		key      string
		expected bool
	}{
		{"SSM_PS_KEY", true},
		{"KMS_ENCRYPTED_KEY", false},
		{"OTHER", false},
	}
	for _, k := range keys {
		assert.Equal(t, k.expected, isSSMHander(k.key), "SSM should be selected as true")
	}
}

func TestIsKMSHander(t *testing.T) {

	var keys = []struct {
		key      string
		custom   string
		expected bool
	}{
		{"SSM_PS_KEY", "KMS_ENCRYPTED_", false},
		{"KMS_ENCRYPTED_KEY", "KMS_ENCRYPTED_", true},
		{"KMS_OTHER_KEY", "KMS_OTHER_", true},
		{"OTHER_KEY", "KMS_OTHER_", false},
		{"OTHER", "KMS_ENCRYPTED_", false},
	}
	for _, k := range keys {
		assert.Equal(t, k.expected, isKMSHandler(k.key, k.custom), "KMS should be selected as true")
	}
}

func TestEnvDriveKMS(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock_provider.NewMockProvider(ctrl)

	m.EXPECT().KMSDecryptEnv("helloworld", "ABCD")

	(&envDriver{
		variables: []string{"KMS_ENCRYPTED_ABCD=helloworld"},
	}).drive(m)
}

func TestEnvDriveKMSCustomPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock_provider.NewMockProvider(ctrl)

	m.EXPECT().KMSDecryptEnv("okayworld", "EFGH")

	(&envDriver{
		variables: []string{
			"KMS_ENCRYPTED_ABCD=helloworld",
			"KMS_OKAY_EFGH=okayworld",
		},
		customPrefix: "KMS_OKAY_",
	}).drive(m)
}

func TestEnvDriveSSM(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock_provider.NewMockProvider(ctrl)

	m.EXPECT().SSMDecryptEnv("helloworld", "ABCD")

	(&envDriver{
		variables: []string{"SSM_PS_ABCD=helloworld"},
	}).drive(m)
}

func TestEnvDriveCoExist(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock_provider.NewMockProvider(ctrl)

	m.EXPECT().SSMDecryptEnv("helloworld", "ABCD")
	m.EXPECT().KMSDecryptEnv("helloworld", "EDCE")
	// m.EXPECT().KMSDecryptEnv("World", "HELLO")

	(&envDriver{
		variables: []string{
			"SSM_PS_ABCD=helloworld",
			"KMS_ENCRYPTED_EDCE=helloworld",
			"HELLO=World",
		},
	}).drive(m)
}
