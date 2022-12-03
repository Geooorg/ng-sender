GOPATH=$(shell go env GOPATH)
FINDFILES=find . \( -path ./proto -o -path ./.git -o -path ./vendor \) -prune -o -type f
XARGS = xargs -0 -r

test:
	go test -count 1 -vet "" -race ./...

clean:
	@go clean

dep:
	go mod tidy -compat=1.19; go mod vendor

todo:
	grep -rnw '.' --exclude-dir=vendor --exclude=makefile -e '// TODO:'

build: clean test
	go build -ldflags="-s -w" .

run:
	go run . serve-http --config config/application.yaml

tidy-go:
	@go mod tidy -compat=1.18

lint-go:
	@golangci-lint run -v -c .golangci.yml

format-go:
	@golangci-lint run --fix -c .golangci-format.yml

image:
	docker build -t ng-receiver:1.0.0 .

cross-compile:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w"  -o service .

build-and-docker: cross-compile image
