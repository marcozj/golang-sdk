package platform

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Secret - Encapsulates a single generic secret
type Secret struct {
	vaultObject
	// VaultData specific APIs
	apiRetrieveSecret             string
	apiMoveSecret                 string
	apiGetChallenge               string
	apiRequestSecretDownloadUrl   string
	apiDownloadSecretFileInChunks string

	SecretName              string          `json:"SecretName,omitempty" schema:"secret_name,omitempty"` // User Name
	SecretText              string          `json:"SecretText,omitempty" schema:"secret_text,omitempty"`
	Type                    string          `json:"Type,omitempty" schema:"type,omitempty"`
	FolderID                string          `json:"FolderId,omitempty" schema:"folder_id,omitempty"`
	ParentPath              string          `json:"ParentPath,omitempty" schema:"parent_path,omitempty"`
	DataVaultDefaultProfile string          `json:"DataVaultDefaultProfile" schema:"default_profile_id"` // Default Secret Challenge Profile (used if no conditions matched)
	ChallengeRules          *ChallengeRules `json:"DataVaultRules,omitempty" schema:"challenge_rule,omitempty"`
	Sets                    []string        `json:"Sets,omitempty" schema:"sets,omitempty"`
	NewParentPath           string          `json:"-"`
	SecretFileName          string          `json:"SecretFileName,omitempty" schema:"secret_filename,omitempty"`
	WorkflowEnabled         bool            `json:"WorkflowEnabled,omitempty" schema:"workflow_enabled,omitempty"`
	//WorkflowSent         bool               `json:"WorkflowSent,omitempty" schema:"workflow_sent,omitempty"`
	WorkflowApprovers      []WorkflowApprover      `json:"WorkflowApprovers,omitempty" schema:"workflow_approver,omitempty"`
	WorkflowDefaultOptions *WorkflowDefaultOptions `json:"WorkflowDefaultOptions,omitempty" schema:"workflow_default_options,omitempty"`
}

// NewSecret is a Secret constructor
func NewSecret(c *restapi.RestClient) *Secret {
	s := Secret{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Secret
	s.SetType = settype.Secret.String()
	s.apiRead = "/ServerManage/GetSecret"
	s.apiCreate = "/ServerManage/AddSecret"
	s.apiDelete = "/ServerManage/DeleteSecret"
	s.apiUpdate = "/ServerManage/UpdateSecret"
	s.apiRetrieveSecret = "/ServerManage/RetrieveSecretContents"
	s.apiMoveSecret = "/ServerManage/MoveSecret"
	s.apiPermissions = "/ServerManage/SetSecretPermissions"
	s.apiGetChallenge = "/ServerManage/GetSecretRightsAndChallenges"
	s.apiRequestSecretDownloadUrl = "ServerManage/RequestSecretDownloadUrl"
	s.apiDownloadSecretFileInChunks = "ServerManage/DownloadSecretFileInChunks"

	return &s
}

// Read function fetches a Secret from source, including attribute values. Returns error if any
func (o *Secret) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	logger.Debugf("Response for Secret from tenant: %v", resp)
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

	// Get challenge profile information
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
	if v, ok := resp.Result["DataVaultDefaultProfile"]; ok {
		o.DataVaultDefaultProfile = v.(string)
	}

	// Fill challenge rules
	if v, ok := resp.Result["Challenges"]; ok {
		challenges := v.(map[string]interface{})
		if challenges["DataVaultDefaultProfile"] != nil {
			o.DataVaultDefaultProfile = challenges["DataVaultDefaultProfile"].(string)
		}
		if r, ok := challenges["DataVaultRules"]; ok {
			challengerules := &ChallengeRules{}
			mapToStruct(challengerules, r.(map[string]interface{}))
			o.ChallengeRules = challengerules
		}
	}

	return nil
}

