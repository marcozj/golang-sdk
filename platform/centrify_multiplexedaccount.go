package platform

import (
	"fmt"
	"strings"

	"github.com/marcozj/golang-sdk/enum/resourcetype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// MultiplexedAccount - Encapsulates a single MultiplexedAccount
type MultiplexedAccount struct {
	vaultObject

	RealAccount1ID  string   `json:"RealAccount1ID,omitempty" schema:"account1_id,omitempty"`
	RealAccount2ID  string   `json:"RealAccount2ID,omitempty" schema:"account2_id,omitempty"`
	RealAccount1    string   `json:"RealAccount1,omitempty" schema:"account1,omitempty"`
	RealAccount2    string   `json:"RealAccount2,omitempty" schema:"account2,omitempty"`
	ActiveAccount   string   `json:"ActiveAccount,omitempty" schema:"active_account,omitempty"`
	RealAccounts    []string `json:"RealAccounts,omitempty" schema:"accounts,omitempty"`
	RealAccount1UPN string   `json:"-"`
	RealAccount2UPN string   `json:"-"`
}

// NewMultiplexedAccount is a MultiplexedAccount constructor
func NewMultiplexedAccount(c *restapi.RestClient) *MultiplexedAccount {
	s := MultiplexedAccount{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.MultiplexAccount
	s.apiRead = "/Subscriptions/GetMultiplexedAccount"
	s.apiCreate = "/Subscriptions/CreateMPAccount"
	s.apiDelete = "/Subscriptions/DeleteMPAccount"
	s.apiUpdate = "/Subscriptions/UpdateMPAccount"
	s.apiPermissions = "/Subscriptions/SetMultiplexedAccountPermissions"

	return &s
}

// Read function fetches a MultiplexedAccount from source, including attribute values. Returns error if any
func (o *MultiplexedAccount) Read() error {
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
	// fill RealAccounts
	if o.RealAccount1ID != "" && o.RealAccount2ID != "" && o.RealAccounts == nil {
		o.RealAccounts = []string{o.RealAccount1ID, o.RealAccount2ID}
	}

	return nil
}

// Create function creates a new MultiplexedAccount
func (o *MultiplexedAccount) Create() (*restapi.StringResponse, error) {
	err := o.resolveAccountID()
	if err != nil {
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
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

// Delete function deletes a MultiplexedAccount
func (o *MultiplexedAccount) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing MultiplexedAccount
func (o *MultiplexedAccount) Update() (*restapi.StringResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	err := o.resolveAccountID()
	if err != nil {
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
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

// Query function returns a single MultiplexedAccount object in map format
func (o *MultiplexedAccount) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM MultiplexedAccount WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns MultiplexedAccount ID by name
func (o *MultiplexedAccount) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("MultiplexedAccount name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving MultiplexedAccount: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves MultiplexedAccount from tenant by name
func (o *MultiplexedAccount) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of MultiplexedAccount %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a MultiplexedAccount by name
func (o *MultiplexedAccount) DeleteByName() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of MultiplexedAccount %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *MultiplexedAccount) resolveAccountID() error {
	// Make sure if either both RealAccount1UPN & RealAccount2UPN are set or both of them are not set at all
	if o.RealAccount1UPN != "" && o.RealAccount2UPN != "" {
		o.RealAccounts = nil
		accts := []string{o.RealAccount1UPN, o.RealAccount2UPN}
		for _, acct := range accts {
			// Breaks account if it is upn <username>@<domain>
			acctparts := strings.Split(acct, "@")
			var acctname, acctdomain string
			acctname = acctparts[0]
			if len(acctparts) > 1 {
				acctdomain = acctparts[1]
			} else {
				return fmt.Errorf("RealAccountxUPN must be in <username>@<domain> format. But it is '%s'", acct)
			}

			account := NewAccount(o.client)
			account.User = acctname
			//var resourceID string
			var err error
			account.ResourceType = resourcetype.Domain.String()
			account.ResourceName = acctdomain

			acctid, err := account.GetIDByName()
			if err != nil {
				logger.Errorf(err.Error())
				return fmt.Errorf(err.Error())
			}
			o.RealAccounts = append(o.RealAccounts, acctid)
		}
	} else if (o.RealAccount1UPN != "" && o.RealAccount2UPN == "") || (o.RealAccount1UPN == "" && o.RealAccount2UPN != "") {
		return fmt.Errorf("Only 1 account is set")
	}

	return nil
}

/*
Get Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-getmultiplexedaccount

	Request body
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

	Responde Result
	{
		"success": true,
		"Result": {
			"RealAccount2ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "My Multiplexed Account",
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"ActiveAccount": "csvc_acct1 (exampoe.com)",
			"RealAccount1ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Issues": "",
			"Description": "My Multiplexed Account -",
			"RealAccount2": "csvc_acct2 (example.com)",
			"RealAccount1": "csvc_acct1 (example.com)"
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Create Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-deletempaccount

	Request body
	{
		"RealAccount1ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"RealAccount2ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Name": "My Multiplexed Account",
		"Description": "My Multiplexed Account",
		"RealAccount1": "csvc_acct1 (example.com)",
		"RealAccount2": "csvc_acct2 (example.com)",
		"RealAccounts": [
			"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		]
	}

	Responde Result
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

Update Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-updatempaccount

	Request body
	{
		"RealAccount2ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Name": "My Multiplexed Account",
		"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"ActiveAccount": "csvc_acct1 (example.com)",
		"RealAccount1ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Description": "My Multiplexed Account -",
		"RealAccount2": "csvc_acct2 (example.com)",
		"RealAccount1": "csvc_acct1 (example.com)",
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"RealAccounts": [
			"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		]
	}

	Responde Result
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

Delete Multiplexed Account
https://developer.centrify.com/reference#post_subscriptions-updatempaccount

	Request body
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Responde Result
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
