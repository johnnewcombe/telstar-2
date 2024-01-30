build:
	env GOOS=darwin GOARCH=amd64  go build -o /usr/local/bin/telstar-emf .
	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-emf-2.0/telstar-emf-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-emf-2.0/telstar-emf-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-emf-2.0/telstar-emf-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-emf-2.0/telstar-emf-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-emf-2.0/telstar-emf-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-emf-2.0/telstar-emf-linux-arm64 .

	zip -v ./telstar-emf-2.0/telstar-emf-2.0.zip  ./telstar-emf-2.0/telstar-emf-macos-amd64 ./telstar-emf-2.0/telstar-emf-macos-arm64 ./telstar-emf-2.0/telstar-emf-windows-386.exe ./telstar-emf-2.0/telstar-emf-windows-amd64.exe ./telstar-emf-2.0/telstar-emf-linux-amd64 ./telstar-emf-2.0/telstar-emf-linux-arm64

	# update the deployment folder with the amd64 version
	cp ./telstar-emf-2.0/telstar-emf-linux-amd64 ../../telstar-2-system/telstar-deployment/opt/telstar/telstar-emf
