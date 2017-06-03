# ReShifter

A cluster admin backup and restore tool for OpenShift:

- etcd is handled via [burry](http://burry.sh)
- projects are handled via `oc export dc`
- other cluster objects like namespaces are handled via [direct API access](https://github.com/kubernetes/client-go)

## Development

Using OpenShift:

```
git clone https://github.com/mhausenblas/reshifter.git

oc new-project reshifter
oc new-build --strategy=docker --name='rs-app' --context-dir='./app/' .
oc start-build rs-app --from-dir .
oc logs -f bc/rs-app

oc run reshifter --image=$REGISTRY_IP:5000/reshifter/rs-app
oc expose dc reshifter --port=8080
oc expose svc/reshifter
http http://reshifter-reshifter.example.com/v1/backup
```

Local test etcd:

```
docker run -d -p 2379:2379 -p 2380:2380 -p 4001:4001 -p 7001:7001 -v $(pwd):/data --name test-etcd elcolio/etcd:2.0.10
```
