CMD = go build -ldflags "-s -w"

compile:
	@test -n "$(version)" || (echo 'Error: missing arg version=<version>' && exit 1)
	GOOS=darwin GOARCH=arm64 $(CMD) -o bin/cl-tool-macos-arm64-$(version)
	GOOS=darwin GOARCH=amd64 $(CMD) -o bin/cl-tool-macos-amd64-$(version)
	GOOS=linux GOARCH=amd64 $(CMD) -o bin/cl-tool-linux-amd64-$(version)
	GOOS=windows GOARCH=amd64 $(CMD) -o bin/cl-tool-windows-amd64-$(version)

clean:
	rm -rf bin cl-tool
