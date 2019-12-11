# Releasing a new version of Shush

When you're ready to perform a new release, please do the following (sorry it's not fully automated):

* Change desired version in `main.go`
* Commit and push to master
* Run `auto/release-docker-image`
* Run `auto/binaries`
* Create a new release in Github with your desired release notes, pointed at the version tag created by `auto/release-docker-image`
* Attach the binaries from `target/*`
