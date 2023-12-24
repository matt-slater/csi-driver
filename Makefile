.PHONY: version

version=$(shell cat VERSION)
commit_hash=$()
ldflags=""

version:
	@echo ${version}

build:
	