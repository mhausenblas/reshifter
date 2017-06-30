# Development

If you plan to fix bugs or contribute to ReShifter, please consider the following.

## Builds and releases

We're following [semantic versioning](http://semver.org/). The canonical ReShifter release version is defined in one place only,
in the [Makefile](https://github.com/mhausenblas/reshifter/blob/master/Makefile).

This version is then used in the Go code, in the Docker image as a tag and for all downstream deployments.

A new release (Linux binary on GitHub and image on quay.io) is cut using the following process:

```
# 1. Generate the binary:
$ make gbuild

# 2. Release on GitHub, using `v$reshifter_version`

# 3. Build a container image locally and push it to quay.io:
$ make crelease
```

## Vendoring

We are using Go [dep](https://github.com/golang/dep) for dependency management.
If you don't have `dep` installed yet, do `go get -u github.com/golang/dep/cmd/dep` now and then:

```
$ dep ensure
```

## Unit tests

In general, for unit tests we use the `go test` command, for example:

```
$ cd pkg/backup/
$ go test -v
```

Please do make sure all unit tests pass before sending in a PR.

## Local testing

The following shows an example (interactive) session against an etcd3-based Kubernetes control plane.

First, launch ReShifter:

```
$ docker run --rm -e "ACCESS_KEY_ID=Q3AM3UQ867SPQQA43P2F" -e "SECRET_ACCESS_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG" --name reshifter -p 8080:8080 quay.io/mhausenblas/reshifter:0.3.0
```

Now, launch a local etcd. Note: use the result of `docker inspect test-etcd | jq -r '.[0].NetworkSettings.IPAddress'`
as the value for the endpoint in the UI/API:

```
$ docker run --rm -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
```

Next we generate some entries in etcd:

```
$ export ETCDCTL_API=3
$ etcdctl --endpoints=http://127.0.0.1:2379 put /kubernetes.io "."
$ etcdctl --endpoints=http://127.0.0.1:2379 put /kubernetes.io/namespaces/kube-system "."
$ etcdctl --endpoints=http://127.0.0.1:2379 put /openshift.io "."
```

Now you can use the UI to create a backup and after restarting etcd3 you can restore it again.

Note that if you want to use etc2, do the following:

```
$ docker run --rm -d -p 2379:2379 --name test-etcd --dns 8.8.8.8 --env ETCD_DEBUG quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd  \
--advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380
$ curl http://127.0.0.1:2379/v2/keys/kubernetes.io/namespaces/kube-system -XPUT -d value="."
$ curl http://127.0.0.1:2379/v2/keys/openshift.io -XPUT -d value="."
```

## Demo

The demo given to the Kubernetes SIG Cluster Lifecycle on 2017-06-27:

```
# Use Minio playground as S3 compatible storage backend:
cd /Users/mhausenblas/Dropbox/dev/work/src/github.com/mhausenblas/reshifter
export ACCESS_KEY_ID=Q3AM3UQ867SPQQA43P2F
export SECRET_ACCESS_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG

# Launch etcd2:
docker run --rm -p 2379:2379 \
           --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 \
          --advertise-client-urls http://0.0.0.0:2379 \
          --listen-client-urls http://0.0.0.0:2379

# Launch ReShifter:
cd /Users/mhausenblas/Dropbox/dev/work/src/github.com/mhausenblas/reshifter
DEBUG=true reshifter

# Populate etcd:
curl http://localhost:2379/v2/keys/kubernetes.io/namespaces/kube-system -XPUT -d \
         value="{\"kind\":\"Namespace\",\"apiVersion\":\"v1\"}"

# Backup via UI:
# Open http://localhost:8080/reshifter/
# Config -> Backup
# Open https://play.minio.io:9000/ and verify with bucket

# Re-start etcd:
docker kill test-etcd
docker run --rm -p 2379:2379 \
           --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 \
          --advertise-client-urls http://0.0.0.0:2379 \
          --listen-client-urls http://0.0.0.0:2379

# Restore via UI:
# Open http://localhost:8080/reshifter/
# Restore

# Query etcd to verify restore:
http http://localhost:2379/v2/keys/kubernetes.io/namespaces/kube-system
```

## etcd key prefixes

### Kubernetes

```
/kubernetes.io/ranges
/kubernetes.io/statefulsets
/kubernetes.io/jobs
/kubernetes.io/horizontalpodautoscalers
/kubernetes.io/events
/kubernetes.io/masterleases
/kubernetes.io/minions
/kubernetes.io/persistentvolumes
/kubernetes.io/configmaps
/kubernetes.io/controllers
/kubernetes.io/deployments
/kubernetes.io/serviceaccounts
/kubernetes.io/services
/kubernetes.io/namespaces
/kubernetes.io/securitycontextconstraints
/kubernetes.io/thirdpartyresources
/kubernetes.io/persistentvolumeclaims
/kubernetes.io/pods
/kubernetes.io/replicasets
/kubernetes.io/secrets
```

### OpenShift

```
/openshift.io/authorization
/openshift.io/buildconfigs
/openshift.io/oauth
/openshift.io/registry
/openshift.io/users
/openshift.io/useridentities
/openshift.io/builds
/openshift.io/deploymentconfigs
/openshift.io/images
/openshift.io/imagestreams
/openshift.io/ranges
/openshift.io/routes
/openshift.io/templates
```
