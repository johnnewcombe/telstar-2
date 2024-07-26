
build:
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/local/bin/telstar-ftse main.go
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-ftse-macos-amd64 main.go
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-ftse-macos-arm64 main.go
	env GOOS=windows GOARCH=amd64  go build -o ./bin/telstar-ftse-windows-amd64.exe main.go
	env GOOS=linux GOARCH=amd64  go build -o ./bin/telstar-ftse-linux-amd64 main.go
	env GOOS=linux GOARCH=arm64  go build -o ./bin/telstar-ftse-linux-arm64 main.go

	zip -v ./bin/telstar-ftse.zip  ./bin/telstar-ftse-macos-amd64 ./bin/telstar-ftse-macos-arm64 ./bin/telstar-ftse-windows-amd64.exe ./bin/telstar-ftse-linux-amd64 ./bin/telstar-ftse-linux-arm64

