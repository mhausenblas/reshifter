mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
project_name := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
app_name ?= $(project_name)-app
public_port ?=8080

.PHONY: up create init build publish destroy

up :  build

create : init build publish

init :
	@oc new-project $(project_name)
	@oc new-app --strategy=docker --name='$(app_name)' --context-dir='./app/' . --output yaml > app.yaml
	@oc apply -f app.yaml

build :
	@oc start-build $(app_name) --from-dir .
	@oc logs -f bc/$(app_name)

publish :
	@oc expose dc $(app_name) --port=$(public_port)
	@oc expose svc/$(app_name)

destroy :
	@rm app.yaml
	@oc delete project $(project_name)
