build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-upload .

	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-upload-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-upload-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-upload-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-upload-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-upload-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-upload-linux-arm64 .

	zip -v ./bin/telstar-upload.zip ./bin/telstar-upload-macos-amd64 ./bin/telstar-upload-macos-arm64 ./bin/telstar-upload-windows-386.exe ./bin/telstar-upload-windows-amd64.exe ./bin/telstar-upload-linux-amd64 ./bin/telstar-upload-linux-arm64

