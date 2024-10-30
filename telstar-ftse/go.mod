module github.com/johnnewcombe/telstar-ftse

go 1.17

// use the local library rather than the one in bitbucket
replace github.com/johnnewcombe/telstar-library => ../telstar-library

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/johnnewcombe/telstar-library v0.0.0-00010101000000-000000000000
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/wagslane/go-password-validator v0.3.0 // indirect
	go.mongodb.org/mongo-driver v1.17.1 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
)
