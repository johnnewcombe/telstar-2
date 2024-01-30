module bitbucket.org/johnnewcombe/telstar-util

go 1.17

// use the local library rather than the one in bitbucket
replace bitbucket.org/johnnewcombe/telstar-library => ../telstar-library

require bitbucket.org/johnnewcombe/telstar-library v0.0.0-20220314174241-c5c80a4acadb

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/wagslane/go-password-validator v0.3.0 // indirect
)
