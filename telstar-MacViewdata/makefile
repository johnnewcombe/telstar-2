build:
	sudo env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-obbs .

	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-macviewdata-2.0/telstar-macviewdata-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-macviewdata-2.0/telstar-macviewdata-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-macviewdata-2.0/telstar-macviewdata-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-macviewdata-2.0/telstar-macviewdata-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-macviewdata-2.0/telstar-macviewdata-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-macviewdata-2.0/telstar-macviewdata-linux-arm64 .

	zip -v ./telstar-macviewdata-2.0/telstar-macviewdata-2.0.zip  ./telstar-macviewdata-2.0/telstar-macviewdata-macos-amd64 ./telstar-macviewdata-2.0/telstar-macviewdata-macos-arm64 ./telstar-macviewdata-2.0/telstar-macviewdata-windows-amd64.exe ./telstar-macviewdata-2.0/telstar-macviewdata-linux-amd64 ./telstar-macviewdata-2.0/telstar-macviewdata-linux-arm64
