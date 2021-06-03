DIRS=$(shell ls -1 -F | grep "/$$" | grep -v vendor)
GLINT=$(shell command -v golangci-lint 2>/dev/null)

###############################################################################
#                            Protocol Buffer targets                          #
###############################################################################
.PHONY: pb
pb:
	@bash ./build/pb_generate.sh -e project

.PHONY: lint
lint:
ifdef GLINT
	@echo "Checking golangci-lint..."
	@golangci-lint run --timeout 10m
else
	@echo "golangci-lint not found, please intall it"
	@exit -1
endif

.PHONY: fmt
fmt:
	@echo "==> Fixing source code with gofmt..."
	@for dir in $(DIRS) ; do `goimports -w $$dir` ; done
	@for dir in $(DIRS) ; do `gofmt -s -w $$dir` ; done
	@echo "==> Fixing source code with gofmt done"

.PHONY: binary
binary:
	@echo "==> Make binary file..."
	@go build -o boot.bin internal/main.go
	@echo "==> Make binary file done"

.PHONY: tools
tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: git-hooks
git-hooks: tools
	@echo "Copy git hooks"
	@find .git/hooks -type l -exec rm {} \;
	@find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;