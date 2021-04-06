package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

type WebApp struct {
	vaultObject
	apiResetAppScript string

	TemplateName string `json:"TemplateName,omitempty" schema:"template_name,omitempty"`
	// Policy menu
	DefaultAuthProfile string          `json:"DefaultAuthProfile,omitempty" schema:"default_profile_id,omitempty"`
	ChallengeRules     *ChallengeRules `json:"AuthRules,omitempty" schema:"challenge_rule,omitempty"`
	PolicyScript       string          `json:"PolicyScript" schema:"policy_script"` // Use script to specify authentication rules (configured rules are ignored)
	// Account Mapping menu
	UserNameStrategy string `json:"UserNameStrategy,omitempty" schema:"username_strategy,omitempty"` // ADAttribute, Fixed or UseScript
	//ADAttribute      string `json:"ADAttribute,omitempty" schema:"ad_attribute,omitempty"`           // Directory service field name. Used when UserNameStrategy=ADAttribute
	Username      string `json:"UserNameArg,omitempty" schema:"username,omitempty"` // Used when UserNameStrategy is ADAttribute or Fixed
	UserMapScript string `json:"UserMapScript" schema:"user_map_script"`            // Used when UserNameStrategy=UseScript
	// Workflow menu
	WorkflowEnabled      bool               `json:"WorkflowEnabled" schema:"workflow_enabled"`
	WorkflowSettings     string             `json:"WorkflowSettings,omitempty" schema:"workflow_settings"` // This is the actual workflow attribute in string format
	WorkflowApproverList []WorkflowApprover `json:"-" schema:"workflow_approver,omitempty"`                // This is used in tf file only
}

func newWebpp(c *restapi.RestClient) *WebApp {
	s := WebApp{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Application
	s.SetType = settype.Application.String()
	s.apiRead = "/SaasManage/GetApplication"
	s.apiCreate = "/SaasManage/ImportAppFromTemplate"
	s.apiDelete = "/SaasManage/DeleteApplication"
	s.apiUpdate = "/SaasManage/UpdateApplicationDE"
	s.apiPermissions = "/SaasManage/SetApplicationPermissions"
	s.apiResetAppScript = "/SaasManage/ResetAppScript"

	return &s
}

/*
// Read function fetches a WebApp from source, including attribute values. Returns error if any
func (o *WebApp) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = o.ID

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
*/
// Create function creates a new WebApp and returns a map that contains creation result
func (o *WebApp) Create() (*restapi.SliceResponse, error) {
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

	return resp, nil
}

/*
// Update function updates an existing WebApp and returns a map that contains update result
func (o *WebApp) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

		err := o.processSpMetaData()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, err
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
*/
// Delete function deletes a WebApp and returns a map that contains deletion result
func (o *WebApp) Delete() (*restapi.SliceResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = []string{o.ID}

	resp, err := o.client.CallSliceAPI(o.apiDelete, queryArg)
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

/*
func (o *WebApp) processWorkflow() error {
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

*/
// Query function returns a single WebApp object in map format
func (o *WebApp) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Application WHERE 1=1 AND AppType='Web'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns vault object ID by name
func (o *WebApp) GetIDByName() (string, error) {
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

/*
// GetByName retrieves vault object from tenant by name
func (o *WebApp) GetByName() error {
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
*/
// DeleteByName deletes a DesktopApp by name
func (o *WebApp) DeleteByName() (*restapi.SliceResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("failed to find ID of WebApp %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *WebApp) ResetAppScript() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = o.ID

	// Attempt to read from an upstream API
	_, err := o.client.CallGenericMapAPI(o.apiResetAppScript, queryArg)

	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}
