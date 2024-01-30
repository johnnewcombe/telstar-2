module bitbucket.org/johnnewcombe/telstar-emf

go 1.17

// use the local library rather than the one in bitbucket
replace bitbucket.org/johnnewcombe/telstar-library => ../telstar-library

require (
	bitbucket.org/johnnewcombe/telstar-library v0.0.0-20220314174241-c5c80a4acadb
	github.com/mmcdole/gofeed v1.1.3
)

require (
	github.com/PuerkitoBio/goquery v1.5.1 // indirect
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mmcdole/goxpp v0.0.0-20181012175147-0068e33feabf // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/wagslane/go-password-validator v0.3.0 // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/text v0.3.7 // indirect
)
