package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

type OidcWebApp struct {
	WebApp

	// Setting menu
	ApplicationID string `json:"ServiceName,omitempty" schema:"application_id,omitempty"`

	OAuthProfile        *OidcProfile `json:"OAuthProfile,omitempty" schema:"oauth_profile,omitempty"`
	Script              string       `json:"Script,omitempty" schema:"script,omitempty"`                   // Script to generate OpenID Connect Authorization and UserInfo responses for this application
	OpenIDConnectScript string       `json:"OpenIDConnectScript,omitempty" schema:"oidc_script,omitempty"` // Read only attribute
}

type OidcProfile struct {
	// Trust menu
	ClientSecret string   `json:"ClientSecret,omitempty" schema:"client_secret,omitempty"` // The OpenID Client Secret for this Identity Provider
	Url          string   `json:"Url,omitempty" schema:"application_url,omitempty"`        // The OpenID Connect Service Provider URL
	Redirects    []string `json:"Redirects,omitempty" schema:"redirects,omitempty"`        // Redirect URI that the Service Provider will specify in the OpenID Connect request to Centrify
	// Read only attributes
	ClientID string `json:"ClientID,omitempty" schema:"client_id,omitempty"` // The OpenID Client ID for this Identity Provider
	Issuer   string `json:"Issuer,omitempty" schema:"issuer,omitempty"`      // The OpenID Connect Issuer URL for this application

	// Tokens menu
	TokenLifetime   string `json:"TokenLifetimeString,omitempty" schema:"token_lifetime,omitempty"`     // 5 hours "5:00:00"
	AllowRefresh    bool   `json:"AllowRefresh,omitempty" schema:"allow_refresh,omitempty"`             // Issue refresh tokens
	RefreshLifetime string `json:"RefreshLifetimeString,omitempty" schema:"refresh_lifetime,omitempty"` // 365 days "365.00:00:00"
}

func NewOidcWebApp(c *restapi.RestClient) *OidcWebApp {
	webapp := newWebpp(c)
	s := OidcWebApp{}
	s.WebApp = *webapp
	s.OAuthProfile = &OidcProfile{}

	return &s
}

func (o *OidcWebApp) Read() error {
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
	// This is annoying. "Script" attribute is used for update but "OpenIDConnectScript" attribute is used for read
	// So, assign value of "OpenIDConnectScript" to "Script"
	o.Script = o.OpenIDConnectScript

	return nil
}

// Create function creates a new WebApp and returns a map that contains creation result
func (o *OidcWebApp) Create() (*restapi.SliceResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = []string{o.TemplateName}
	logger.Debugf("Generated Map for Create(): %+v", queryArg)

	resp, err := o.client.CallSliceAPI(o.apiCreate, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	// Assign ID after successful creation so that the same object can be used for subsequent update operation
	o.ID = resp.Result[0].(map[string]interface{})["_RowKey"].(string)

	// After creation, read it back. This is to retrieve genreated ClientID attribute
	obj := NewOidcWebApp(o.client)
	obj.ID = o.ID
	obj.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	o.OAuthProfile.ClientID = obj.OAuthProfile.ClientID

	return resp, nil
}

// Create function creates a new WebApp and returns a map that contains creation result
func (o *OidcWebApp) CreateComplete() error {
	_, err := o.Create()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	// Perform update
	_, err = o.Update()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	// Assign to Set
	if len(o.Sets) > 0 {
		err := o.AddToSetsByName(o.Sets)
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
	}

	// Assign permissions
	if len(o.Permissions) > 0 {
		err := ResolvePermissions2(o.client, o.Permissions, o.ValidPermissions)
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
		_, err = o.SetPermissions(false)
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
	}

	return nil
}

// Update function updates an existing WebApp and returns a map that contains update result
func (o *OidcWebApp) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	err := o.processWorkflow()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
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

func (o *OidcWebApp) processWorkflow() error {
	// Resolve guid of each approver
	if o.WorkflowEnabled && o.WorkflowApproverList != nil {
		err := ResolveWorkflowApprovers(o.client, o.WorkflowApproverList)
		if err != nil {
			return err
		}
		// Due to historical reason, WorkflowSettings attribute is not in json format rather it is in string so need to perform conversion
		// Convert approvers from struct to string so that it can be assigned to the actual attribute used for privision.
		wfApprovers := FlattenWorkflowApprovers(o.WorkflowApproverList)
		o.WorkflowSettings = "{\"WorkflowApprover\":" + wfApprovers + "}"
	}
	return nil
}

// GetIDByName returns vault object ID by name
func (o *OidcWebApp) GetIDByName() (string, error) {
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
func (o *OidcWebApp) GetByName() error {
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
func (o *OidcWebApp) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Application WHERE 1=1 AND AppType='Web' AND WebAppType='OpenIDConnect'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.ApplicationID != "" {
		query += " AND ServiceName='" + o.ApplicationID + "'"
	}

	return queryVaultObject(o.client, query)
}
