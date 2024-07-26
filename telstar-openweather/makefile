build:
	env GOOS=darwin GOARCH=amd64  go build -o /opt/telstar/telstar-openweather .
	env GOOS=darwin GOARCH=amd64  go build -o ./bin/telstar-openweather-macos-amd64 .
	env GOOS=darwin GOARCH=arm64  go build -o ./bin/telstar-openweather-macos-arm64 .
	env GOOS=windows GOARCH=386 go build -o ./bin/telstar-openweather-windows-386.exe .
	env GOOS=windows GOARCH=amd64 go build -o ./bin/telstar-openweather-windows-amd64.exe .
	env GOOS=linux GOARCH=amd64   go build -o ./bin/telstar-openweather-linux-amd64 .
	env GOOS=linux GOARCH=arm64   go build -o ./bin/telstar-openweather-linux-arm64 .

	zip -v ./bin/telstar-openweather.zip  ./bin/telstar-openweather-macos-amd64 ./bin/telstar-openweather-macos-arm64 ./bin/telstar-openweather-windows-amd64.exe ./bin/telstar-openweather-linux-amd64 ./bin/telstar-openweather-linux-arm64
