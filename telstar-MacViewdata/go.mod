module telstar-gnome

go 1.17

// use the local library rather than the one in bitbucket
replace github.com/johnnewcombe/telstar-library => ../telstar-library

require github.com/johnnewcombe/telstar-library v0.0.0-00010101000000-000000000000
