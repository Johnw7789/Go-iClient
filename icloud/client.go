package icloud

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Client struct {
	HttpClient tls_client.HttpClient

	OtpChannel chan string

	Username string
	Password string

	authToken  string
	trustToken string
	frameId    string
	clientId   string
	authAttr   string
	serviceKey string
	sessionID  string
	scnt       string
}

// * NewClient intializes a new http client and returns a new icloud client
func NewClient(username, password string, sniff bool) (*Client, error) {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithCookieJar(jar),
		tls_client.WithClientProfile(profiles.Chrome_124),
	}

	if sniff {
		options = append(options, tls_client.WithProxyUrl("http://localhost:8888"))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		HttpClient: client,
		OtpChannel: make(chan string),
		Username:   username,
		Password:   password,
	}, nil
}
