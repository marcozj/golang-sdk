package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/workflowtype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

type GlobalWorkflow struct {
	client    *restapi.RestClient
	apiRead   string
	apiUpdate string

	ID       string                 `json:"ID,omitempty" schema:"id,omitempty"`
	Type     string                 `json:"key,omitempty" schema:"type,omitempty"` // wf, agentAuthWorkflow, secretsWorkflow, privilegeElevationWorkflow
	Settings *GlobalWorkflowSetting `json:"settings,omitempty" schema:"settings,omitempty"`
}

type GlobalWorkflowSetting struct {
	Enabled        bool               `json:"Enabled,omitempty" schema:"enabled,omitempty"`
	DefaultOptions string             `json:"DefaultOptions,omitempty" schema:"default_options,omitempty"`
	Approvers      string             `json:"Approvers,omitempty" schema:"approvers,omitempty"`
	ApproverList   []WorkflowApprover `json:"-" schema:"approver,omitempty"`
}

func NewGlobalWorkflow(c *restapi.RestClient, wfType string) (*GlobalWorkflow, error) {
	s := GlobalWorkflow{}
	s.client = c
	s.apiRead = "/ServerManage/GetSettings"
	s.apiUpdate = "/ServerManage/UpdateSettings"
	s.Type = wfType
	s.Settings = &GlobalWorkflowSetting{}
	//types := []string{workflowtype.AccountWorkflow.String(), workflowtype.AgentAuthWorkflow.String(), workflowtype.SecretsWorkflow.String(), workflowtype.PrivilegeElevationWorkflow.String()}
	if wfType != workflowtype.AccountWorkflow.String() && wfType != workflowtype.AgentAuthWorkflow.String() &&
		wfType != workflowtype.SecretsWorkflow.String() && wfType != workflowtype.PrivilegeElevationWorkflow.String() {
		errormsg := fmt.Sprintf("invalid workflow type %s", wfType)
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	return &s, nil
}

// Read function fetches global workflow settings
func (o *GlobalWorkflow) Read() error {
	var queryArg = make(map[string]interface{})
	queryArg["key"] = o.Type
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

	wfSettings := &GlobalWorkflowSetting{}
	mapToStruct(wfSettings, resp.Result)
	o.Settings = wfSettings

	return nil
}

func (o *GlobalWorkflow) Update() (*restapi.GenericMapResponse, error) {
	if o.Settings.Enabled && o.Settings.ApproverList != nil {
		err := ResolveWorkflowApprovers(o.client, o.Settings.ApproverList)
		if err != nil {
			return nil, err
		}
		// Due to historical reason, Approvers attribute is not in json format rather it is in string so need to perform conversion
		// Convert approvers from struct to string so that it can be assigned to the actual attribute used for privision.
		o.Settings.Approvers = FlattenWorkflowApprovers(o.Settings.ApproverList)
		//logger.Debugf("Converted approvers: %+v", o.Settings.Approvers)

		if o.Settings.DefaultOptions == "" {
			o.Settings.DefaultOptions = "{\"GrantMin\":60}"
		}
	}

	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

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

func (o *GlobalWorkflow) Delete() error {
	o.Settings.Enabled = false
	o.Settings.Approvers = ""

	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	logger.Debugf("Generated Map for Update(): %+v", queryArg)

	resp, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}

	return nil
}
