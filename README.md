# ReShifter

[![Go Report Card](https://goreportcard.com/badge/github.com/mhausenblas/reshifter)](https://goreportcard.com/report/github.com/mhausenblas/reshifter)
[![godoc](https://godoc.org/github.com/mhausenblas/reshifter?status.svg)](https://godoc.org/github.com/mhausenblas/reshifter)
[![Docker Repository on Quay](https://quay.io/repository/mhausenblas/reshifter/status "Docker Repository on Quay")](https://quay.io/repository/mhausenblas/reshifter)

A cluster admin backup and restore tool for Kubernetes distros such as OpenShift, using etcd to query and manipulate the state of all objects.

[![Screen cast: Introducing ReShifter](images/reshifter-cli.png)](https://www.useloom.com/share/e590aedeb95b441fb23ab4f9e9e80c32 "Introducing ReShifter")  

Supported:

- Cluster: Kubernetes 1.5 compatible distros
- App: modern browsers

Index:

- [Using it](#using-it)
  - [Deploy it on OpenShift](#deploy-it-on-openshift)
  - [Deploy it on vanilla Kubernetes](#deploy-it-on-vanilla-kubernetes)
  - [HTTP API](#http-api)
- [Testbed](#testbed)
  - [End-to-end tests](#end-to-end-tests)
  - [Synthetic tests](#synthetic-tests)
  - [Cluster dumps](#cluster-dumps)
- [Extending it](#extending-it)
  - [Vendoring](#vendoring)
  - [Unit tests](#unit-tests)

### Status and roadmap

See [Trello board](https://trello.com/b/iOrEdJQ3/reshifter).

## Using it

### Deploy it on OpenShift

Note: requires an OpenShift 1.5 cluster and [oc](https://github.com/openshift/origin/releases/tag/v1.5.1) installed, locally.

```
$ make init
$ make publish
```

### Deploy it on vanilla Kubernetes

TBD.

### HTTP API

```
GET  /                    … the ReShifter Web UI
GET  /metrics             … Prometheus metrics
GET  /v1/version          … list ReShifter status and version
POST /v1/backup           … start backup
GET  /v1/backup/$BACKUPID … download backup $BACKUPID
POST /v1/restore          … start restore
GET  /v1/explorer         … auto-discovery of etcd and Kubernetes
```


## Testbed

There are three types of integration/perf/scale tests in ReShifter, all available in [testbed](testbed/):

- end-to-end tests
- synthetic tests
- tests based on cluster dumps

### End-to-end tests

For end-to-end tests do the following. Note that each might take up to 30s and that you MUST execute them from within the `testbed/` directory:

```
$ cd testbed/
$ e2e-etcd2.sh
$ e2e-etcd3.sh
```

The end-to-end tests have the following dependencies:

- Docker CE (tested with v1.17)
- [etcdctl](https://github.com/coreos/etcd/tree/master/etcdctl)
- [http](https://httpie.org)
- [jq](https://stedolan.github.io/jq/)

The end-to-end test matrix is as follows:

|version   | insecure  | secure       |
| --------:| --------- | ------------ |
| 2.x      | available | available*   |
| 3.x      | available | available**  |

Legend:

- `*` … based on the etcd2 [security flags](https://coreos.com/etcd/docs/latest/v2/configuration.html#security-flags) and the etcd2 [security model](https://coreos.com/etcd/docs/latest/v2/security.html)
- `**` … based on the etcd3 [security flags](https://coreos.com/etcd/docs/latest/op-guide/configuration.html#security-flags) and the etcd3 [security model](https://coreos.com/etcd/docs/latest/op-guide/security.html)

See also the notes on [setting up etcd in a secure way](testbed/certs/README.md) for more details on how to change or extend these tests.

### Synthetic tests

For synth tests, execute `testbed/gen-synth-testbed.sh`, which creates a number of Kubernetes objects and requires access to a Kubernetes cluster.

### Cluster dumps

TBD.


## Extending it

To extend ReShifter or fix issues, please consider the following.

### Vendoring

We are using Go [dep](https://github.com/golang/dep) for dependency management.
If you don't have `dep` installed yet, do `go get -u github.com/golang/dep/cmd/dep` now and then:

```
$ dep ensure
```

### Unit tests

In general, for unit tests we use the `go test` command, for example:

```
$ cd pkg/backup/
$ go test -v
```

Please do make sure all unit tests pass before sending in a PR.
