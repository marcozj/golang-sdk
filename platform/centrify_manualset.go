package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// ManualSet - Encapsulates a single Generic ManualSet
type ManualSet struct {
	vaultObject
	apiUpdateMembers       string
	apiMemberPermissions   string
	ValidMemberPermissions map[string]string

	ObjectType        string `json:"ObjectType,omitempty" schema:"type,omitempty"`
	SubObjectType     string `json:"SubObjectType,omitempty" schema:"subtype,omitempty"`
	CollectionType    string `json:"CollectionType,omitempty" schema:"collection_type,omitempty"`
	MemberPermissions []Permission
}

// setMember is a simple struct used for performing Set related API call
type setMember struct {
	MemberType string
	Table      string
	Key        string
}

// NewManualSet is a ManualSet constructor
func NewManualSet(c *restapi.RestClient) *ManualSet {
	s := ManualSet{}
	s.CollectionType = "ManualBucket"
	s.ValidPermissions = ValidPermissionMap.Set
	s.client = c
	s.apiRead = "/Collection/GetCollection"
	s.apiCreate = "/Collection/CreateManualCollection"
	s.apiDelete = "/Collection/DeleteCollection"
	s.apiUpdate = "/Collection/UpdateCollection"
	s.apiUpdateMembers = "/Collection/UpdateMembersCollection"
	s.apiPermissions = "/Collection/SetCollectionPermissions"

	return &s
}

// NewManualSetWithType is another ManualSet constructor that initialise memberpermissions api endpiont
func NewManualSetWithType(c *restapi.RestClient, setType string) (*ManualSet, error) {
	s := NewManualSet(c)
	s.ObjectType = setType
	var err error
	err = s.ResolveValidMemberPerms()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return s, nil
}

