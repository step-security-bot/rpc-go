module rpc

go 1.20

// uncomment if developing with go-wsman-messages locally
replace github.com/open-amt-cloud-toolkit/go-wsman-messages => ../go-wsman-messages

require (
	github.com/gorilla/websocket v1.5.0
	github.com/ilyakaznacheev/cleanenv v1.4.2
	github.com/open-amt-cloud-toolkit/go-wsman-messages v1.3.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	golang.org/x/sys v0.10.0
)

require golang.org/x/crypto v0.11.0 // indirect

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
	software.sslmate.com/src/go-pkcs12 v0.2.1
)
