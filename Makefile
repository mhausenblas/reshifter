git_version := `git rev-parse HEAD`
app_name := reshifter-app

.PHONY: gbuild hub init build publish destroy

gbuild :
	@GOOS=linux GOARCH=amd64 go build .

hub :
	@docker build -t mhausenblas/reshifter:$(git_version) .
	@docker push quay.io/mhausenblas/reshifter:$(git_version)
	# rm reshifter

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
