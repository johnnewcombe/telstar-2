module bitbucket.org/johnnewcombe/telstar-ftse

go 1.17

// use the local library rather than the one in bitbucket
replace bitbucket.org/johnnewcombe/telstar-library => ../telstar-library

require (
	bitbucket.org/johnnewcombe/telstar-library v0.0.0-00010101000000-000000000000
	github.com/PuerkitoBio/goquery v1.8.0
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/wagslane/go-password-validator v0.3.0 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
)
