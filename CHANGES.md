# CHANGES

## 1.5.3

- Upgraded the AWS SDK KMS component to v1.18.4.

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
