# ReShifter

A cluster admin backup and restore tool for OpenShift:

- etcd is handled via [burry](http://burry.sh)
- projects are handled via `oc export dc`
- other cluster objects like namespaces are handled via [direct API access](https://github.com/kubernetes/client-go)

## Development

Using OpenShift and the here provided [Makefile](Makefile):

```
# get a local copy of this repo:
$ git clone https://github.com/mhausenblas/reshifter.git && cd $_

# create the project and app in OpenShift, build and deploy it:
$ make create

# whenever something changes locally (builds and updates):
$ make up
```

For a smoother workflow, that is, if you want changes automatically trigger builds and updates of the app, use the here provided script [outline.sh](outline.sh) (which I've aliased like so: `alias outline='./outline.sh'`):

```
$ outline create

$ outline up

CTRL-Z -> to background, now you can edit your code and when you save it, the app is updated
```

### Use a local etcd for testing

```
docker run -d -p 2379:2379 -p 2380:2380 -p 4001:4001 -p 7001:7001 -v $(pwd):/data --name test-etcd elcolio/etcd:2.0.10
```
