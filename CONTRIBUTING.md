# Contributing to ReShifter

Your contributions are welcome, no matter how small or big and in which area: issues, code, docs, a blog post or a tweet.

## Reporting issues

Please provide a description of your environment (cloud/on-premises, K8S version/distro, etcd version) and what
the target is (rcli or the app). See this wonderful example of a helpful [report](https://github.com/mhausenblas/reshifter/issues/2#issue-242929110).

## Pull requests and testing

We're looking forward to your PR, no matter if you spotted a grammar problem in the documentation or if you want to provide new test cases.
Start by forking this repo and we suggest that you create a new branch for each PR.
Note that we have a very comprehensive [test bed](https://github.com/mhausenblas/reshifter/tree/master/testbed) in addition to the unit tests
of the packages, which are executed via `make gtest`. Please do make sure all tests pass before your send in a PR.
