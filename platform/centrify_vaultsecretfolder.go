package platform

import (
	"fmt"
	"strings"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// SecretFolder - Encapsulates a single generic secret folder
type SecretFolder struct {
	vaultObject
	// SecretFolder specific APIs
	apiGetChallenge        string
	apiMoveFolder          string
	apiMemberPermissions   string
	ValidMemberPermissions map[string]string

	Type                            string          `json:"Type,omitempty" schema:"type,omitempty"`        // Can only be Folder
	ParentID                        string          `json:"Parent,omitempty" schema:"parent_id,omitempty"` // ID of parent folder
	ParentPath                      string          `json:"ParentPath,omitempty" schema:"parent_path,omitempty"`
	CollectionMembersDefaultProfile string          `json:"CollectionMembersDefaultProfile" schema:"default_profile_id"` // Default Secret Challenge Profile (used if no conditions matched)
	ChallengeRules                  *ChallengeRules `json:"CollectionMembersRules,omitempty" schema:"challenge_rule,omitempty"`
	MemberPermissions               []Permission
	NewParentPath                   string `json:"-"`
}

// NewSecretFolder is a SecretFolder constructor
func NewSecretFolder(c *restapi.RestClient) *SecretFolder {
	s := SecretFolder{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Folder
	s.ValidMemberPermissions = ValidPermissionMap.Secret
	s.SetType = settype.Secret.String()
	s.apiRead = "/ServerManage/GetSecretFolder"
	s.apiCreate = "/ServerManage/AddSecretsFolder"
	s.apiDelete = "/ServerManage/DeleteSecretsFolder"
	s.apiUpdate = "/ServerManage/UpdateSecretsFolder"
	s.apiGetChallenge = "/ServerManage/GetSecretsFolderRightsAndChallenges"
	s.apiMoveFolder = "/ServerManage/MoveFolder"
	s.apiPermissions = "/ServerManage/SetSecretsFolderPermissions"
	s.apiMemberPermissions = "/ServerManage/SetSecretCollectionPermissions"
	s.Type = "Folder"

	return &s
}

// Read function fetches a SecretFolder from source, including attribute values. Returns error if any
func (o *SecretFolder) Read() error {
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

	// Loop through respond results and grab the first record
	var results = resp.Result["Results"].([]interface{})
	if len(results) < 1 {
		// Make sure error message contains "not exist"
		logger.Debugf("Returning SecretFolder does not exist in tenant")
		return fmt.Errorf("SecretFolder does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return fmt.Errorf("There are more than one SecretFolder with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})

	mapToStruct(o, row)

	// Get challenge profile
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}
	//logger.Debugf("Challenges result: %+v", resp)
	if v, ok := resp.Result["Challenges"]; ok {
		challenges := v.(map[string]interface{})
		if challenges["CollectionMembersDefaultProfile"] != nil {
			o.CollectionMembersDefaultProfile = challenges["CollectionMembersDefaultProfile"].(string)
		}
		// Fill challenge rules
		if r, ok := challenges["CollectionMembersRules"]; ok {
			challengerules := &ChallengeRules{}
			mapToStruct(challengerules, r.(map[string]interface{}))
			o.ChallengeRules = challengerules
		}
	}

	return nil
}

// Create function creates a new SecretFolder and returns a map that contains creation result
func (o *SecretFolder) Create() (*restapi.StringResponse, error) {
	// Resolve ParentID if only ParentPath is provided
	err := o.resolveParentID()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
	if err != nil {
		return nil, err
	}
	queryArg["updateChallenges"] = false

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

// Delete function deletes a SecretFolder and returns a map that contains deletion result
func (o *SecretFolder) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing SecretFolder and returns a map that contains update result
func (o *SecretFolder) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	// Resolve ParentID if only ParentPath is provided or NewParentPath is provided for moving into another folder
	err := o.resolveParentID()
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
	queryArg["updateChallenges"] = true

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

// MoveFolder function moves an existing SecretFolder to another folder
func (o *SecretFolder) MoveFolder() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	// Resolve ParentID if only ParentPath is provided or NewParentPath is provided for moving into another folder
	err := o.resolveParentID()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["targetFolderId"] = o.ParentID
	//queryArg["updateChallenges"] = true

	logger.Debugf("Generated Map for MoveFolder(): %+v", queryArg)

	resp, err := o.client.CallBoolAPI(o.apiMoveFolder, queryArg)
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

// Query function returns a single SecretFolder object in map format
func (o *SecretFolder) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Sets WHERE 1=1 AND ObjectType='DataVault' AND CollectionType='Phantom'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.ParentPath != "" {
		query += " AND ParentPath='" + o.ParentPath + "'"
	}

	return queryVaultObject(o.client, query)
}

// SetMemberPermissions sets member permissions. isRemove indicates whether to remove all permissions instead of setting permissions
func (o *SecretFolder) SetMemberPermissions(isRemove bool) (*restapi.GenericMapResponse, error) {
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

// GetIDByName returns Secret folder ID by name
func (o *SecretFolder) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Secret folder name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		errormsg := fmt.Sprintf("Failed to retrieve secret folder '%s' in '%s'. %v", o.Name, o.ParentPath, err)
		logger.Errorf(errormsg)
		return "", fmt.Errorf(errormsg)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves Secret folder from tenant by name
func (o *SecretFolder) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf(err.Error())
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a Secret folder by name
func (o *SecretFolder) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf(err.Error())
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *SecretFolder) resolveParentID() error {
	// If NewParentPath is set, this is called from directly API
	// It means we want to change folder so need to recaculate FolderID
	if o.NewParentPath != "" {
		o.ParentPath = o.NewParentPath
		o.ParentID = ""
	}

	if o.ParentID == "" && o.ParentPath != "" {
		path := strings.Split(o.ParentPath, "\\")
		folder := NewSecretFolder(o.client)
		// folder name is the last in split slice
		folder.Name = path[len(path)-1]
		if len(path) > 1 {
			folder.ParentPath = strings.TrimSuffix(o.ParentPath, fmt.Sprintf("\\%s", path[len(path)-1]))
		}
		var err error
		o.ParentID, err = folder.GetIDByName()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	API to manage vault secret

	Read folder
	https://developer.centrify.com/reference#post_servermanage-getsecretfolder

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RRFormat": true
		}

		Respond result
		{
			"success": true,
			"Result": {
				"IsAggregate": false,
				"Count": 1,
				"Columns": [...],
				"FullCount": 1,
				"Results": [
					{
						"Entities": [
							{
								"Type": "DataVaultFolder",
								"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
								"IsForeignKey": false
							}
						],
						"Row": {
							"_entitycontext": "W/\"datetime'2020-08-16T02%3A08%3A49.301252Z'\"",
							"Name": "Test Folder",
							"_metadata": {
								"Version": 1,
								"IndexingVersion": 1
							},
							"ACL": "true",
							"_PartitionKey": "XXXXXX",
							"WhenCreated": "/Date(1597543728366)/",
							"_Timestamp": "/Date(1597543728399)/",
							"_encryptkeyid": "XXXXXX",
							"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"Description": "Test Folder",
							"Type": "Folder",
							"ObjectType": "DataVault",
							"SecretName": "Test Folder",
							"CollectionType": "Phantom",
							"_TableName": "collections",
							"Filters": "FolderId",
							"ParentPath": ""
						}
					}
				],
				"ReturnID": ""
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Add folder
	https://developer.centrify.com/reference#post_servermanage-addsecretsfolder

		Request body format
		{
			"Name": "Test Folder",
			"Description": "Test Folder",
			"Type": "Folder",
			"updateChallenges": false
		}
		or
		{
			"Parent": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Sub folder",
			"Description": "Sub folder",
			"Type": "Folder",
			"updateChallenges": false
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

	Update folder
	https://developer.centrify.com/reference#post_servermanage-updatesecretsfolder

		Request body format
		{
			"_entitycontext": "W/\"datetime'2020-08-16T02%3A17%3A58.6860435Z'\"",
			"Name": "Sub folder",
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"ACL": "true",
			"_PartitionKey": "XXXXXX",
			"_encryptkeyid": "XXXXXX",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": "Sub folder of something",
			"Type": "Folder",
			"ObjectType": "DataVault",
			"SecretName": "Sub folder",
			"CollectionType": "Phantom",
			"_TableName": "collections",
			"Filters": "FolderId",
			"Parent": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"ParentPath": "Folder 1",
			"Rights": "View, Edit, Delete, Grant, Add",
			"CollectionMembersDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": true
		}

		Respond result
		{
			"success": true,
			"Result": {
				"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete folder
	https://developer.centrify.com/reference#post_servermanage-deletesecretsfolder

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": false,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}
		or
		{
			"success": false,
			"Result": {
				"ChallengeId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": "You must respond to a challenge to proceed.",
			"MessageID": null,
			"Exception": null,
			"ErrorID": "ChallengeRequired",
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Get Rights and Challenges

	Request body format
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"Rights": "View, Edit, Delete, Grant, Add",
			"Challenges": {
				"CollectionMembersRules": {
					"_UniqueKey": "Condition",
					"_Value": [],
					"Enabled": true,
					"_Type": "RowSet"
				},
				"CollectionMembersDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Move folder to another folder
	https://developer.centrify.com/reference#post_servermanage-movefolder

	Request body format
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"targetFolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": false,
		"Result": {
			"ChallengeId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		},
		"Message": "You must respond to a challenge to proceed.",
		"MessageID": null,
		"Exception": null,
		"ErrorID": "ChallengeRequired",
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
	or
	{
		"success": true,
		"Result": true,
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
*/
