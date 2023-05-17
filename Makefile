build-osx:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/hlsdl_osx ./cmd/hlsdl

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/hlsdl_linux ./cmd/hlsdl

build-windows:
	env GOOS=windows GOARCH=amd64 go build -o ./bin/hlsdl_windows.exe ./cmd/hlsdl

build: build-osx build-linux build-windows