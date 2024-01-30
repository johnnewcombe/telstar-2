module bitbucket.org/johnnewcombe/telstar

go 1.16

// use the local library rather than the one in bitbucket
replace bitbucket.org/johnnewcombe/telstar-library => ../telstar-library

require (
	bitbucket.org/johnnewcombe/telstar-library v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.1
	github.com/go-chi/docgen v1.2.0 // indirect
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.3.0
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/lestrrat-go/jwx v1.1.0
	github.com/wagslane/go-password-validator v0.3.0 // indirect
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	rsc.io/quote v1.5.2 // indirect
)
