tools: tools/lint

tools/lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2
	go install github.com/segmentio/golines@v0.13.0
	go install github.com/daixiang0/gci@v0.13.7
	go install github.com/securego/gosec/v2/cmd/gosec@v2.22.9

lint: tools/lint
	golangci-lint run

lint/fix: tools/lint
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default .

test:
	go test ./...

test/v:
	go test ./... -v 

test/race:
	go test ./... -race

test/cov/report:
	go test ./... -coverpkg=./... -covermode=atomic -coverprofile=coverage.out -race

test/cov: test/cov/report
	go tool cover -html=coverage.out -o=coverage.html
	make test/cov/open

test/cov/open:
	open ./coverage.html
