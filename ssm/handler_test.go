package ssm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	"github.com/realestate-com-au/shush/ssm/mock_ssm"
)

func TestClient(t *testing.T) {
	output := Client("ap-southeast-2")
	var s *ssm.SSM
	assert.IsType(t, s, output, "Client should be SSM")
	assert.Equal(t, "https://ssm.ap-southeast-2.amazonaws.com", output.Endpoint, "Should be ssm API endpoint")
}
func TestEncrypt(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		parameterKeyName = "KEY-NAME"
		parameterType    = "String"
		parameterValue   = "ThisIsNotASecret"
		overwrite        = true
	)
	defer ctrl.Finish()

	m := mock_ssm.NewMockAWSIface(ctrl)

	m.EXPECT().PutParameter(gomock.Eq(&ssm.PutParameterInput{
		Name:      &(parameterKeyName),
		Type:      &(parameterType),
		Value:     &(parameterValue),
		Overwrite: &overwrite,
	}))

	(&Handler{
		Service:          m,
		ParameterKeyName: parameterKeyName,
		ParameterValue:   parameterValue,
	}).Encrypt()
}

func TestEncryptbyKMS(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		parameterKeyName = "KEY-NAME"
		parameterType    = "SecureString"
		parameterValue   = "ThisIsASecret"
		kmsKeyID         = "ABCDEFG"
		overwrite        = true
	)
	defer ctrl.Finish()

	m := mock_ssm.NewMockAWSIface(ctrl)

	m.EXPECT().PutParameter(gomock.Eq(&ssm.PutParameterInput{
		Name:      &(parameterKeyName),
		KeyId:     &(kmsKeyID),
		Type:      &(parameterType),
		Value:     &(parameterValue),
		Overwrite: &overwrite,
	}))

	(&Handler{
		Service:          m,
		KMSKeyID:         kmsKeyID,
		ParameterKeyName: parameterKeyName,
		ParameterValue:   parameterValue,
	}).Encrypt()
}

func TestDecrypt(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		parameterKeyName = "KEY-NAME"
		withDecryption   = true
	)
	defer ctrl.Finish()

	m := mock_ssm.NewMockAWSIface(ctrl)

	m.EXPECT().GetParameter(gomock.Eq(&ssm.GetParameterInput{
		Name:           &(parameterKeyName),
		WithDecryption: &withDecryption,
	})).Return(mock_ssm.CustomMockParameterOutput(), nil)

	output, _ := (&Handler{
		Service:          m,
		ParameterKeyName: parameterKeyName,
	}).Decrypt()

	assert.Equal(t, "This is a secret", output, "The parameter value should be equal")
}

func TestDecryptEnv(t *testing.T) {
	ctrl := gomock.NewController(t)

	var (
		parameterKeyName = "SSM_PS_KEY-NAME"
		plainTextKey     = "KEY-NAME"
		withDecryption   = true
	)

	defer ctrl.Finish()

	m := mock_ssm.NewMockAWSIface(ctrl)

	m.EXPECT().GetParameter(gomock.Eq(&ssm.GetParameterInput{
		Name:           &(parameterKeyName),
		WithDecryption: &withDecryption,
	})).Return(mock_ssm.CustomMockParameterOutput(), nil)

	(&Handler{
		Service:          m,
		PlaintextKey:     plainTextKey,
		ParameterKeyName: parameterKeyName,
	}).DecryptEnv()

	assert.Equal(t, "This is a secret", os.Getenv(plainTextKey),
		"The environment variable should be set")
}
