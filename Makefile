APP_NAME = eks-shared-usage-calculator
REPO_NAME = ${APP_NAME}

GOOS = linux
GOARCH = amd64
BINARY_VERSION = $(shell awk -F\" '/version string = "/ {print $$2}' main.go)
BINARY_NAME = $(REPO_NAME)-$(GOOS)-$(GOARCH)-$(BINARY_VERSION)

VERSIONS = \
	export AWS_REGION_OVERRIDE=$(AWS_REGION_OVERRIDE) && \
	export BINARY_VERSION=$(BINARY_VERSION) && \
	export GO_VERSION=1.19 && \
	export AWS_CLI_VERSION=2.2.8 && \
	export PRETTIER_VERSION=sha256:2bee846df0094f583714a49a6962619e7d9e192efcb18003280fa29cdf51cd25

ifdef CI
	include .cicd/Makefile
else
	include .local/Makefile
endif
