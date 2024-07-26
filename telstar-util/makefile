build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-util .

	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-util-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-util-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-util-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-util-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-util-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-util-linux-arm64 .

	zip -v ./bin/telstar-util.zip ./bin/telstar-util-macos-amd64 ./bin/telstar-util-macos-arm64 ./bin/telstar-util-windows-386.exe ./bin/telstar-util-windows-amd64.exe ./bin/telstar-util-linux-amd64 ./bin/telstar-util-linux-arm64

