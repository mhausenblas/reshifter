reshifter_version := 0.3.5
git_version := `git rev-parse HEAD`
app_name := reshifter-app
main_dir := `pwd`

.PHONY: gbuild cbuild cpush release init build publish destroy

gbuild :
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.releaseVersion=$(reshifter_version)" .
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/mhausenblas/reshifter/rcli/cmd.releaseVersion=$(reshifter_version)" -o ./rcli/rcli-linux rcli/main.go
	@go build -ldflags "-X github.com/mhausenblas/reshifter/rcli/cmd.releaseVersion=$(reshifter_version)" -o ./rcli/rcli-macos rcli/main.go

ginstall :
	@GOOS=linux GOARCH=amd64 go install -ldflags "-X main.releaseVersion=$(reshifter_version)" .

cbuild :
	@docker build --build-arg rversion=$(reshifter_version) -t quay.io/mhausenblas/reshifter:$(reshifter_version) .

cpush :
	@docker push quay.io/mhausenblas/reshifter:$(reshifter_version)

crelease : cbuild cpush

clean :
	@rm reshifter

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
