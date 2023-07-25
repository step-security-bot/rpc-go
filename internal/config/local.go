package config

type (
	Config struct {
		Password          string
		IEEE8021XSettings `yaml:"ieee801xConfig"`
		ACMSettings `yaml:"acmactivate"`
	}
	IEEE8021XSettings struct {
		Name                   string `yaml:"name"`
		AuthenticationMethod   int    `yaml:"authenticationMethod"`
		EncryptionMethod       int    `yaml:"encryptionMethod"`
		SSID                   string `yaml:"ssid"`
		Username               string `yaml:"username"`
		AuthenticationProtocol int    `yaml:"authenticationProtocol"`
		Priority               int    `yaml:"priority"`
		ClientCert             string `yaml:"clientCert"`
		CACert                 string `yaml:"caCert"`
		PrivateKey             string `yaml:"privateKey"`
	}
	ACMSettings struct {
		AMTPassword string `yaml:"amtPassword"`
		MEBxPassword string `yaml:"mebxPassword"`
		ProvisioningCert string `yaml:"provisioningCert"`
		ProvisioningCertPwd string `yaml:"provisioningCertPwd"`
	}
)
