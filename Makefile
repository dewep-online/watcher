SHELL=/bin/bash

.PHONY: run_back
run_back:
	go run -race cmd/watcher/main.go run "echo 'test'"

.PHONY: build_back build_font
build_back:
	bash scripts/build.sh amd64

.PHONY: linter
linter:
	bash scripts/linter.sh

.PHONY: tests
tests:
	bash scripts/tests.sh

.PHONY: ci
ci:
	bash scripts/ci.sh

.PHONY: ci
deb: 
	deb-builder build