build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-emf .
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-emf-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-emf-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-emf-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-emf-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-emf-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-emf-linux-arm64 .

	zip -v ./bin/telstar-emf.zip  ./bin/telstar-emf-macos-amd64 ./bin/telstar-emf-macos-arm64 ./bin/telstar-emf-windows-386.exe ./bin/telstar-emf-windows-amd64.exe ./bin/telstar-emf-linux-amd64 ./bin/telstar-emf-linux-arm64