// Create function creates a new Secret and returns a map that contains creation result
func (o *Secret) Create() (*restapi.StringResponse, error) {
	// Resolve FolderID if only ParentPath is provided
	err := o.resolveFolderdID()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	err = o.processWorkflow()
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
	queryArg["updateChallenges"] = false
	if o.WorkflowEnabled {
		queryArg["WorkflowSent"] = true
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

// Delete function deletes a Secret and returns a map that contains deletion result
func (o *Secret) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing Secret and returns a map that contains update result
func (o *Secret) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	// Resolve FolderID if only ParentPath is provided or NewParentPath is provided for moving into another folder
	err := o.resolveFolderdID()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	err = o.processWorkflow()
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
	// Need to always send this when workflow is turned on and off
	queryArg["WorkflowSent"] = true

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

// MoveSecret function moves an existing Secret to another folder
func (o *Secret) MoveSecret() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["targetFolderId"] = o.FolderID
	//queryArg["updateChallenges"] = true

	logger.Debugf("Generated Map for MoveFolder(): %+v", queryArg)

	resp, err := o.client.CallBoolAPI(o.apiMoveSecret, queryArg)
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

// Query function returns a single Secret object in map format
func (o *Secret) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM DataVault WHERE 1=1"
	if o.SecretName != "" {
		query += " AND SecretName='" + o.SecretName + "'"
	}
	// ParentPath should always be added
	query += " AND ParentPath='" + o.ParentPath + "'"

	return queryVaultObject(o.client, query)
}

func (o *Secret) checkoutSecret() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Description"] = "Checkout by golang SDK"

	resp, err := o.client.CallGenericMapAPI(o.apiRetrieveSecret, queryArg)
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

// CheckoutSecret checks out secret from vault
func (o *Secret) CheckoutSecret() (string, error) {
	// To retrieve secret, we must know its ID
	// In order to know ID, we must know SecretName + ParentPath
	if o.ID == "" {
		err := o.GetByName()
		if err != nil {
			logger.Errorf(err.Error())
			return "", fmt.Errorf("Failed to find secret %s. %v", o.SecretName, err)
		}
	}

	// Check again if ID is known
	if o.ID == "" {
		return "", fmt.Errorf("Missing ID for secret %s in %s", o.SecretName, o.ParentPath)
	}

	resp, err := o.checkoutSecret()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving secret content for %s: %s", o.SecretName, err)
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return "", fmt.Errorf(errmsg)
	}
	if p, ok := resp.Result["SecretText"]; ok {
		return p.(string), nil
	}

	return "", fmt.Errorf("Failed to retrieve secret %s", o.SecretName)
}

