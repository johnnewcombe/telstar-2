# based on the fyne-cross image from lucor e.g.
#  go get -d github.com/fyne-io/fyne-cross/cmd/fyne-cross
# bbs.fozztexx.com (TELNET) UN:TELSTAR PW:J16
version = 1.6.2

build:
	# Update fyne-cross to use the latest docker images e.g.

    # Deprecated: go get -u github.com/fyne-io/fyne-cross
    # go get fyne.io/fyne/v2/app@v2.2.1
	./gen.sh bundle.go

	#go build

	# MacOS AMD & ARM
	fyne-cross darwin -app-id=com.glasstty.telstar-client -arch=amd64,arm64

	#cd ./fyne-cross/dist && chmod -R 777 licences
	#rm -r licences
	#go-licenses check .
	#go-licenses save . --save_path="licences"
	#go-licenses save "bitbucket.org/johnnewcombe/telstar-client" --save_path="tmp/licenses"

	# MacOS/AMD copy config files to app bundle resources
	#cp nxtel.yml ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/Resources
	#cp ccl4.yml ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/Resources
	#cp teefax.yml ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/Resources
	#cp telstar.yml ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/Resources
	#cp telstar-serial.yml ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/Resources
	#cp EOL-BBS.yml ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/Resources
	#cp -R ./fyne-cross/dist/licences ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents

	# MacOS/ARM copy config files to app bundle resources
	#cp nxtel.yml ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/Resources
	#cp ccl4.yml ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/Resources
	#cp teefax.yml ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/Resources
	#cp telstar.yml ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/Resources
	#cp telstar-serial.yml ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/Resources
	#cp EOL-BBS.yml ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/Resources
	#cp -R ./fyne-cross/dist/licences ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents

	# Update bin folder for completeness
	cp ./fyne-cross/dist/darwin-amd64/telstar-client.app/Contents/MacOS/telstar-client ./fyne-cross/bin/darwin-amd64
	cp ./fyne-cross/dist/darwin-arm64/telstar-client.app/Contents/MacOS/telstar-client ./fyne-cross/bin/darwin-arm64

	# zip MAC-OS Versions and remove the .app version (copy to applications first)
	cd ./fyne-cross/dist/darwin-amd64 && cp -r -v ./telstar-client.app /Applications
	cd ./fyne-cross/dist/darwin-amd64 && zip -r -v ./telstar-client.app.zip ./telstar-client.app
	cd ./fyne-cross/dist && rm -R ./darwin-amd64/telstar-client.app
	cd ./fyne-cross/dist/darwin-arm64 && zip -r -v ./telstar-client.app.zip ./telstar-client.app
	cd ./fyne-cross/dist && rm -R ./darwin-arm64/telstar-client.app


	# Windows
	fyne-cross windows -pull -app-id=com.glasstty.telstar-client -arch=amd64
	fyne-cross windows -pull -app-id=com.glasstty.telstar-client -arch=386



	# Linux AMD64

	#######################################################################
	# NOTE NOTE due to the packaging error (see below) the command is followed
    # by ';true' to prevent the error from aborting Make. A manual cp to dist
    # has also been added.
	# Issue Raised: https://github.com/fyne-io/fyne-cross/issues/85
	#######################################################################
	#fyne-cross linux -pull -app-id=com.glasstty.telstar-client -arch=amd64
	fyne-cross linux -pull -app-id=com.glasstty.telstar-client -arch=amd64; true
	cp ./fyne-cross/tmp/linux-amd64/telstar-client.tar.xz ./fyne-cross/dist/linux-amd64/

	# Linux ARM64

	#######################################################################
	# NOTE NOTE due to the packaging error (see below) the command is followed
    # by ';true' to prevent the error from aborting Make. A manual cp to dist
    # has also been added.
	# Issue Raised: https://github.com/fyne-io/fyne-cross/issues/85
	#######################################################################
	fyne-cross linux -pull -app-id=com.glasstty.telstar-client -arch=arm64; true
	cp ./fyne-cross/tmp/linux-arm64/telstar-client.tar.xz ./fyne-cross/dist/linux-arm64/

	# Linux ARM

	#######################################################################
	# NOTE NOTE due to the packaging error (see below) the command is followed
    # by ';true' to prevent the error from aborting Make. A manual cp to dist
    # has also been added.
	# Issue Raised: https://github.com/fyne-io/fyne-cross/issues/85
	#######################################################################
	fyne-cross linux -pull -app-id=com.glasstty.telstar-client -arch=arm; true
	cp ./fyne-cross/tmp/linux-arm/telstar-client.tar.xz ./fyne-cross/dist/linux-arm/


	# Connection Files
	#cp nxtel.yml ./fyne-cross/dist/example-connection-files/
	#cp ccl4.yml ./fyne-cross/dist/example-connection-files/
	#cp teefax.yml ./fyne-cross/dist/example-connection-files/
	#cp telstar.yml ./fyne-cross/dist/example-connection-files/
	#cp telstar-serial.yml ./fyne-cross/dist/example-connection-files/
	#cp EOL-BBS.yml ./fyne-cross/dist/example-connection-files/

	# remove previous .zip file
	cd ./fyne-cross && touch telstar-client.zip
	cd ./fyne-cross && rm telstar-client.zip

	# zip up the client for download convenience
	cp -R ./fyne-cross/dist/ ./fyne-cross/telstar-client-$(version)
	cd ./fyne-cross && zip -vr ./telstar-client.zip ./telstar-client-$(version)
	#cd ./fyne-cross && zip -vr ./example-connection-files.zip ./telstar-client-$(version)/example-connection-files

	# remove previous .zip file
	cd ./fyne-cross && touch telstar-client.zip
	cd ./fyne-cross && rm telstar-client.zip

	# zip up the client for download convenience
	cp -R ./fyne-cross/dist/ ./fyne-cross/telstar-client-$(version)
	cd ./fyne-cross && zip -vr ./telstar-client.zip ./telstar-client-$(version)
	#cd ./fyne-cross && zip -vr ./example-connection-files.zip ./telstar-client-$(version)/example-connection-files
