package webcookie

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	log "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// GetClient creates REST client
func (c *WebCookie) GetClient() (*restapi.RestClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Munge on the service a little bit, force it to have no trailing / and always start with https://
	url, err := url.Parse(c.Service)
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	url.Path = ""
	c.Service = url.String()

	if c.SkipCertVerify {
		// Ignore certificate error for on-prem deployment
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Client = &http.Client{Transport: tr}
	} else {
		c.Client = &http.Client{}
	}
	c.Client.Jar = jar

	log.Debugf("Start authentication...\n")
	authResp, err := c.startAuthentication()
	if err != nil {
		return nil, err
	}

	if !authResp.Success {
		return nil, fmt.Errorf("Failed to start authentication")
	}

	token, err := c.advanceAuthentication(authResp)
	if err != nil {
		return nil, err
	}

	var clientFactory restapi.HttpClientFactory = func() *http.Client {
		return &http.Client{}
	}
	if c.SkipCertVerify {
		// Ignore certificate error for on-prem deployment
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		clientFactory = func() *http.Client {
			return &http.Client{Transport: tr}
		}
	}
	restClient, err := restapi.GetNewRestClient(c.Service, clientFactory)
	if err != nil {
		return nil, err
	}

	restClient.Headers["Authorization"] = "Bearer " + token
	return restClient, nil
}
