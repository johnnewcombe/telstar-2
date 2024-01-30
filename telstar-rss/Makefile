build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-rss .
	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-rss-2.0/telstar-rss-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-rss-2.0/telstar-rss-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-rss-2.0/telstar-rss-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-rss-2.0/telstar-rss-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-rss-2.0/telstar-rss-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-rss-2.0/telstar-rss-linux-arm64 .

	zip -v ./telstar-rss-2.0/telstar-rss-2.0.zip  ./telstar-rss-2.0/telstar-rss-macos-amd64 ./telstar-rss-2.0/telstar-rss-macos-arm64 ./telstar-rss-2.0/telstar-rss-windows-386.exe ./telstar-rss-2.0/telstar-rss-windows-amd64.exe ./telstar-rss-2.0/telstar-rss-linux-amd64 ./telstar-rss-2.0/telstar-rss-linux-arm64

	# update the deployment folder with the amd64 version
	cp ./telstar-rss-2.0/telstar-rss-linux-amd64 ../../telstar-2-system/telstar-deployment/opt/telstar/telstar-rss
	cp -R ./data/templates/* ../../telstar-2-system/telstar-deployment/opt/telstar/cron/telstar/data/templates/