// Read function fetches a ManualSet from source, including attribute values. Returns error if any
func (o *ManualSet) Read() error {
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

// Create function creates a new Manual Set and returns a map that contains creation result
func (o *ManualSet) Create() (*restapi.StringResponse, error) {
	// If ObjectType is "Application", SubObjectType to be set to either "Desktop" or "Web"
	// If SubObjectType is not set, the set will be visible for both Web and Desktop applications
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

// Delete function deletes a Manual Set and returns a map that contains deletion result
func (o *ManualSet) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing Manual Set and returns a map that contains update result
func (o *ManualSet) Update() (*restapi.GenericMapResponse, error) {
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

// Query function returns a single Set object in map format
func (o *ManualSet) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Sets WHERE 1=1"
	if o.ObjectType != "" {
		query += " AND ObjectType='" + o.ObjectType + "'"
	}
	if o.CollectionType != "" {
		query += " AND CollectionType='" + o.CollectionType + "'"
	}
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// UpdateSetMembers adds or removes members from the ManualSet
func (o *ManualSet) UpdateSetMembers(ids []string, action string) (*restapi.StringResponse, error) {
	var resp *restapi.StringResponse
	var err error
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	if action != "add" && action != "remove" {
		return nil, fmt.Errorf("Update Set action must be either 'add' or 'remove'")
	}
	if o.ObjectType == "" {
		return nil, fmt.Errorf("Set object type is empty")
	}

	for _, v := range ids {
		var members []setMember
		var queryArg = make(map[string]interface{})
		queryArg["id"] = o.ID // This is the ID of Set
		setData := setMember{
			Key:        v,
			MemberType: "Row",
			Table:      o.ObjectType,
		}
		members = append(members, setData)
		queryArg[action] = members
		resp, err = o.client.CallStringAPI(o.apiUpdateMembers, queryArg)
		if err != nil {
			logger.Errorf(err.Error())
			return nil, err
		}
		if !resp.Success {
			errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
			logger.Errorf(errmsg)
			return nil, fmt.Errorf(errmsg)
		}
	}

	return resp, nil
}

// SetMemberPermissions sets member permissions. isRemove indicates whether to remove all permissions instead of setting permissions
func (o *ManualSet) SetMemberPermissions(isRemove bool) (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	var permissions []map[string]interface{}
	for _, v := range o.MemberPermissions {
		var permission = make(map[string]interface{})
		permission, err := generateRequestMap(v)
		if isRemove {
			permission["Rights"] = "None"
		}
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if len(permissions) > 0 {
		var queryArg = make(map[string]interface{})
		queryArg["ID"] = o.ID
		queryArg["PVID"] = o.ID
		queryArg["RowKey"] = o.ID
		queryArg["Grants"] = permissions
		logger.Debugf("Generated Map for SetMemberPermissions(): %+v", queryArg)
		resp, err := o.client.CallGenericMapAPI(o.apiMemberPermissions, queryArg)
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
	return nil, nil
}

// ResolveValidMemberPerms returns member permission list and set member permission API endpoint according to type of resource
func (o *ManualSet) ResolveValidMemberPerms() error {
	switch o.ObjectType {
	case settype.System.String():
		o.ValidMemberPermissions = ValidPermissionMap.WinNix
		o.apiMemberPermissions = "/ServerManage/SetResourceCollectionPermissions"
	case settype.Database.String():
		o.ValidMemberPermissions = ValidPermissionMap.Database
		o.apiMemberPermissions = "/ServerManage/SetDatabaseCollectionPermissions"
	case settype.Domain.String():
		o.ValidMemberPermissions = ValidPermissionMap.Domain
		o.apiMemberPermissions = "/ServerManage/SetDomainCollectionPermissions"
	case settype.Account.String():
		o.ValidMemberPermissions = ValidPermissionMap.Account
		o.apiMemberPermissions = "/ServerManage/SetAccountCollectionPermissions"
	case settype.Secret.String():
		o.ValidMemberPermissions = ValidPermissionMap.Secret
		o.apiMemberPermissions = "/ServerManage/SetSecretCollectionPermissions"
	case settype.SSHKey.String():
		o.ValidMemberPermissions = ValidPermissionMap.SSHKey
		o.apiMemberPermissions = "/ServerManage/SetSSHKeyCollectionPermissions"
	case settype.Service.String():
		o.ValidMemberPermissions = ValidPermissionMap.Service
		o.apiMemberPermissions = "/Subscriptions/SetSubscriptionCollectionPermissions"
	case settype.Application.String():
		o.ValidMemberPermissions = ValidPermissionMap.Application
		o.apiMemberPermissions = "/SaasManage/SetApplicationCollectionPermissions"
	case settype.ResourceProfile.String():
		o.ValidMemberPermissions = ValidPermissionMap.Set
		o.apiMemberPermissions = "/ResourceProfile/SetResourceProfileCollectionPermissions"
	case settype.CloudProvider.String():
		o.ValidMemberPermissions = ValidPermissionMap.Set
		o.apiMemberPermissions = "/CloudProvider/SetCloudProviderCollectionPermissions"
	default:
		return fmt.Errorf("Invalid Set type")
	}

	return nil
}

// GetIDByName returns set ID by name
func (o *ManualSet) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Set name must be provided")
	}
	if o.ObjectType == "" {
		return "", fmt.Errorf("Set type must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving set: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves set from tenant by name
func (o *ManualSet) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of set %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a Set by name
func (o *ManualSet) DeleteByName() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of Set %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

// ValidateMemberPermissions returns validated list of rights
func (o *ManualSet) ValidateMemberPermissions(perms []string) ([]string, error) {
	return ConvertToValidList(perms, o.ValidMemberPermissions)
}

/*
	API to manage set

	Get Set
	https://developer.centrify.com/reference#post_collection-getcollection

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "XXXXX",
				"_TableName": "collections",
				"_Timestamp": "/Date(1569301295080)/",
				"ObjectType": "Server",
				"ACL": "true",
				"Name": "LAB Systems",
				"Filters": ...
				"_PartitionKey": "XXXXX",
				"WhenCreated": "/Date(1569301294975)/",
				"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"_entitycontext": "W/\"datetime'2019-09-24T05%3A01%3A35.0809483Z'\"",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"ParentPath": "",
				"CollectionType": "ManualBucket",
				"MembersFile": "/sys/buckets/exxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx.json",
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
				}
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Create Set
	https://developer.centrify.com/reference#post_collection-createmanualcollection

		Request body format
		{
			"ObjectType": "Server",
			"Name": "Test Set 1",
			"Description": "Test Set 1",
			"CollectionType": "ManualBucket",
		}

		Respond result
		{
			"success": true,
			"Result": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update Set
	https://developer.centrify.com/reference#post_collection-updatecollection

		Request body format
		{
			"Rights": "View, Edit, Delete, Grant",
			"_TableName": "collections",
			"_encryptkeyid": "XXXXX",
			"ObjectType": "Server",
			"ACL": "true",
			"Name": "Test Set 1",
			"_PartitionKey": "XXXXX",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_entitycontext": "W/\"datetime'2020-08-01T09%3A49%3A55.7338109Z'\"",
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"CollectionType": "ManualBucket",
			"Description": "Test Set 2",
			"MembersFile": "/sys/buckets/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx.json",
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"sql": "undefined"
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

	Delete Set
	https://developer.centrify.com/reference#post_collection-deletecollection

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Updates set members
	https://developer.centrify.com/reference#post_collection-updatememberscollection
	Request body format
	{
		"id": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"add": [
			{
				"MemberType": "Row",
				"Table": "Server",
				"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			}
		]
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

	Set Collection Permissions
	https://developer.centrify.com/reference#post_collection-setcollectionpermissions

		Request body format
		{
			"Grants": [
				{
					"Principal": "LAB Infrastructure Admins",
					"PType": "Role",
					"Rights": "Grant,View,Edit,Delete",
					"PrincipalId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
				}
			],
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"PVID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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
