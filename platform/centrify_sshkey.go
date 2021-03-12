package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// SSHKey - Encapsulates a single generic SSHKey
type SSHKey struct {
	vaultObject
	apiGetChallenge string
	apiRetrieve     string

	SSHKeysDefaultProfileID string          `json:"SshKeysDefaultProfile,omitempty" schema:"default_profile_id,omitempty"` // Default SSH Key Challenge Profile
	ChallengeRules          *ChallengeRules `json:"SshKeysRules,omitempty" schema:"challenge_rule,omitempty"`
	KeyFormat               string          `json:"KeyFormat,omitempty" schema:"key_format,omitempty"`
	KeyLength               int             `json:"KeyLength,omitempty" schema:"key_length,omitempty"`
	KeyType                 string          `json:"KeyType,omitempty" schema:"key_type,omitempty"`
	IsManaged               bool            `json:"IsManaged,omitempty" schema:"is_managed,omitempty"`
	Description             string          `json:"Comment,omitempty" schema:"description,omitempty"`
	PrivateKey              string          `json:"PrivateKey,omitempty" schema:"private_key,omitempty"`
	Passphrase              string          `json:"Passphrase,omitempty" schema:"passphrase,omitempty"`
	KeyPairType             string          `json:"KeyPairType,omitempty" schema:"key_pair_type,omitempty"` // Which key to retrieve from the pair, must be either PublicKey, PrivateKey, or PPK
}

// NewSSHKey is a SSHKey constructor
func NewSSHKey(c *restapi.RestClient) *SSHKey {
	s := SSHKey{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.SSHKey
	s.SetType = settype.SSHKey.String()
	s.apiRead = "/ServerManage/GetSshKeyInfo"
	s.apiCreate = "/ServerManage/AddSshKey"
	s.apiDelete = "/ServerManage/DeleteSshKey"
	s.apiUpdate = "/ServerManage/UpdateSshKey"
	s.apiRetrieve = "/ServerManage/RetrieveSshKey"
	s.apiGetChallenge = "/ServerManage/GetSshKeyRightsAndChallenges"
	s.apiPermissions = "/ServerManage/SetSSHKeyPermissions"

	return &s
}

// Read function fetches a SSHKey from source
func (o *SSHKey) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	logger.Debugf("Response for SSHKey from tenant: %v", resp)
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

	// Get SSH Key challenge profile information
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
	if v, ok := resp.Result["Challenges"]; ok {
		challenges := v.(map[string]interface{})
		if p, ok := challenges["SshKeysDefaultProfile"]; ok {
			o.SSHKeysDefaultProfileID = p.(string)
		}
		// Fill challenge rules
		if r, ok := challenges["SshKeysRules"]; ok {
			challengerules := &ChallengeRules{}
			mapToStruct(challengerules, r.(map[string]interface{}))
			o.ChallengeRules = challengerules
		}
	}

	return nil
}

// Create function creates a new SSHKey and returns a map that contains creation result
func (o *SSHKey) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	// Special handling of challenge checkout profile
	queryArg["updateChallenges"] = true

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

// Delete function deletes a SSHKey and returns a map that contains deletion result
func (o *SSHKey) Delete() (*restapi.StringResponse, error) {
	return o.deleteObjectStringAPI("")
}

// Update function updates an existing SSHKey and returns a map that contains update result
func (o *SSHKey) Update() (*restapi.GenericMapResponse, error) {
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
	// Special handling of password checkout profile
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

// Query function returns a single SSHKey object in map format
func (o *SSHKey) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM SshKeys WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// RetriveSSHKey retrieves SSH Key from vault
func (o *SSHKey) RetriveSSHKey() (string, error) {
	if o.ID == "" && o.Name == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return "", fmt.Errorf(errormsg)
	}
	// If SSHKey name is provided, try to find out its ID
	if o.Name != "" {
		var err error
		o.ID, err = o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return "", err
		}
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	if o.KeyPairType == "" {
		return "", fmt.Errorf("Missing KeyPairType. It must be PublicKey, PrivateKey, or PPK")
	}
	queryArg["KeyPairType"] = o.KeyPairType
	if o.Passphrase != "" {
		queryArg["Passphrase"] = o.Passphrase
	}
	if o.KeyFormat != "" {
		queryArg["KeyFormat"] = o.KeyFormat
	}

	resp, err := o.client.CallStringAPI(o.apiRetrieve, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}

	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return "", fmt.Errorf(errmsg)
	}

	return resp.Result, nil
}

// GetIDByName returns SSHKey ID by name
func (o *SSHKey) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("SSHKey name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving SSHKey: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves sshkey from tenant by name
func (o *SSHKey) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of sshkey %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a sshkey by name
func (o *SSHKey) DeleteByName() (*restapi.StringResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of sshkey %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
	Fetch SSH Key
	https://developer.centrify.com/reference#post_servermanage-getsshkeyinfo

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond Result
		{
			"success": true,
			"Result": {
				"_entitycontext": "W/\"datetime'2020-08-24T12%3A16%3A42.6342428Z'\"",
				"KeyFormat": "PEM",
				"LastUpdated": "/Date(1598271402593)/",
				"IsManaged": false,
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
				},
				"_TableName": "sshkeys",
				"_encryptkeyid": "XXXXX",
				"_PartitionKey": "XXXXX",
				"_RowKey": "5xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Revision": 0,
				"_Timestamp": "/Date(1582552067897)/",
				"CreatedBy": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Name": "my_ami_key",
				"KeyType": "RSA",
				"Comment": "my AWS AMI key",
				"Created": "/Date(1582552067684)/",
				"ACL": "true",
				"KeyLength": 2048,
				"State": "Active",
				"LastUpdatedBy": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Add SSH Key
		https://developer.centrify.com/reference#post_servermanage-addsshkey

		Request body format
		{
			"Name": "Test Key",
			"Comment": "Test Key",
			"jsutil-enhancecheckbox-34778-inputEl": false,
			"PrivateKey": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCA...AwOQ==\n-----END RSA PRIVATE KEY-----",
			"Type": "Manual"
		}

		Respond Result
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

	Update SSH Key
	https://developer.centrify.com/reference#post_servermanage-updatesshkey

		Request body format
		{
			"IsManaged": false,
			"IsFavorite": false,
			"Revision": 0,
			"CreatedBy": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Test Key",
			"KeyType": "RSA",
			"Comment": "Test Key 1",
			"KeyLength": 2048,
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"State": "Active",
			"jsutil-enhancecheckbox-35013-inputEl": false
		}

		Respond Result
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

	Delete SSH Key
	https://developer.centrify.com/reference#post_servermanage-deletesshkey

		Reqeust body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond Result
		{
			"success": true,
			"Result": "Successfully deleted SSH Key Test Key",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Retrieve SSH Key
	https://developer.centrify.com/reference#post_servermanage-retrievesshkey

	Reqeust body format
		{
			"KeyPairType": "PrivateKey",
			"KeyFormat": "PEM",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond Result
		{
			"success": true,
			"Result": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCA...AwOQ==\n-----END RSA PRIVATE KEY-----",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}
*/
