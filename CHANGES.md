# CHANGES

## 1.5.5

- When encrypting, warn if plaintext begins or ends with whitespace.  Warning can be silenced with `--no-warn-whitespace` or `-w`.

## 1.5.4

- Upgraded Go from 1.18 to 1.20.
- Upgraded base docker image to alpine:3.18.2.
- Upgraded the AWS SDK KMS component to v1.24.1.
- Upgraded the AWS SDK config component to v1.18.32.
- Upgraded the Google/uuid library to v1.3.0.
- Upgraded the urfave/cli library to v1.22.14.

## 1.5.3

- Upgraded base docker image to alpine:3.16.2
- Upgraded the AWS SDK KMS component to v1.18.4.
- Upgraded Go from 1.16 to 1.18.

## 1.5.2

- Bugfix: fully qualified KMS Key ARNs were being treated as aliases, this meant
  cross account encryption wasn't possible.

## 1.5.1

- Support for darwin/arm64.
- Multi-architecture docker images.

## 1.5.0

- Upgraded the AWS SDK from 1 to 2.
- Upgraded Go from 1.13 to 1.16.

## 1.4.1

- Use build image to reduce Docker target image size from 160MB to approximately 20MB.
- Add ARM64 binary build target.
