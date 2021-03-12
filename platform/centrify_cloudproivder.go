package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// CloudProvider - Encapsulates a cloud provider
type CloudProvider struct {
	vaultObject
	apiAddAccount            string
	apiSetAccountPermissions string

	CloudAccountID                            string          `json:"CloudAccountId,omitempty" schema:"cloud_account_id,omitempty"`
	Type                                      string          `json:"Type,omitempty" schema:"type,omitempty"`
	EnableUnmanagedPasswordRotation           bool            `json:"EnableUnmanagedPasswordRotation,omitempty" schema:"enable_interactive_password_rotation,omitempty"`
	EnableUnmanagedPasswordRotationPrompt     bool            `json:"EnableUnmanagedPasswordRotationPrompt,omitempty" schema:"prompt_change_root_password,omitempty"`
	EnableUnmanagedPasswordRotationReminder   bool            `json:"EnableUnmanagedPasswordRotationReminder,omitempty" schema:"enable_password_rotation_reminders,omitempty"`
	UnmanagedPasswordRotationReminderDuration int             `json:"UnmanagedPasswordRotationReminderDuration,omitempty" schema:"password_rotation_reminder_duration,omitempty"`
	ChallengeRules                            *ChallengeRules `json:"LoginRules,omitempty" schema:"challenge_rule,omitempty"`              // CloudProvider Login Challenge Rules
	LoginDefaultProfile                       string          `json:"LoginDefaultProfile,omitempty" schema:"default_profile_id,omitempty"` // Default CloudProvider Login Profile (used if no conditions matched)
}

// NewCloudProvider is a CloudProvider constructor
func NewCloudProvider(c *restapi.RestClient) *CloudProvider {
	s := CloudProvider{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Generic
	s.SetType = settype.CloudProvider.String()
	s.apiRead = "/CloudProvider/GetCloudProvider"
	s.apiCreate = "/CloudProvider/AddCloudProvider"
	s.apiDelete = "/CloudProvider/DeleteCloudProviders"
	s.apiUpdate = "/CloudProvider/UpdateCloudProvider"
	s.apiPermissions = "/CloudProvider/SetCloudProviderPermissions"
	s.apiAddAccount = "/ServerManage/AddAccount"
	s.apiSetAccountPermissions = "ServerManage/SetAccountPermissions"

	return &s
}

// Read function fetches a CloudProvider from source, including attribute values. Returns error if any
func (o *CloudProvider) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

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

// Create function creates a new CloudProvider and returns a map that contains creation result
func (o *CloudProvider) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	logger.Debugf("Generated Map for Create(): %+v", queryArg)

	resp, err := o.client.CallStringAPI(o.apiCreate, queryArg)
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
	o.ID = resp.Result

	return resp, nil
}

