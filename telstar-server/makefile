build:
	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-server-2.0/telstar-server-macos-amd64 telstar.go
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-server-2.0/telstar-server-macos-arm64 telstar.go
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-server-2.0/telstar-server-windows-amd64.exe telstar.go
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-server-2.0/telstar-server-linux-amd64 telstar.go
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-server-2.0/telstar-server-linux-arm64 telstar.go

	zip -v ./telstar-server-2.0/telstar-server-2.0.zip  ./telstar-server-2.0/telstar-server-macos-amd64 ./telstar-server-2.0/telstar-server-macos-arm64 ./telstar-server-2.0/telstar-server-windows-amd64.exe ./telstar-server-2.0/telstar-server-linux-amd64 ./telstar-server-2.0/telstar-server-linux-arm64 ./telstar-server-2.0/telstar-server-linux-arm
