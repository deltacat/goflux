include .env

# Go related variables.
PKGS		:= $(shell go list ./...)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# get revision (git hash)
REVER := $(shell git describe --always |sed -e "s/^v//")

## help: this help
.PHONY: all help
all: help
help: Makefile
	@echo
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## req: Install development requirements
.PHONY: req
req:
	@echo "Installing requirements ..."
	go get golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint

## test: 执行代码检查和单元测试
.PHONY: test
test:
	@echo running tests
	@rm -f coverage.out
	@go test -p 1 -v -cover $(PKGS) -coverprofile coverage.out
	@echo done

## lint: 执行代码检查
.PHONY: lint
lint:
	@echo running code inspection
	@echo $(PKGS)
	@for pkg in $(PKGS) ; do \
		golint $$pkg ; \
	done
	@go vet $(PKGS)
	@echo done