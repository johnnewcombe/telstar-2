build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-upload .

	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-upload-2.0/telstar-upload-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-upload-2.0/telstar-upload-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-upload-2.0/telstar-upload-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-upload-2.0/telstar-upload-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-upload-2.0/telstar-upload-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-upload-2.0/telstar-upload-linux-arm64 .

	zip -v ./telstar-upload-2.0/telstar-upload-2.0.zip ./telstar-upload-2.0/telstar-upload-macos-amd64 ./telstar-upload-2.0/telstar-upload-macos-arm64 ./telstar-upload-2.0/telstar-upload-windows-386.exe ./telstar-upload-2.0/telstar-upload-windows-amd64.exe ./telstar-upload-2.0/telstar-upload-linux-amd64 ./telstar-upload-2.0/telstar-upload-linux-arm64

	# update the deployment folder with the amd64 version
	cp ./telstar-upload-2.0/telstar-upload-linux-amd64 ../../telstar-2-system/telstar-deployment/opt/telstar/telstar-upload
