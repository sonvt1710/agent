.ONESHELL:
.DELETE_ON_ERROR:
export SHELL     := bash
export SHELLOPTS := pipefail:errexit
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rule

include docs.mk

.PHONY: docs
docs: check-cloudwatch-integration

check-cloudwatch-integration:
	$(PODMAN) run -v $(shell git rev-parse --show-toplevel):/repo -v $(shell pwd):/docs -w /repo golang:1.24.2-bullseye go run static/integrations/cloudwatch_exporter/docs/doc.go check /docs/sources/static/configuration/integrations/cloudwatch-exporter-config.md
	$(PODMAN) run -v $(shell git rev-parse --show-toplevel):/repo -v $(shell pwd):/docs -w /repo golang:1.24.2-bullseye go run static/integrations/cloudwatch_exporter/docs/doc.go check /docs/sources/flow/reference/components/prometheus.exporter.cloudwatch.md

generate-cloudwatch-integration:
	$(PODMAN) run -v $(shell git rev-parse --show-toplevel):/repo -v $(shell pwd):/docs -w /repo golang:1.24.2-bullseye go run static/integrations/cloudwatch_exporter/docs/doc.go generate

sources/assets/hierarchy.svg: sources/operator/hierarchy.dot
	cat $< | $(PODMAN) run --rm -i nshine/dot dot -Tsvg > $@
