compile:
	echo "Compiling for every OS and Platform"
	GOOS=darwin GOARCH=arm64 go build -o bin/cl-tool-macos-arm64-$(version)
	GOOS=darwin GOARCH=amd64 go build -o bin/cl-tool-macos-amd64-$(version)
	GOOS=linux GOARCH=amd64 go build -o bin/cl-tool-linux-amd64-$(version)
	GOOS=windows GOARCH=amd64 go build -o bin/cl-tool-windows-amd64-$(version)
