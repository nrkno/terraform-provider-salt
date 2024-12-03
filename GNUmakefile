.PHONY: build
build: fmt
	go install

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./...

.PHONY: fmt
fmt:
	golangci-lint run --fix

.PHONY: check
check:
	pre-commit run -a
