
# watcher

[![Coverage Status](https://coveralls.io/repos/github/dewep-online/watcher/badge.svg?branch=master)](https://coveralls.io/github/dewep-online/watcher?branch=master)
[![Release](https://img.shields.io/github/release/dewep-online/watcher.svg?style=flat-square)](https://github.com/dewep-online/watcher/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/dewep-online/watcher)](https://goreportcard.com/report/github.com/dewep-online/watcher)
[![CI](https://github.com/dewep-online/watcher/actions/workflows/ci.yml/badge.svg)](https://github.com/dewep-online/watcher/actions/workflows/ci.yml)

# install

```go
 go get -u github.com/dewep-online/watcher/cmd/... 
```

# how use

with default interval = 10s

```bash
watcher run <command>
```

with custom interval = 30s

```bash
watcher run --interval=30 <command>
```

# example

```bash
watcher run ping 1.1.1.1
```