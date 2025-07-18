.PHONY: build runexe run

build:
	echo "Size before build:"; ls -la example |grep example; ls -lh example |grep example; echo "\n\nSize after build:"; go build --ldflags "-s -w" -o example/example ./example; ls -la example |grep example; ls -lh example |grep example

runexe:
	./example/example

run:
	go run ./example/main.go
