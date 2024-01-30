build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-telesoftware .

	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-telesoftware-2.0/telstar-telesoftware-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-telesoftware-2.0/telstar-telesoftware-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-telesoftware-2.0/telstar-telesoftware-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-telesoftware-2.0/telstar-telesoftware-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-telesoftware-2.0/telstar-telesoftware-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-telesoftware-2.0/telstar-telesoftware-linux-arm64 .

	zip -v ./telstar-telesoftware-2.0/telstar-telesoftware-2.0.zip  ./telstar-telesoftware-2.0/telstar-telesoftware-macos-amd64 ./telstar-telesoftware-2.0/telstar-telesoftware-macos-arm64 ./telstar-telesoftware-2.0/telstar-telesoftware-windows-386.exe ./telstar-telesoftware-2.0/telstar-telesoftware-windows-amd64.exe ./telstar-telesoftware-2.0/telstar-telesoftware-linux-amd64 ./telstar-telesoftware-2.0/telstar-telesoftware-linux-arm64

