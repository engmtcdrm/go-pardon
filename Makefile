.PHONY: build runexe run

build:
	go build --ldflags "-s -w" -o example/example ./example; ls -la example |grep example; ls -lh example |grep example

runexe:
	./example/example

run:
	go run ./example/main.go
