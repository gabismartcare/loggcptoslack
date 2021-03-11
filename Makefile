fmt:
	@go fmt $(shell go list ./... | grep -v /vendor/)

test:
	@go vet $(shell go list ./... | grep -v /vendor/)
	@go test -short $(shell go list ./... | grep -v /vendor/)

dep:
	@go mod download

package:
	@go build ./...

install:
	@go install ./...

feature-start :
	@read -p "Enter feature name: " module; \
	git flow feature start $$module

hotfix-start :
	@read -p "Enter version: " module; \
	git flow hotfix start $$module

release-start :
	@read -p "Enter version: " module; \
	git flow release start $$module

support-start :
	@read -p "Enter version: " module; \
	git flow support start $$module

feature-finish : fmt
	git flow feature finish -S

hotfix-finish : fmt
	git flow hotfix finish -S

release-finish : fmt
	git flow release finish -S
