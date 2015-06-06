# Shush!

"shush" is a small tool that can be used to encrypt and decrypt secrets, using the AWS "Key Management Service" (KMS).

## Usage

Encrypt secrets like this:

    shush encrypt KEY-ID < secrets.txt > secrets.encrypted

Later, you can decrypt them:

    shush decrypt < secrets.encrypted > secrets.txt

KEY-ID can be the id or ARN of a KMS master key, or alias prefixed by "alias/".  See document on [Encrypt](http://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html) for more details.

Naturally, suitable AWS credentials must be provided (via environment variables, command-line options, or EC2 instance profile).

## Limitations

"shush" can only encrypt small amounts of data; up to 4 KB.
