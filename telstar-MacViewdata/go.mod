module telstar-gnome

go 1.17

// use the local library rather than the one in bitbucket
replace bitbucket.org/johnnewcombe/telstar-library => ../telstar-library

require (
	bitbucket.org/johnnewcombe/telstar-library v0.0.0-00010101000000-000000000000
	github.com/PuerkitoBio/goquery v1.8.0
)
