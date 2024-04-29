package restclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// target registry
const ()

// REST API request
type Api struct {
	method    string
	endpoint  string
	arguments []interface{}
	body      interface{}
}

// newRequest returns new restClient API request */
func NewRequest(method string, endpoint string, body interface{}, arguments ...interface{}) *Api {

	// log.Println("rest_service.go - handleRequest - request.body: ", body)
	return &Api{
		method:    method,
		endpoint:  endpoint,
		arguments: arguments,
		body:      body,
	}
}

// RestError represents a Schema Registry HTTP Error response
type RestError struct {
	Code    int    `json:"error_code"`
	Message string `json:"message"`
}

// Error implements the errors.Error interface
func (err *RestError) Error() string {
	return fmt.Sprintf("request failed error code: %d: %s", err.Code, err.Message)
}

type restService struct {
	url     *url.URL
	headers http.Header
	*http.Client
}

// newRestService returns a new REST client
func NewRestService(conf *Config, contentType string) (*restService, error) {
	urlConf := conf.TargetURL
	u, err := url.Parse(urlConf)

	if err != nil {
		return nil, err
	}

	headers, err := newAuthHeader(u, conf)
	if err != nil {
		return nil, err
	}

	fmt.Println("In NewRestService, see the headers: ", headers)

	headers.Add("Content-Type", contentType)
	if err != nil {
		return nil, err
	}

	transport, err := configureTransport(conf)
	if err != nil {
		return nil, err
	}

	timeout := conf.RequestTimeoutMs

	return &restService{
		url:     u,
		headers: headers,
		Client: &http.Client{
			Transport: transport,
			Timeout:   time.Duration(timeout) * time.Millisecond,
		},
	}, nil
}

// configureTransport returns a new Transport
func configureTransport(conf *Config) (*http.Transport, error) {

	// Exposed for testing purposes only. In production properly formed certificates should be used
	// https://tools.ietf.org/html/rfc2818#section-3
	tlsConfig := &tls.Config{}
	if err := configureTLS(conf, tlsConfig); err != nil {
		return nil, err
	}

	timeout := conf.ConnectionTimeoutMs

	return &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(timeout) * time.Millisecond,
		}).Dial,
		TLSClientConfig: tlsConfig,
	}, nil
}

// configureTLS populates tlsConf
func configureTLS(conf *Config, tlsConf *tls.Config) error {
	certFile := conf.SslCertificateLocation
	keyFile := conf.SslKeyLocation
	caFile := conf.SslCaLocation
	unsafe := conf.SslDisableEndpointVerification

	var err error
	if certFile != "" {
		if keyFile == "" {
			return errors.New(
				"SslKeyLocation needs to be provided if using SslCertificateLocation")
		}
		var cert tls.Certificate
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		tlsConf.Certificates = []tls.Certificate{cert}
	}

	if caFile != "" {
		if unsafe {
			log.Println("WARN: endpoint verification is currently disabled. " +
				"This feature should be configured for development purposes only")
		}
		var caCert []byte
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return err
		}

		tlsConf.RootCAs = x509.NewCertPool()
		if !tlsConf.RootCAs.AppendCertsFromPEM(caCert) {
			return fmt.Errorf("could not parse certificate from %s", caFile)
		}
	}

	tlsConf.BuildNameToCertificate()

	return err
}

// newAuthHeader returns a base64 encoded userinfo string identified on the configured credentials source
func newAuthHeader(service *url.URL, conf *Config) (http.Header, error) {
	// Remove userinfo from url regardless of source to avoid confusion/conflicts
	defer func() {
		service.User = nil
	}()

	source := conf.BasicAuthCredentialsSource

	header := http.Header{}

	var err error
	switch strings.ToUpper(source) {
	case "URL":
		err = configureURLAuth(service, header)
	case "BEARER":
		err = configureSpecificAuth(conf, header)
	default:
		err = fmt.Errorf("unrecognized value for basic.auth.credentials.source %s", source)
	}
	return header, err
}

// configureURLAuth copies the url userinfo into a basic HTTP auth authorization header
func configureURLAuth(service *url.URL, header http.Header) error {
	header.Add("Authorization", fmt.Sprintf("Basic %s", encodeBasicAuth(service.User.String())))
	return nil
}

func configureSpecificAuth(conf *Config, header http.Header) error {
	token := strings.TrimSpace(string(conf.AuthCredentialsDatas.Token))
	token = strings.Trim(token, `"`)
	header.Add("Authorization", fmt.Sprintf("%s %s", string(conf.AuthCredentialsDatas.AuthType), token))
	return nil
}

// encodeBasicAuth adds a basic http authentication header to the provided header
func encodeBasicAuth(userinfo string) string {
	return base64.StdEncoding.EncodeToString([]byte(userinfo))
}
