

.PHONY: test
test: 
	go test ./...

.PHONY: cleantest
cleantest:
	go clean -testcache

.PHONY: example
example:
	go run ./_examples/up-and-down

.PHONY: acme-like
acme-like:
	go run ./_examples/acme-like

.PHONY: list-records
list-records:
	go run ./_examples/list-records

