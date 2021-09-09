# Releasing a new version of Shush

When you're ready to perform a new release, please do the following (sorry it's not fully automated):

* Log into an AWS account and locate an enabled KMS key which has an alias
* In your terminal, export the key ID as `SHUSH_KEY` and the alias as `SHUSH_ALIAS`
* Run `./auto/test` with the AWS credentials of the account that has the key

Assuming the tests pass and you want to proceed with the release...

* Change desired version in `main.go`
* Commit and push to master
* On REA's internal CI tool of choice, find the build for `stackup-ci` and trigger a new build from `HEAD`. This will always grab the latest code from here and release the rest of the parts.
