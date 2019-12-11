# Releasing a new version of Shush

When you're ready to perform a new release, please do the following (sorry it's not fully automated):

* Log into an AWS account and locate an enabled KMS key which has an alias
* In your terminal, export the key ID as `SHUSH_KEY` and the alias as `SHUSH_ALIAS`
* Run `make test` with the AWS credentials of the account that has the key

Assuming the tests pass and you want to proceed with the release...

* Change desired version in `main.go`
* Commit and push to master
* Log into docker hub `docker login`. Your Docker ID needs to be in the `realestate` organisation on Docker hub.
* Run `auto/release-docker-image`
* Run `auto/binaries`
* Create a new release in Github with your desired release notes, pointed at the version tag created by `auto/release-docker-image`
* Attach the binaries from `target/*`
