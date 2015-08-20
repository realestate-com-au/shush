# Shush!

"shush" is a small tool that can be used to encrypt and decrypt secrets, using the AWS "Key Management Service" (KMS).

## Usage

Encrypt secrets like this:

    shush encrypt KEY-ID < secrets.txt > secrets.encrypted

Later, you can decrypt them:

    shush decrypt < secrets.encrypted > secrets.txt

KEY-ID can be the id or ARN of a KMS master key, or alias prefixed by "alias/".  See documentation on [Encrypt](http://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html) for more details.

Appropriate AWS credentials must be provided by one of the [mechanisms supported by aws-sdk-go](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials), e.g. environment variables, or EC2 instance profile.

When used within EC2, `shush` selects the appropriate region automatically.  
Outside EC2, you'll need to specify is, via `--region` or by setting `$AWS_DEFAULT_REGION`.

### Limitations

"shush" can only encrypt small amounts of data; up to 4 KB.

## Installation

    $ go get github.com/realestate-com-au/shush

Binaries for official releases may be downloaded from the [releases page on GitHub](https://github.com/realestate-com-au/shush/releases).

## Examples

### Encrypt a password

Encrypt user input:

    echo -n "Enter password: "
    ENCRYPTED_PASSWORD=$(shush encrypt alias/app-secrets)

and later:

    some-command --password $(shush decrypt "$ENCRYPTED_PASSWORD")

### Bulk encryption of secrets

Encrypt some environment variables, as though they were arguments to `env(1)`:

    shush encrypt alias/app-secrets 'FOO=1 BAR=2' > secrets

and later:

    env $(shush decrypt < secrets) some-command
