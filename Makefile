.PHONY: version build lint fmt

project=csi-driver
version=$(shell cat VERSION)
commitHash=$(shell git rev-parse --short HEAD)
ldFlags=-ldflags="-X 'main.version=${version}' -X 'main.commit=${commitHash}'"
trimpath=-trimpath
buildFlags=${trimpath} ${ldFlags}
imageName=mattslater/${project}-linux:${version}

version:
	@echo ${version}

lint:
	golangci-lint run ./...

fmt:
	gofumpt -l -w .

build:
	GOOS=darwin GOARCH=arm64 go build ${buildFlags} -o build/package/${project}/${project}-darwin cmd/${project}/${project}.go
	GOOS=linux GOARCH=amd64 go build ${buildFlags} -o build/package/${project}/${project}-linux cmd/${project}/${project}.go
	docker build --platform linux/amd64 build/package/${project} -t ${imageName}

deploy:
	kind load docker-image ${imageName}
	kubectl apply -f deployments/

deploy_example:
	kubectl apply -f examples/