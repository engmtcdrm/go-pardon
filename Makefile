.PHONY: build runexe run test testv

build:
	@
	cd examples; \
	echo "Size before build:"; \
	ls -la examples |grep examples; \
	ls -lh examples |grep examples; \
	echo "\n\nSize after build:"; \
	go build --ldflags "-s -w" -o examples; \
	ls -la examples |grep examples; \
	ls -lh examples |grep examples

check-build-files:
	@go list -f '{{.GoFiles}}' . ./tui ./keys

runexe:
	@./examples/examples

run:
	@cd examples; \
	go run .; \
	cd ..

test:
	@go test ./...

testv:
	@go test -v ./...

testcover:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
