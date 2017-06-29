# Development notes

## local development

```
# launch ReShifter:
docker run --rm -e "ACCESS_KEY_ID=Q3AM3UQ867SPQQA43P2F" -e "SECRET_ACCESS_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG" --name reshifter -p 8080:8080 quay.io/mhausenblas/reshifter:0.2.4

# launch test etcd, note: use the result of the last command as the endpoint in the UI/API:
docker run --rm -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
curl http://localhost:2379/v2/keys/kubernetes.io/namespaces/kube-system -XPUT -d value="."
docker inspect test-etcd | jq -r '.[0].NetworkSettings.IPAddress'
```

## demo

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
