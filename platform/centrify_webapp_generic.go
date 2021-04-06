package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

type GenericWebApp struct {
	WebApp

	Url string `json:"Url" schema:"url"` // The URL to bookmark
	// Advanced menu
	HostNameSuffix  string `json:"HostNameSuffix" schema:"hostname_suffix"`                                  // The host name suffix for the url of the login form, for example, acme.com.
	UsernameField   string `json:"UsernameField,omitempty" schema:"username_field,omitempty"`                // The CSS Selector for the user name field in the login form, for example, input#login-username.
	PasswordField   string `json:"PasswordField,omitempty" schema:"password_field,omitempty"`                // The CSS Selector for the password field in the login form, for example, input#login-password.
	SubmitField     string `json:"SubmitField,omitempty" schema:"submit_field,omitempty"`                    // The CSS Selector for the Submit button in the login form, for example, input#login-button. This entry is optional. It is required only if you cannot submit the form by pressing the enter key.
	FormField       string `json:"FormField,omitempty" schema:"form_field,omitempty"`                        // The CSS Selector for the form field of the login form, for example, form#loginForm.
	CorpIdField     string `json:"CorpIdField,omitempty" schema:"additional_login_field,omitempty"`          // The CSS Selector for any Additional Login Field required to login besides username and password, such as Company name or Agency ID. For example, the selector could be input#login-company-id. This entry is required only if there is an additional login field besides username and password.
	CorpIdentifier  string `json:"CorpIdentifier,omitempty" schema:"additional_login_field_value,omitempty"` // The value for the Additional Login Field. For example, if there is an additional login field for the company name, enter the company name here. This entry is required if Additional Login Field is set.
	SelectorTimeout int    `json:"SelectorTimeout,omitempty" schema:"selector_timeout,omitempty"`            // Use this field to indicate the number of milliseconds to wait for the expected input selectors to load before timing out on failure. A zero or negative number means no timeout.
	Order           string `json:"Order,omitempty" schema:"order,omitempty"`                                 // Use this field to specify the order of login if it is not username, password and submit.
	// For Browser Extension (advanced) app only
	Script string `json:"Script,omitempty" schema:"script,omitempty"` // Script to log the user in to this application
	// "UserPassScript": "@GenericUserPass" for User-Password app
	UseLoginPwAdAttr    bool   `json:"UseLoginPwAdAttr" schema:"use_ad_login_pw"` // Use the login password supplied by the user (Active Directory users only)
	Password            string `json:"Password,omitempty" schema:"password,omitempty"`
	UseLoginPwUseScript bool   `json:"UseLoginPwUseScript" schema:"use_ad_login_pw_by_script"`
}

func NewGenericWebApp(c *restapi.RestClient) *GenericWebApp {
	webapp := newWebpp(c)
	s := GenericWebApp{}
	s.WebApp = *webapp

	return &s
}

func (o *GenericWebApp) Read() error {
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

// Create function creates a new WebApp and returns a map that contains creation result
func (o *GenericWebApp) Create() (*restapi.SliceResponse, error) {
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

	return resp, nil
}

// Create function creates a new WebApp and returns a map that contains creation result
func (o *GenericWebApp) CreateComplete() error {
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
func (o *GenericWebApp) Update() (*restapi.GenericMapResponse, error) {
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

func (o *GenericWebApp) processWorkflow() error {
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
func (o *GenericWebApp) GetIDByName() (string, error) {
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
func (o *GenericWebApp) GetByName() error {
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
func (o *GenericWebApp) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Application WHERE 1=1 AND AppType='Web' AND WebAppType='UsernamePassword'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}
