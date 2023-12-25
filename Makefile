.PHONY: version build

project=csi-driver
version=$(shell cat VERSION)
commitHash=$(shell git rev-parse --short HEAD)
ldFlags=-ldflags="-X 'main.version=${version}' -X 'main.commit=${commitHash}'"

imageName=mattslater/${project}-linux:${version}

version:
	@echo ${version}

build:
	GOOS=darwin GOARCH=arm64 go build ${ldFlags} -o build/package/${project}/${project}-darwin cmd/${project}/${project}.go
	GOOS=linux GOARCH=amd64 go build ${ldFlags} -o build/package/${project}/${project}-linux cmd/${project}/${project}.go
	docker build --platform linux/amd64 build/package/${project} -t ${imageName} 