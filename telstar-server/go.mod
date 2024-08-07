module github.com/johnnewcombe/telstar

go 1.18

// use the local library rather than the one in bitbucket
replace github.com/johnnewcombe/telstar-library => ../telstar-library

require (
	github.com/go-chi/chi/v5 v5.0.1
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-chi/render v1.0.3
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/johnnewcombe/telstar-library v0.0.0-20220314174241-c5c80a4acadb
	github.com/lestrrat-go/jwx v1.1.0
	go.mongodb.org/mongo-driver v1.16.0
	golang.org/x/crypto v0.22.0
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/goccy/go-json v0.3.5 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.7 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.0 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/wagslane/go-password-validator v0.3.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	olympos.io/encoding/edn v0.0.0-20200308123125-93e3b8dd0e24 // indirect
)
