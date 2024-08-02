build:
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-server-macos-amd64 telstar.go
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-server-macos-arm64 telstar.go
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-server-windows-amd64.exe telstar.go
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-server-linux-amd64 telstar.go
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-server-linux-arm64 telstar.go

	# Removed as MongoDB needs 64 bit
	#env GOOS=linux GOARCH=arm   go build -o ./bin/telstar-server-linux-arm telstar.go

	zip -v ./bin/telstar-server.zip  ./bin/telstar-server-macos-amd64 ./bin/telstar-server-macos-arm64 ./bin/telstar-server-windows-amd64.exe ./bin/telstar-server-linux-amd64 ./bin/telstar-server-linux-arm64
