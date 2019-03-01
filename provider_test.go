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
		assert.Equal(t, k.custom, KMSPrefix, "Custom prefix overwrite KMSPrefix")
	}
}

func TestEnvDrive(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock_provider.NewMockProvider(ctrl)

	m.EXPECT().KMSDecryptEnv("helloworld", "KMS_ENCRYPTED_ABCD")

	(&envDriver{
		variables: []string{"KMS_ENCRYPTED_ABCD=helloworld"},
	}).drive(m)
}
