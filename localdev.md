# Setup

- `pre-commit install` to install hooks
- Install the Pip3 package globally if you have not already `xmlformatter` i.e. `pip3 install xmlformatter --break-system-packages`
- Install Go v1.24.1. https://go.dev/doc/install
- Ensure that your `~/go/bin` is on `$PATH`
- Install goimports: `go install golang.org/x/tools/cmd/goimports@latest`
- Install golangci-lint: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.0.2`
- Install gofumpt: `go install mvdan.cc/gofumpt@latest`
- Install dlv for local debugging: `go install github.com/go-delve/delve/cmd/dlv@latest`, although personally I wouldn't bother, as the current state of VSCode+Delve is not great (as of 2025-04-05)
- Run this to exec all the pre-commit hooks on the entire project: `pre-commit run --all-files`
- Copy `.env.example` to `.env` and configure it with your settings
