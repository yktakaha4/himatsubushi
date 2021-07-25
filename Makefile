#!/usr/bin/make -f

MAKEFLAGS+=--warn-undefined-variables
SHELL:=/bin/bash
.SHELLFLAGS:=-eu -o pipefail -c
.DEFAULT_GOAL:=help
.SILENT:

# all targets are phony
.PHONY: $(shell egrep -o ^[a-zA-Z_-]+: $(MAKEFILE_LIST) | sed 's/://')

taskA: ## executes task A
	echo 'Starting $@'
	sleep 1
	echo 'Finished $@'

taskB: ## executes task B
	echo 'Starting $@'
	sleep 1
	echo 'Finished $@'

taskC: taskB ## executes task C (depends task B)
	echo 'Starting $@'
	sleep 1
	echo 'Finished $@'

help: ## Print this help
	echo 'Usage: make [target]'
	echo ''
	echo 'Targets:'
	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
