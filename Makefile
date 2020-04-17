.PHONY: build test

build: test
	go build ./...

test:
	go test -v -race ./...

%.zip: %
	cd ./pkg && rm -rf postqueue2json-$@ && zip -rD postqueue2json-$@ $<

linux: test
	rm -rf pkg/linux
	mkdir -p pkg/linux
	GOOS=linux go build -o pkg/linux/postqueue2json .
	cd ./pkg && zip -rD postqueue2json-linux.zip linux

darwin: test
	rm -rf pkgdarwin
	mkdir -p pkg/darwin
	GOOS=darwin go build -o pkg/darwin/postqueue2json .

windows: test
	rm -rf pkg/windows
	mkdir -p pkg/windows
	GOOS=windows go build -o pkg/windows/postqueue2json.exe .

release: darwin.zip linux.zip windows.zip