// Update function updates a existing CloudProvider and returns a map that contains update result
func (o *CloudProvider) Update() (*restapi.StringResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	logger.Debugf("Generated Map for Update(): %+v", queryArg)

	resp, err := o.client.CallStringAPI(o.apiUpdate, queryArg)
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

// Delete function deletes a CloudProvider and returns a string result that contains deletion result
func (o *CloudProvider) Delete() (*restapi.StringResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["Ids"] = []string{o.ID}
	queryArg["SaveToSecrets"] = false
	queryArg["SkipIfHasAppsOrServices"] = true

	resp, err := o.client.CallStringAPI(o.apiDelete, queryArg)
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

// Query function returns a single CloudProvider object in map format
func (o *CloudProvider) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM CloudProviders WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.CloudAccountID != "" {
		query += " AND CloudAccountId='" + o.CloudAccountID + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns CloudProvider ID by name
func (o *CloudProvider) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("CloudProvider name must be provided")
	}
	if o.CloudAccountID == "" {
		return "", fmt.Errorf("CloudProvider account id must be provided")
	}

	result, err := o.Query()
	if err != nil {
		errmsg := fmt.Sprintf("Error retrieving cloud provider: %s", err)
		logger.Errorf(errmsg)
		return "", fmt.Errorf(errmsg)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves CloudProvider from tenant by name
func (o *CloudProvider) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return fmt.Errorf("Failed to find ID of cloud provider %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a CloudProvider by name
func (o *CloudProvider) DeleteByName() (*restapi.StringResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return nil, fmt.Errorf("Failed to find ID of cloud provider %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

/*
Fetch Cloud Provider
	Request body format
	{
		"ID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"_entitycontext": "W/\"datetime'2021-01-10T02%3A11%3A07.3502066Z'\"",
			"Description": "jkljk",
			"EnableUnmanagedPasswordRotationReminder": true,
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"EnableUnmanagedPasswordRotation": true,
			"_TableName": "cloudproviders",
			"_encryptkeyid": "xxxxxxx",
			"_PartitionKey": "xxxxxxx",
			"_RowKey": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"UnmanagedPasswordRotationReminderDuration": 10,
			"Type": "Aws",
			"_Timestamp": "/Date(1610204132968)/",
			"Name": "My AWS",
			"EnableUnmanagedPasswordRotationPrompt": true,
			"ACL": "true",
			"ID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"CloudAccountId": "xxxxxxxxxx"
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Create Cloud Provider
	Request body format
	{
		"CloudAccountId": "xxxxxxxxxxxx",
		"Name": "Demo AWS",
		"Description": "Demo AWS Object",
		"Type": "Aws",
		"EnableUnmanagedPasswordRotation": true,
		"EnableUnmanagedPasswordRotationPrompt": true,
		"EnableUnmanagedPasswordRotationReminder": true,
		"UnmanagedPasswordRotationReminderDuration": 10
	}
	Respond result
	{
		"success": true,
		"Result": "312309c5-253f-4567-91f7-8dc07f002711",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

	Set CloudProvider Permissions
		Request body format
		{
			"Grants": [
				{
					"Principal": "System Administrator",
					"PType": "Role",
					"Rights": "View,Edit",
					"PrincipalId": "sysadmin"
				}
			],
			"ID": "9f60a10d-0426-49b9-a2d3-88aed47494f0"
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

	Add Root Account
	Request body format
	{
		"User": "xxxxxxxxx",
		"Password": "xxxxxxxxxxx",
		"CredentialType": "Password",
		"IsRootAccount": "true",
		"CloudProviderId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

	Set Account Permissions
	Request body format
	{
		"Grants": [
			{
				"Principal": "System Administrator",
				"PType": "Role",
				"Rights": "View,Login,Naked",
				"PrincipalId": "sysadmin"
			}
		],
		"ID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
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

Update Cloud Provider
	Request body format
	{
		"_encryptkeyid": "xxxxxxx",
		"_TableName": "cloudproviders",
		"ACL": "true",
		"Name": "My AWS",
		"_PartitionKey": "xxxxxxx",
		"ID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"_entitycontext": "W/\"datetime'2020-12-10T03%3A30%3A16.9539277Z'\"",
		"_RowKey": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"CloudAccountId": "xxxxxxxxxxx",
		"Type": "Aws",
		"_metadata": {
			"Version": 1,
			"IndexingVersion": 1
		},
		"EnableUnmanagedPasswordRotation": true,
		"updateChallenges": true,
		"Description": ""
	}

	Respond result
	{
		"success": true,
		"Result": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Delete Root Account
/ServerManage/DeleteAccounts
	Request body format
	{
		"Ids": [
			"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
		],
		"SaveToSecrets": false,
		"SecretName": "Bulk Account Delete - xxxxxxx@xxxxxxx.xxx - 1-9-2021 10:47:35 PM",
		"SetQuery": "",
		"RunSync": true,
		"SkipIfHasAppsOrServices": true
	}

	Respond result
	{
		"success": true,
		"Result": "Successfully deleted account xxxxxxx@xxxxxxx.xxx",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Delete Cloud Provider
	Request body format
	{
		"Ids": [
			"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
		],
		"SaveToSecrets": false,
		"SecretName": "Cloud Provider Delete - xxxxxxx@xxxxxxx.xxx - 1-9-2021 10:54:14 PM",
		"SetQuery": "",
		"RunSync": true,
		"SkipIfHasAppsOrServices": true
	}

	Respond result
	{
		"success": true,
		"Result": "Successfully deleted cloud provider Demo AWS",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

*/
