reshifter_version := 0.2.4
git_version := `git rev-parse HEAD`
app_name := reshifter-app

.PHONY: gbuild cbuild cpush registry init build publish destroy

gbuild :
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.releaseVersion=$(reshifter_version)" .

cbuild :
	@docker build -t quay.io/mhausenblas/reshifter:$(reshifter_version) .
	@rm reshifter

cpush :
	@docker push quay.io/mhausenblas/reshifter:$(reshifter_version)

registry : gbuild cbuild cpush

init :
	@oc new-project reshifter
	@oc new-app --strategy=docker --name='$(app_name)' . --output yaml > app.yaml
	@oc apply -f app.yaml

build :
	@oc start-build $(app_name) --from-dir .
	@oc logs -f bc/$(app_name)

publish :
	@oc expose svc/$(app_name)

destroy :
	@rm app.yaml
	@oc delete project reshifter
