# function-extra-resources
[![CI](https://github.com/crossplane/function-template-go/actions/workflows/ci.yml/badge.svg)](https://github.com/crossplane/function-template-go/actions/workflows/ci.yml)

A function for selecting extra resources via [composition function][functions]s in [Go][go].

**just a fork of the function template ATM**

# Local dev.

Use air to iterate quickly.

run `~/go/bin/air -- --insecure --debug --address localhost:9443`.

To start this server process

Then run `./run.sh` to run a function-extra-resources example.

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --tag=runtime

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build -f package --embed-runtime-image=runtime
```

[functions]: https://docs.crossplane.io/latest/concepts/composition-functions
[go]: https://go.dev
[function guide]: https://docs.crossplane.io/knowledge-base/guides/write-a-composition-function-in-go
[package docs]: https://pkg.go.dev/github.com/crossplane/function-sdk-go
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli
