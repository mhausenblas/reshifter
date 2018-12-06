reshifter_version := 0.3.22
git_version := `git rev-parse HEAD`
app_name := reshifter-app
main_dir := `pwd`

.PHONY: gtest gbuild gbuildcli gbuildapp gclean cbuild cpush release init build publish destroy

###############################################################################
# commands related to Go testing and builds --> binaries (HTTP API and CLI tool)
gtest :
	@echo Testing the library. This will take ca. 3 min to complete so get a cuppa tea for now ...
	go test -short -run Test* ./pkg/discovery
	go test -short -run Test* ./pkg/backup
	go test -short -run Test* ./pkg/restore

gbuild : gbuildcli gbuildapp

gbuildcli :
	go build -ldflags "-X github.com/mhausenblas/reshifter/rcli/cmd.releaseVersion=$(reshifter_version)" -o ./rcli-macos rcli/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/mhausenblas/reshifter/rcli/cmd.releaseVersion=$(reshifter_version)" -o ./rcli-linux rcli/main.go

gbuildapp :
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/mhausenblas/reshifter/app/handler.releaseVersion=$(reshifter_version)" -o ./reshifter app/main.go

gclean :
	@rm reshifter
	@rm rcli-macos
	@rm rcli-linux


###############################################################################
# commands related to container image builds
crelease : cbuild cpush

cbuild :
	@docker build --build-arg rversion=$(reshifter_version) -t quay.io/mhausenblas/reshifter:$(reshifter_version) .

cpush :
	@docker push quay.io/mhausenblas/reshifter:$(reshifter_version)


###############################################################################
# commands related to OpenShift-specific build and deployment

init :
	@oc new-project reshifter
	@oc new-app --strategy=docker --name='$(app_name)' --context-dir='./app/' .
	@oc set env bc/$(app_name) rversion=$(reshifter_version)
	# @oc apply -f app.yaml

build :
	@oc start-build $(app_name) --from-dir .
	@oc logs -f bc/$(app_name)

publish :
	@oc expose svc/$(app_name)

destroy :
	@oc delete project reshifter
