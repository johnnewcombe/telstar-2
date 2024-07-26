build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-obbs .
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-macviewdata-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-macviewdata-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-macviewdata-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-macviewdata-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-macviewdata-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-macviewdata-linux-arm64 .

	zip -v ./bin/telstar-macviewdata.zip  ./bin/telstar-macviewdata-macos-amd64 ./bin/telstar-macviewdata-macos-arm64 ./bin/telstar-macviewdata-windows-amd64.exe ./bin/telstar-macviewdata-linux-amd64 ./bin/telstar-macviewdata-linux-arm64
