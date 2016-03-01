# Shush!

"shush" is a small tool that can be used to encrypt and decrypt secrets, using the AWS "Key Management Service" (KMS).

## Usage

### Encrypting things

Encrypt secrets like this:

    shush encrypt KEY-ID < secret.txt > secret.encrypted

The output of `encrypt` is Base64-encoded ciphertext.

KEY-ID can be the id or ARN of a KMS master key, or alias prefixed by "alias/".  See documentation on [Encrypt](http://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html) for more details.

Plaintext input can also be provided on the command-line, e.g.

    shush encrypt KEY-ID 'this is a secret' > secret.encrypted

### Decrypting things

Encrypted secrets are easy to decrypt, like this:

    shush decrypt < secret.encrypted > secret.txt

There's no need to specify a KEY-ID here, as it's encoded in the ciphertext.

### Credentials and region

Appropriate AWS credentials must be provided by one of the [mechanisms supported by aws-sdk-go](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials), e.g. environment variables, or EC2 instance profile.

When used within EC2, `shush` selects the appropriate region automatically.  
Outside EC2, you'll need to specify it, via `--region` or by setting `$AWS_DEFAULT_REGION`.

### Encryption context

"shush" supports KMS [encryption contexts](http://docs.aws.amazon.com/kms/latest/developerguide/encryption-context.html), which may be used to scope use of a key.  The same encryption context must be provided when encrypting and decrypting.

    shush --context app=myapp encrypt KEY-ID secret.txt > secret.encrypted
    shush --context app=myapp decrypt < secret.encrypted > secret.txt

### Limitations

"shush" can only encrypt small amounts of data; up to 4 KB.

## Installation

Binaries for official releases may be downloaded from the [releases page on GitHub](https://github.com/realestate-com-au/shush/releases).

If you want to compile it from source, try:

    $ go get github.com/realestate-com-au/shush

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

## See also

If you dislike 8Gb binary files, and happen to have a Ruby interpreter handy,
"ssssh" is a drop-in replacement for "shush":

* https://github.com/mdub/ssssh

Or, you can just use `bash`, `base64`, and the AWS CLI:

    base64 -d < secrets.encrypted > /tmp/secrets.bin
    aws kms decrypt --ciphertext-blob fileb:///tmp/secrets.bin --output text --query Plaintext | base64 -d > secrets.txt

## License

Copyright (c) 2015 REA Group Ltd.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

## Contributing

Source-code for shush is [on Github](https://github.com/realestate-com-au/shush).
