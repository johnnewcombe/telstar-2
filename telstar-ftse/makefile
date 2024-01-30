
build:
	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-ftse-2.0/local/bin/telstar-ftse main.go
	env GOOS=darwin GOARCH=amd64  go build -o ./telstar-ftse-2.0/telstar-ftse-macos-amd64 main.go
	env GOOS=darwin GOARCH=arm64  go build -o ./telstar-ftse-2.0/telstar-ftse-macos-arm64 main.go
	env GOOS=windows GOARCH=amd64  go build -o ./telstar-ftse-2.0/telstar-ftse-windows-amd64.exe main.go
	env GOOS=linux GOARCH=amd64  go build -o ./telstar-ftse-2.0/telstar-ftse-linux-amd64 main.go
	env GOOS=linux GOARCH=arm64  go build -o ./telstar-ftse-2.0/telstar-ftse-linux-arm64 main.go

	zip -v ./telstar-ftse-2.0/telstar-ftse-2.0.zip  ./telstar-ftse-2.0/telstar-ftse-macos-amd64 ./telstar-ftse-2.0/telstar-ftse-macos-arm64 ./telstar-ftse-2.0/telstar-ftse-windows-amd64.exe ./telstar-ftse-2.0/telstar-ftse-linux-amd64 ./telstar-ftse-2.0/telstar-ftse-linux-arm64

	# update the deployment folder with the amd64 version
	cp ./telstar-ftse-2.0/telstar-ftse-linux-amd64 ../../telstar-2-system/telstar-deployment/opt/telstar/telstar-ftse
	cp ./data/templates/* ../../telstar-2-system/telstar-deployment/opt/telstar/data/templates/ftse
