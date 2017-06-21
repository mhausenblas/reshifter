# ReShifter

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

## Vendoring

Using Go [dep](https://github.com/golang/dep) for dependency management.

### Testing

In general:

```
$ go test -v
```

See [test/](test/) directory for `etc-*` scripts.
