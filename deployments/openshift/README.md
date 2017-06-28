# Deploying ReShifter on OpenShift

```
oc new-app reshifter-template.yaml -p SOURCE_REPOSITORY_URL=https://github.com/mhausenblas/reshifter
oc get builds -w
oc logs -f bc/reshifter
oc get all
```
