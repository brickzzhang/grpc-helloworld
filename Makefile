DIRS=$(shell ls -1 -F | grep "/$$" | grep -v vendor)
GLINT=$(shell command -v golangci-lint 2>/dev/null)

###############################################################################
#                            Protocol Buffer targets                          #
###############################################################################
.PHONY: pb
pb:
# NOTICE: Needed to be modified (such as api/services), otherwise it will use `api` path by default
	@bash ./build/pb_generate.sh -e project -p api

.PHONY: lc-pb
lc-pb:
# NOTICE: Needed to be modified (such as api/services), otherwise it will use `api` path by default
	@bash ./build/pb_generate.sh -e local -p api

###############################################################################
#                            Formation targets                                #
###############################################################################

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
git-hooks:
	@echo "Copy git hooks"
	@find .git/hooks -type l -exec rm {} \;
	@find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

###############################################################################
#                            CICD targets                                     #
###############################################################################

.PHONY: build
build: pb
	@go mod vendor
	@go build -o boot.bin internal/main.go

.PHONY: lc-build
lc-build: lc-pb
	@go mod vendor
	@go build -o boot.bin internal/main.go

.PHONY: run
run:
	./boot.bin -f configs/config.yaml

###############################################################################
#                            client targets                                   #
###############################################################################

.PHONY: cli-build
cli-build: pb
	@cd client && go build -o client.bin -mod=vendor ./

.PHONY: lc-cli-build
lc-cli-build: lc-pb
	@cd client && go build -o client.bin -mod=vendor ./

.PHONY: cli-run
cli-run:
	client/client.bin -f configs/config.yaml
