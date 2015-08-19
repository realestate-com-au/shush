# Shush!

"shush" is a small tool that can be used to encrypt and decrypt secrets, using the AWS "Key Management Service" (KMS).

## Usage

Encrypt secrets like this:

    shush encrypt KEY-ID < secrets.txt > secrets.encrypted

Later, you can decrypt them:

    shush decrypt < secrets.encrypted > secrets.txt

KEY-ID can be the id or ARN of a KMS master key, or alias prefixed by "alias/".  See documentation on [Encrypt](http://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html) for more details.

Appropriate AWS credentials must be provided by one of the [mechanisms support by aws-sdk-go](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials), e.g. environment variables, or EC2 instance profile.

When used within EC2, `shush` selects the appropriate region automatically.  Outside EC2, you'll need to set `$AWS_DEFAULT_REGION`.

### Limitations

"shush" can only encrypt small amounts of data; up to 4 KB.

## Installation

    $ go get github.com/realestate-com-au/shush

Binaries for official releases may be downloaded from the [releases page on GitHub](https://github.com/realestate-com-au/shush/releases).
