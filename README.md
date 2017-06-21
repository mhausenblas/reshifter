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

### Monitoring

Prometheus metrics are available via `/metrics`.

## Extending it

To extend ReShifter or fixing issues, please consider the following.

### Vendoring

We are using Go [dep](https://github.com/golang/dep) for dependency management.

### Testing

In general:

```
$ go test -v
```

For end-to-end tests, Docker is required to provide a local etcd instance.
