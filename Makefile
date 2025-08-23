.PHONY: build runexe run test testv

build:
	echo "Size before build:"; ls -la example |grep example; ls -lh example |grep example; echo "\n\nSize after build:"; go build --ldflags "-s -w" -o example/example ./example; ls -la example |grep example; ls -lh example |grep example

check-build-files:
	go list -f '{{.GoFiles}}' . ./tui ./keys

runexe:
	./example/example

run:
	go run ./example/main.go

test:
	go test ./...

testv:
	go test -v ./...