// GetIDByName returns Secret ID by name
func (o *Secret) GetIDByName() (string, error) {
	if o.SecretName == "" {
		return "", fmt.Errorf("Secret name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving secret: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves Secret from tenant by name
func (o *Secret) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of secret %s. %v", o.SecretName, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a Secret by name
func (o *Secret) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of secret %s. %v", o.SecretName, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *Secret) resolveFolderdID() error {
	// If NewParentPath is set, this is called from directly API
	// It means we want to change folder so need to recaculate FolderID
	if o.NewParentPath != "" {
		o.ParentPath = o.NewParentPath
		o.FolderID = ""
	}

	if o.FolderID == "" && o.ParentPath != "" {
		path := strings.Split(o.ParentPath, "\\")
		folder := NewSecretFolder(o.client)
		// folder name is the last in split slice
		folder.Name = path[len(path)-1]
		if len(path) > 1 {
			folder.ParentPath = strings.TrimSuffix(o.ParentPath, fmt.Sprintf("\\%s", path[len(path)-1]))
		}
		var err error
		o.FolderID, err = folder.GetIDByName()
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Secret) DownloadSecretFile(saveToHome bool) (string, error) {
	if o.ID == "" {
		err := o.GetByName()
		if err != nil {
			logger.Errorf(err.Error())
			return "", fmt.Errorf("Failed to find secret %s. %v", o.SecretName, err)
		}
	}

	if o.Type != "File" || o.SecretFileName == "" {
		return "", fmt.Errorf("This secret type is '%s' and SecretFileName is missing.\n", o.Type)
	}

	var queryArg = make(map[string]interface{})
	queryArg["secretID"] = o.ID

	resp, err := o.client.CallGenericMapAPI(o.apiRequestSecretDownloadUrl, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return "", fmt.Errorf(errmsg)
	}

	secretfilepath := resp.Result["FilePath"].(string)
	savedfilepath := o.SecretFileName
	if saveToHome {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		savedfilepath = user.HomeDir + "/" + savedfilepath
	}
	var downloadArg = make(map[string]interface{})
	o.client.DownloadFile(o.apiDownloadSecretFileInChunks+"?FilePath="+secretfilepath, downloadArg, savedfilepath)
	if err != nil {
		return "", err
	}

	return savedfilepath, nil
}

// CheckoutSecretAndFile checks out secret from vault and supports file type secret
func (o *Secret) CheckoutSecretAndFile(saveToHome bool) (string, error) {
	if o.ID == "" {
		err := o.GetByName()
		if err != nil {
			logger.Errorf(err.Error())
			return "", fmt.Errorf("Failed to find secret %s. %v", o.SecretName, err)
		}
	}

	// Check again if ID is known
	if o.ID == "" {
		return "", fmt.Errorf("Missing ID for secret %s in %s", o.SecretName, o.ParentPath)
	}

	// if it is text secret, retrieve secret text
	if o.Type == "Text" {
		resp, err := o.checkoutSecret()
		if err != nil {
			logger.Errorf(err.Error())
			return "", fmt.Errorf("Error retrieving secret content for %s: %s", o.SecretName, err)
		}
		if !resp.Success {
			errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
			logger.Errorf(errmsg)
			return "", fmt.Errorf(errmsg)
		}
		if p, ok := resp.Result["SecretText"]; ok {
			return p.(string), nil
		}
	} else if o.Type == "File" {
		filename, err := o.DownloadSecretFile(saveToHome)
		if err != nil {
			return "", err
		}
		// For File type, return file name
		return filename, nil
	}

	return "", fmt.Errorf("Failed to retrieve secret %s", o.SecretName)
}

func (o *Secret) processWorkflow() error {
	// Resolve guid of each approver
	if o.WorkflowEnabled && o.WorkflowApprovers != nil {
		err := ResolveWorkflowApprovers(o.client, o.WorkflowApprovers)
		if err != nil {
			return err
		}
		if o.WorkflowDefaultOptions == nil {
			o.WorkflowDefaultOptions = &WorkflowDefaultOptions{
				GrantMin: 60,
			}
		}
	}

	return nil
}

/*
	API to manage vault secret

	Read Secret
	https://developer.centrify.com/reference#post_servermanage-getsecret

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RRFormat": true,
			"Args": {
				"PageNumber": 1,
				"Limit": 1,
				"PageSize": 1,
				"Caching": -1
			}
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "XXXXXX",
				"_TableName": "datavault",
				"_Timestamp": "/Date(1584413116338)/",
				"WhenContentsReplaced": "/Date(1584413116309)/",
				"ACL": "true",
				"_PartitionKey": "XXXXXX",
				"WhenCreated": "/Date(1582558666855)/",
				"_entitycontext": "W/\"datetime'2020-03-17T02%3A45%3A16.3380444Z'\"",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"WhenUpdated": "/Date(1584413116309)/",
				"ParentPath": "Folder 1\\Folder level 2",
        		"FolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Description": "admin@example.com",
				"SecretName": "Centrify PAS Admin Credential",
				"Type": "Text",
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

	Add Secret
	https://developer.centrify.com/reference#post_servermanage-addsecret

		Request body format
		{
			"SecretName": "Access key",
			"Description": "AWS access key",
			"SecretText": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Type": "Text",
			"SetID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": false
		}
		or
		{
			"FolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"SecretName": "Another secret",
			"Description": "Another secret",
			"SecretText": "xxxxxxxxxxxxx",
			"Type": "Text",
			"updateChallenges": false
		}
		or
		{
			"FolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxc",
			"SecretName": "File1",
			"SecretFilePath": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"SecretFileSize": "38.003 KB",
			"SecretFilePassword": "xxxxxxx",
			"Type": "File",
			"Description": "",
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

	Update Secret
	https://developer.centrify.com/reference#post_servermanage-updatesecret

		Request body format
		{
			"SecretName": "Access key",
			"Description": "AWS access key",
			"SecretText": "xxxxxxxxxxxxx",
			"Type": "Text",
			"SetID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": true,
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"DataVaultDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}
		or
		{
			"SecretText": "xxxxxxxxxxxxx",
			"SecretName": "Random secret",
			"Type": "Text",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Delete Secret
	https://developer.centrify.com/reference#post_servermanage-deletesecret

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
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

	Retrieve Secret content
	https://developer.centrify.com/reference#post_servermanage-retrievesecretcontents

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_encryptkeyid": "XXXXXX",
				"_TableName": "datavault",
				"_Timestamp": "/Date(1592380339832)/",
				"ACL": "true",
				"_PartitionKey": "XXXXXX",
				"WhenCreated": "/Date(1592380339057)/",
				"_entitycontext": "W/\"datetime'2020-06-17T07%3A52%3A19.8321511Z'\"",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"ParentPath": "",
				"Description": "A random secret",
				"SecretName": "Randon secret",
				"Type": "Text",
				"SecretText": "xxxxxxxxxxx",
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

	Move Secret to another folder
	https://developer.centrify.com/reference#post_servermanage-movesecret

	Request body format
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"targetFolderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
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
