# Set this to the build version
version := 2.5.RC1

# Remember to set the ENV REFRESHED_AT variable in the associated docker files.

build:
	echo $(version) > ./telstar-server/version.txt
	make -C ./telstar-server
	make -C ./telstar-util
	make -C ./telstar-rss
	make -C ./telstar-upload
	make -C ./telstar-emf
# Requires SUDOcd
	make -C ./telstar-openweather

#   Uncomment this line when fyne-cross is fixed
#	make -C ./telstar-client

docker: #build
	# Build the architectures
	# this adds the latest tag to this release
	docker build -f Dockerfile.amd64 --rm --no-cache --tag johnnewcombe/telstar:latest --tag johnnewcombe/telstar:amd64-$(version) .
	docker build -f Dockerfile.arm64v8 --rm --no-cache --tag johnnewcombe/telstar:arm64v8-$(version) .

docker-push: #docker
	docker push johnnewcombe/telstar:amd64-$(version)
	docker push johnnewcombe/telstar:arm64v8-$(version)

	# TODO Look at creating a manifest list so that the correct architecture version is pulled from docker hub automatically.
	# docker manifest create johnnewcombe/telstar:latest johnnewcombe/telstar:amd64-2.0.0 johnnewcombe/telstar:arm64v8-2.0.0
	docker push johnnewcombe/telstar:latest

