
# Check richgo does exist.
ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go get github.com/kyoh86/richgo")
endif

.PHONY: test build sync codecov test-app

default: all

all: test build

test: sync
	$(info _______________________running tests_______________________)
	 richgo test -v ./...

codecov: sync
	$(info __________________running tests coverage___________________)
	sh build/script/coverage.sh

sync: 
	$(info _________________downloading dependencies___________________)
	go get -v ./...
