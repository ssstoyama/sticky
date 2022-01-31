.PHONY: test
test:
	go test -v -cover -race .

.PHONY: ci
ci:
	go test -race .