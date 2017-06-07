# ReShifter

A cluster admin backup and restore tool for OpenShift:

- etcd is handled via code based on [burry](http://burry.sh)
- projects are handled via `oc export dc`
- other cluster objects like namespaces are handled via [direct API access](https://github.com/kubernetes/client-go)

## Local etcd for testing

See [test/](test/) directory for `etc-*` scripts.
