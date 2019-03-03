# Shush!

"shush" is a small tool that can be used to encrypt and decrypt secrets, using the AWS "Key Management Service" (KMS).

## Usage

### Encrypting things

**KMS**

Encrypt secrets like this:

    shush encrypt KEY-ID < secret.txt > secret.encrypted

The output of `encrypt` is Base64-encoded ciphertext.

KEY-ID can be the id or ARN of a KMS master key, or alias prefixed by "alias/".  See documentation on [Encrypt](http://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html) for more details.

Plaintext input can also be provided on the command-line, e.g.

    shush encrypt KEY-ID 'this is a secret' > secret.encrypted

**SSM Parameter Store**

Create a parameter without encryption:

    shush encryptssm PARAMETER-NAME <value>

Create and encrypt a parameter like this:

    shush encryptssm --kms KEY-ID PARAMETER-NAME <value>

### Decrypting things

Encrypted secrets are easy to decrypt, like this:

**KMS**

    shush decrypt < secret.encrypted > secret.txt

**SSM Parameter Store**

    shush decryptssm PARAMETER-NAME

There's no need to specify a KMS KEY-ID here, as it's encoded in the ciphertext.

### Credentials and region

Appropriate AWS credentials must be provided by one of the [mechanisms supported by aws-sdk-go](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials), e.g. environment variables, or EC2 instance profile.

When used within EC2, `shush` selects the appropriate region automatically.  
Outside EC2, you'll need to specify it, via `--region` or by setting `$AWS_DEFAULT_REGION`.

### Encryption context

"shush" supports KMS [encryption contexts](http://docs.aws.amazon.com/kms/latest/developerguide/encryption-context.html), which may be used to scope use of a key.  The same encryption context must be provided when encrypting and decrypting.

    shush --context app=myapp encrypt KEY-ID secret.txt > secret.encrypted
    shush --context app=myapp decrypt < secret.encrypted > secret.txt

SSM Parameter store feature does not support Encryption Contexts, yet.

### Limitations

"shush" can only encrypt small amounts of data; up to 4 KB.

## Use as a command shim

"shush exec" is a command shim that makes it easy to use secrets transported
via the (Unix) process environment.  It decrypts any environment variables
prefixed by "`KMS_ENCRYPTED_`", and executes a specified command.

`SSM_PS_` for SSM parameter store environment variables

For example:

```
$ export KMS_ENCRYPTED_DB_PASSWORD=$(shush encrypt alias/app-secrets 'seasame')
$ shush exec -- env | grep PASSWORD
KMS_ENCRYPTED_DB_PASSWORD=CiAbQLOo2VC4QTV/Ng986wsDSJ0srAe6oZnLyzRT6pDFWRKOAQEBAgB4G0CzqNlQuEE1fzYPfOsLA0idLKwHuqGZy8s0U+qQxVkAAABlMGMGCSqGSIb3DQEHBqBWMFQCAQAwTwYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAzfFR0tsHRq18JUhMcCARCAImvuMNYuHUut3BT7sZs9a31qWcmOBUBXYEsD+kx2GxUxBPE=
DB_PASSWORD=seasame
```

In this example, "shush exec":

- found `$KMS_ENCRYPTED_DB_PASSWORD` in the environment
- decrypted the contents
- put the result in `$DB_PASSWORD`
- executed `env`

"shush exec" works well as an entrypoint for Docker images, e.g.

    # Include "shush" to decode KMS_ENCRYPTED_STUFF
    RUN curl -sL -o /usr/local/bin/shush \
        https://github.com/realestate-com-au/shush/releases/download/v1.4.0/shush_linux_amd64 \
     && chmod +x /usr/local/bin/shush
    ENTRYPOINT ["/usr/local/bin/shush", "exec", "--"]

Use the same command "shush exec" for SSM parameter store decryption.


**Ec2 instance IAM**

```yaml
iam_policy_statements:
- Effect: Allow
  Action:
  - ssm:GetParameter
  Resource:
  - arn:aws:ssm:ap-southeast-2:xxxxxxxxxxxx:parameter/KEY-NAME-1
  - arn:aws:ssm:ap-southeast-2:xxxxxxxxxxxx:parameter/KEY-NAME-2
```

## Installation

Binaries for official releases may be downloaded from the [releases page on GitHub](https://github.com/realestate-com-au/shush/releases).

If you want to compile it from source, try:

    $ go get github.com/realestate-com-au/shush
    
For Unix/Linux users, you can install `shush` using the following command. You may want to change the version number in the command below from `v1.4.0` to whichever version you want:

```
curl -sL -o /usr/local/bin/shush \
    https://github.com/realestate-com-au/shush/releases/download/v1.4.0/shush_linux_amd64 \
 && chmod +x /usr/local/bin/shush
```

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
