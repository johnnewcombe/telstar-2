build:
	#env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-telesoftware .

	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-telesoftware-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-telesoftware-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-telesoftware-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-telesoftware-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-telesoftware-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-telesoftware-linux-arm64 .

	zip -v ./bin/telstar-telesoftware.zip  ./bin/telstar-telesoftware-macos-amd64 ./bin/telstar-telesoftware-macos-arm64 ./bin/telstar-telesoftware-windows-386.exe ./bin/telstar-telesoftware-windows-amd64.exe ./bin/telstar-telesoftware-linux-amd64 ./bin/telstar-telesoftware-linux-arm64

