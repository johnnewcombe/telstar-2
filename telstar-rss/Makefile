build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-rss .
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-rss-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-rss-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-rss-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-rss-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-rss-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-rss-linux-arm64 .

	zip -v ./bin/telstar-rss.zip  ./bin/telstar-rss-macos-amd64 ./bin/telstar-rss-macos-arm64 ./bin/telstar-rss-windows-386.exe ./bin/telstar-rss-windows-amd64.exe ./bin/telstar-rss-linux-amd64 ./bin/telstar-rss-linux-arm64

