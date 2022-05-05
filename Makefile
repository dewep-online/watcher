SHELL=/bin/bash

.PHONY: run
run:
	go run -race cmd/watcher/main.go run ping 1.1.1.1

.PHONY: build build_font
build:
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

.PHONY: deb
deb: 
	deb-builder build