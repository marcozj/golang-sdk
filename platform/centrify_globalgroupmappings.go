package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// GroupMappings - Encapsulates Glboal Group Mappings
type GroupMappings struct {
	BulkUpdate bool `json:"-"`
	//Mappings   []GroupMapping `json:"Mappings,omitempty" schema:"mappings,omitempty"`
	Mappings []GroupMapping `json:"Mappings,omitempty" schema:"attribute_group,omitempty"`

	client     *restapi.RestClient
	apiRead    string
	apiCreate  string
	apiDelete  string
	apiUpdates string
}

// GroupMapping represents individual group mapping
type GroupMapping struct {
	AttributeValue string `json:"AttributeValue,omitempty" schema:"attribute_value,omitempty"`
	GroupName      string `json:"GroupName,omitempty" schema:"group_name,omitempty"`
}

// NewGroupMappings is a GroupMappings constructor
func NewGroupMappings(c *restapi.RestClient) *GroupMappings {
	s := GroupMappings{}
	s.client = c
	s.apiRead = "/Federation/GetGlobalGroupAssertionMappings"
	s.apiCreate = "/Federation/AddGlobalGroupAssertionMapping"
	s.apiDelete = "/Federation/DeleteGlobalGroupAssertionMapping"
	s.apiUpdates = "/Federation/UpdateGlobalGroupAssertionMappings"

	return &s
}

// Read function fetches Global Group mappings from tenant
func (o *GroupMappings) Read() error {

	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
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

// Create adds list of group mappings
func (o *GroupMappings) Create() error {
	err := o.createOrDelete(o.apiCreate)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes list of group mappings
func (o *GroupMappings) Delete() error {
	err := o.createOrDelete(o.apiDelete)
	if err != nil {
		return err
	}
	return nil
}

func (o *GroupMappings) createOrDelete(api string) error {
	for _, v := range o.Mappings {
		var queryArg = make(map[string]interface{})
		queryArg[v.AttributeValue] = v.GroupName
		logger.Debugf("Generated Map for %s: %+v", api, queryArg)

		resp, err := o.client.CallStringAPI(api, queryArg)
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
		if !resp.Success {
			errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
			logger.Errorf(errmsg)
			return fmt.Errorf(errmsg)
		}
	}
	return nil
}

func (o *GroupMappings) Update() error {
	var mappings []interface{}
	for _, v := range o.Mappings {
		var mapping = make(map[string]interface{})
		mapping["AttributeValue"] = v.AttributeValue
		mapping["GroupName"] = v.GroupName
		mappings = append(mappings, mapping)
	}

	var queryArg = make(map[string]interface{})
	queryArg["Mappings"] = mappings

	resp, err := o.client.CallGenericMapAPI(o.apiUpdates, queryArg)
	logger.Debugf("Generated Map for %s: %+v", o.apiUpdates, queryArg)
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

/*
Get mappings
	Request body format
	{
		"Args": {
			"PageNumber": 1,
			"PageSize": 100000,
			"Limit": 100000,
			"SortBy": "",
			"direction": "False",
			"Caching": -1
		}
	}

	Respond result
	{
		"success": true,
		"Result": {
			"Mappings": [
				{
					"AttributeValue": "Test Federated Group",
					"GroupName": "Azure PAS Users"
				},
				{
					"AttributeValue": "Test 2",
					"GroupName": "Okta PAS Admin"
				}
			]
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Add mapping
	Request body format
	{
		"Group from IdP": "PAS Group Name"
	}

	Respond result
	{
		"success": true,
		"Result": null,
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Delete mapping
	Request body format
	{
		"Group from IdP": "PAS Group Name"
	}

	Respond result
	{
		"success": true,
		"Result": null,
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
*/
