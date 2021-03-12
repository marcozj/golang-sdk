package utils

import (
	"fmt"
	"strings"

	"github.com/marcozj/golang-sdk/dmc"
	"github.com/marcozj/golang-sdk/enum/authenticationtype"
	"github.com/marcozj/golang-sdk/oauth"
	"github.com/marcozj/golang-sdk/restapi"
	"github.com/marcozj/golang-sdk/webcookie"
)

// VaultClient represents vault client structure
type VaultClient struct {
	client   *restapi.RestClient // Authenticated REST client
	AuthType string              // Authentication type
	URL      string              // Tenant URL
	AppID    string              // OAuth2 application id
	Scope    string              // OAuth2 or DMC scope definition name
	Token    string              // OAuth2 or DMC token
	User     string              // User to run the command as (or OAuth2 client if requesting a token)
	Password string              // Password for user (or OAuth2 client secret if requesting a token)
	Skipcert bool                // Whether to skip certificate validation
	Debug    bool
}

// authenticate authenticates to tenant and save reset client
func (c *VaultClient) authenticate() error {
	var restClient *restapi.RestClient
	var err error
	switch strings.ToLower(c.AuthType) {
	case authenticationtype.OAuth2.String():
		call := oauth.OauthClient{
			Service:        c.URL,
			AppID:          c.AppID,
			Scope:          c.Scope,
			Token:          c.Token,
			ClientID:       c.User,
			ClientSecret:   c.Password,
			SkipCertVerify: c.Skipcert,
		}
		restClient, err = call.GetClient()
		if err != nil {
			return fmt.Errorf("Unable to get oauth rest client: %v", err)
		}
	case authenticationtype.UsernamePassword.String():
		call := webcookie.WebCookie{}
		call.Service = c.URL
		call.ClientID = c.User
		call.ClientSecret = c.Password
		call.SkipCertVerify = c.Skipcert

		restClient, err = call.GetClient()
		if err != nil {
			return fmt.Errorf("Unable to get simple rest client: %v", err)
		}
	case authenticationtype.DelegatedMachineCredential.String():
		call := dmc.DMC{}
		call.Service = c.URL
		call.Scope = c.Scope
		call.Token = c.Token
		call.SkipCertVerify = c.Skipcert

		restClient, err = call.GetClient()
		if err != nil {
			return fmt.Errorf("Unable to get DMC rest client: %v", err)
		}
	default:
		return fmt.Errorf("Invalid authentication type: %v", c.AuthType)
	}
	c.client = restClient
	return nil
}

// GetClient returns REST client
func (c *VaultClient) GetClient() (*restapi.RestClient, error) {
	if c.client == nil {
		err := c.authenticate()
		if err != nil {
			return nil, err
		}
	}
	return c.client, nil
}
