# reshifter
A cluster admin backup and restore tool for OpenShift



## Development

Using OpenShift:

```
git clone https://github.com/mhausenblas/reshifter.git

oc new-project reshifter
oc new-build --strategy=docker --name='rs-app' --context-dir='./app/'
oc start-build rs-app --from-dir .
oc logs -f bc/rs-app

oc run reshifter --image=$REGISTRY_IP:5000/reshifter/rs-app
oc expose dc reshifter --port=8080
oc expose svc/reshifter
http http://rs-app-reshifter.example.com/v1/backup
```
