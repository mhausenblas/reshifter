# Testbed

There are three types of integration/perf/scale tests in ReShifter, all available in [testbed](https://github.com/mhausenblas/reshifter/tree/master/testbed/):

- end-to-end tests
- synthetic tests
- tests based on cluster dumps

## End-to-end tests

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

See also the notes on [setting up etcd in a secure way](https://github.com/mhausenblas/reshifter/tree/master/testbed/certs/README.md) for more details on how to change or extend these tests.

## Synthetic tests

For synth tests, execute `testbed/gen-synth-testbed.sh`, which creates a number of Kubernetes objects and requires access to a Kubernetes cluster.

## Cluster dumps

TBD.
