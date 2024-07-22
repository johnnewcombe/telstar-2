build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-util .

	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-util-2.5/telstar-util-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-util-2.5/telstar-util-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-util-2.5/telstar-util-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-util-2.5/telstar-util-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-util-2.5/telstar-util-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-util-2.5/telstar-util-linux-arm64 .

	zip -v ./telstar-util-2.5/telstar-util-2.0.zip ./telstar-util-2.5/telstar-util-macos-amd64 ./telstar-util-2.5/telstar-util-macos-arm64 ./telstar-util-2.5/telstar-util-windows-386.exe ./telstar-util-2.5/telstar-util-windows-amd64.exe ./telstar-util-2.5/telstar-util-linux-amd64 ./telstar-util-2.5/telstar-util-linux-arm64

	# update the deployment folder with the amd64 version
	cp ./telstar-util-2.5/telstar-util-linux-amd64 ../../telstar-2-system/telstar-deployment/opt/telstar/telstar-util
