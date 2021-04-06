package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/webapp/oauth/applicationtemplate"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

type OauthWebApp struct {
	WebApp

	// Setting menu
	ApplicationID string        `json:"ServiceName,omitempty" schema:"application_id,omitempty"`
	OAuthProfile  *OAuthProfile `json:"OAuthProfile,omitempty" schema:"oauth_profile,omitempty"`
	Script        string        `json:"Script,omitempty" schema:"script,omitempty"` // Script to customize JWT token creation for this application
}

type OAuthProfile struct {
	// General Usage menu
	TargetIsUs bool `json:"TargetIsUs,omitempty" schema:"target_is_us,omitempty"` // Set to true for OAuth Client. Set to false for OAuth Server
	//ClientIDType      string   `json:"ClientIDType,omitempty" schema:"clientid_type,omitempty"` // anything, list, confidential
	ClientIDType      int      `json:"ClientIDType,omitempty" schema:"clientid_type,omitempty"`
	Issuer            string   `json:"Issuer,omitempty" schema:"issuer,omitempty"`
	Audience          string   `json:"Audience,omitempty" schema:"audience,omitempty"`
	AllowedClients    []string `json:"AllowedClients,omitempty" schema:"allowed_clients,omitempty"`      // Applicable if ClientIDType is list
	AllowPublic       bool     `json:"AllowPublic,omitempty" schema:"allow_public,omitempty"`            // Set to true if ClientIDType is list
	MustBeOauthClient bool     `json:"MustBeOauthClient,omitempty" schema:"must_oauth_client,omitempty"` // Applicable if ClientIDType is confidential
	Redirects         []string `json:"Redirects,omitempty" schema:"redirects,omitempty"`
	// Tokens menu
	TokenType       string `json:"TokenType,omitempty" schema:"token_type,omitempty"`                   // JwtRS256, Opaque
	AllowedAuth     string `json:"AllowedAuth,omitempty" schema:"allowed_auth,omitempty"`               // AuthorizationCode,Implicit,ClientCreds,ResourceCreds
	TokenLifetime   string `json:"TokenLifetimeString,omitempty" schema:"token_lifetime,omitempty"`     // 5 hours "5:00:00"
	AllowRefresh    bool   `json:"AllowRefresh,omitempty" schema:"allow_refresh,omitempty"`             // Issue refresh tokens
	RefreshLifetime string `json:"RefreshLifetimeString,omitempty" schema:"refresh_lifetime,omitempty"` // 365 days "365.00:00:00"
	// Scope menu
	ConfirmAuthorization bool         `json:"Confirm,omitempty" schema:"confirm_authorization,omitempty"`       // User must confirm authorization request
	AllowScopeSelect     bool         `json:"AllowScopeSelect,omitempty" schema:"allow_scope_select,omitempty"` // Allow scope selection
	KnownScopes          []OAuthScope `json:"KnownScopes,omitempty" schema:"scope,omitempty"`
}

type OAuthScope struct {
	Name            string   `json:"Scope,omitempty" schema:"name,omitempty"`
	Description     string   `json:"Description,omitempty" schema:"description,omitempty"`
	AllowedRestAPIs []string `json:"AllowedRest,omitempty" schema:"allowed_rest_apis,omitempty"`
}

func NewOauthWebApp(c *restapi.RestClient) *OauthWebApp {
	webapp := newWebpp(c)
	s := OauthWebApp{}
	s.WebApp = *webapp

	return &s
}

func (o *OauthWebApp) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = o.ID

	logger.Debugf("Generated Map for Read(): %+v", queryArg)
	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)

	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}

	mapToStruct(o, resp.Result)

	return nil
}

// Update function updates an existing WebApp and returns a map that contains update result
func (o *OauthWebApp) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	o.processOauthProfile()

	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	queryArg["_RowKey"] = o.ID

	logger.Debugf("Generated Map for Update(): %+v", queryArg)

	resp, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	return resp, nil
}

// GetIDByName returns vault object ID by name
func (o *OauthWebApp) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("%s name must be provided", GetVarType(o))
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("error retrieving %s: %s", GetVarType(o), err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves vault object from tenant by name
func (o *OauthWebApp) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("failed to find ID of %s %s. %v", GetVarType(o), o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// Query function returns a single WebApp object in map format
func (o *OauthWebApp) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Application WHERE 1=1 AND AppType='Web' AND WebAppType='OAuth'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.ApplicationID != "" {
		query += " AND ServiceName='" + o.ApplicationID + "'"
	}

	return queryVaultObject(o.client, query)
}

func (o *OauthWebApp) processOauthProfile() {
	// If this is OAuth client, force TargetIsUs attribute to be true otherwise it will be OAuth Server
	//if o.TemplateName == applicationtemplate.OAuth2Client.String() && !o.OAuthProfile.TargetIsUs {
	if o.TemplateName == applicationtemplate.OAuth2Client.String() {
		o.OAuthProfile.TargetIsUs = true
		//} else if o.TemplateName == applicationtemplate.OAuth2Server.String() && o.OAuthProfile.TargetIsUs {
	} else if o.TemplateName == applicationtemplate.OAuth2Server.String() {
		o.OAuthProfile.TargetIsUs = false
	}

	// If ClientIDType is list or anything, set AllowPublic to true
	if o.OAuthProfile.ClientIDType == 0 {
		o.OAuthProfile.AllowPublic = true
	} else {
		o.OAuthProfile.AllowPublic = false
	}
}
