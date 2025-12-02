GOLANGCI_LINT_VERSION?=2.6.2
GOLANGCI_LINT_SHA256?=499c864b5fd9841c4fa8e80b5e2be30f73f085cf186f1b111ff81a2783b7de12
GOLANGCI_LINT=/usr/local/bin/golangci-lint

$(GOLANGCI_LINT):
	curl -sSLO https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz
	shasum -a 256 golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz | grep "^${GOLANGCI_LINT_SHA256}  " > /dev/null
	tar -xf golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz
	sudo mv golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64/golangci-lint /usr/local/bin/golangci-lint
	rm -rf golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64*

.PHONY: test
test:
	@echo "==> Running tests"
	go test -v

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run
