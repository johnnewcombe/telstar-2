build:
	sudo env GOOS=darwin GOARCH=amd64  go build -o /opt/telstar/telstar-openweather .

	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-openweather-2.0/telstar-openweather-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-openweather-2.0/telstar-openweather-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./telstar-openweather-2.0/telstar-openweather-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./telstar-openweather-2.0/telstar-openweather-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./telstar-openweather-2.0/telstar-openweather-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./telstar-openweather-2.0/telstar-openweather-linux-arm64 .

	zip -v ./telstar-openweather-2.0/telstar-openweather-2.0.zip  ./telstar-openweather-2.0/telstar-openweather-macos-amd64 ./telstar-openweather-2.0/telstar-openweather-macos-arm64 ./telstar-openweather-2.0/telstar-openweather-windows-amd64.exe ./telstar-openweather-2.0/telstar-openweather-linux-amd64 ./telstar-openweather-2.0/telstar-openweather-linux-arm64

	# update the deployment folder with the amd64 version
	# cp ./telstar-openweather-2.0/telstar-openweather-linux-amd64 ../../telstar-2-system/telstar-deployment/opt/telstar/plugins/openweather/telstar-openweather
	# cp ./README.md ../../telstar-2-system/telstar-deployment/opt/telstar/plugins/openweather/README.md
