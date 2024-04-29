package restclient

import (
// "fmt"
)

type AuthType string

const (
	Basic  AuthType = "Basic"
	Bearer AuthType = "Bearer"
)

type AuthData struct {
	AuthType AuthType
	Token    string
}

type Config struct {
	// SchemaRegistryURL determines the URL of the service to reach
	TargetURL string

	// BasicAuthUserInfo specifies the user info in the form of {username}:{password}.
	BasicAuthUserInfo string
	// BasicAuthCredentialsSource specifies how to determine the credentials, one of URL, USER_INFO
	BasicAuthCredentialsSource string

	// AuthCreadentialsDatas holds the creadentials datas
	AuthCredentialsDatas AuthData

	// SslCertificateLocation specifies the location of SSL certificates.
	SslCertificateLocation string
	// SslKeyLocation specifies the location of SSL keys.
	SslKeyLocation string
	// SslCaLocation specifies the location of SSL certificate authorities.
	SslCaLocation string
	// SslDisableEndpointVerification determines whether to disable endpoint verification.
	SslDisableEndpointVerification bool

	// ConnectionTimeoutMs determines the connection timeout in milliseconds.
	ConnectionTimeoutMs int
	// RequestTimeoutMs determines the request timeout in milliseconds.
	RequestTimeoutMs int
}

func NewConfig(url string, authDatas ...AuthData) *Config {
	c := &Config{}
	c.TargetURL = url

	c.BasicAuthUserInfo = ""
	c.BasicAuthCredentialsSource = "URL"
	if len(authDatas) > 0 {
		c.BasicAuthCredentialsSource = string(authDatas[0].AuthType)
		c.AuthCredentialsDatas = authDatas[0]
	}

	c.SslCertificateLocation = ""
	c.SslKeyLocation = ""
	c.SslCaLocation = ""
	c.SslDisableEndpointVerification = false

	c.ConnectionTimeoutMs = 10000
	c.RequestTimeoutMs = 10000

	return c
}
