SUBDIRS := $(wildcard ./cmd/*/.)

default: help

build: $(SUBDIRS)  ## Builds the application 
img: $(SUBDIRS)    ## Build Docker Image

test:  ## Run all tests
	@go clean --testcache && go run gotest.tools/gotestsum@latest --junitfile results.xml

cover:  ## Run test coverage suite
	@mkdir -p cover_results
	@go test ./... --coverprofile=cover_results/cov.out
	@go tool cover --html=cover_results/cov.out -o cover_results/cov.html
	
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'

.PHONY: $(SUBTARGETS) $(SUBDIRS)