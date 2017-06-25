# ReShifter

[![Go Report Card](https://goreportcard.com/badge/github.com/mhausenblas/reshifter)](https://goreportcard.com/report/github.com/mhausenblas/reshifter)
[![godoc](https://godoc.org/github.com/mhausenblas/reshifter?status.svg)](https://godoc.org/github.com/mhausenblas/reshifter)

A cluster admin backup and restore tool for OpenShift. Traverses etcd and stores the state of all objects.

Supported:

- Cluster: Kubernetes 1.5
- Client: Linux, macOS

## Using it

### Deployment

TBD

### API

```
GET  /reshifter           … the ReShifter Web UI
GET  /metrics             … Prometheus metrics
GET  /v1/version          … list ReShifter status and version
POST /v1/backup           … start backup
GET  /v1/backup/$BACKUPID … download backup $BACKUPID
POST /v1/restore          … start restore
```

### Status

- [x] support insecure and secure etcd v2
- [x] create backup, stored in local filesystem
- [x] restore backup from local filesystem
- [x] download backup
- [ ] offer a simple UI
- [ ] upload backup in S3
- [ ] restore backup from S3
- [ ] support insecure and secure etcd v3
- [ ] deployment


### Monitoring

Prometheus metrics are available via `/metrics`.

## Extending it

To extend ReShifter or fixing issues, please consider the following.

### Vendoring

We are using Go [dep](https://github.com/golang/dep) for dependency management.

### Testing

In general, for unit tests we use the `go test` command, for example

```
$ cd pkg/etcd/
$ go test -v
```

For end-to-end tests, run `./test-e2e-etcd*.sh`. Note that full end-to-end tests require the following:

- Docker CE (tested with v1.17)
- [http](https://httpie.org)
- [jq](https://stedolan.github.io/jq/)

The test matrix is as follows:

|version   | insecure  | secure       |
| --------:| --------- | ------------ |
| 2.x      | available | available*   |
| 3.x      | available | available**  |

Legend:

- `*` … based on the etcd2 [security flags](https://coreos.com/etcd/docs/latest/v2/configuration.html#security-flags) and the etcd2 [security model](https://coreos.com/etcd/docs/latest/v2/security.html)
- `**` … based on the etcd3 [security flags](https://coreos.com/etcd/docs/latest/op-guide/configuration.html#security-flags) and the etcd3 [security model](https://coreos.com/etcd/docs/latest/op-guide/security.html)

See also the notes on [setting up etcd in a secure way](certs/README.md).
