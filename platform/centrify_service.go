package platform

import (
	"fmt"
	"strings"

	"github.com/marcozj/golang-sdk/enum/computerclass"
	"github.com/marcozj/golang-sdk/enum/resourcetype"
	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Service - Encapsulates a single Service
type Service struct {
	vaultObject

	SystemID               string `json:"ComputerID,omitempty" schema:"system_id,omitempty"`
	SystemName             string `json:"-"` // Use by SDK call
	ServiceType            string `json:"Type,omitempty" schema:"service_type,omitempty"`
	Name                   string `json:"WindowsServiceName,omitempty" schema:"service_name,omitempty"`
	EnableManagement       bool   `json:"IsActive" schema:"enable_management"`
	AdminAccountID         string `json:"PushCreds,omitempty" schema:"admin_account_id,omitempty"`
	AdminAccountUPN        string `json:"-"` // Use by SDK call
	MultiplexedAccountID   string `json:"AccountID,omitempty" schema:"multiplexed_account_id,omitempty"`
	MultiplexedAccountName string `json:"-"` // Use by SDK call
	RestartService         bool   `json:"RestartService" schema:"restart_service"`
	RestartTimeRestriction bool   `json:"RestartTimeRestriction" schema:"restart_time_restriction"`
	DaysOfWeek             string `json:"DaysOfWeek,omitempty" schema:"days_of_week,omitempty"`
	RestartStartTime       string `json:"RestartStartTime,omitempty" schema:"restart_start_time,omitempty"`
	RestartEndTime         string `json:"RestartEndTime,omitempty" schema:"restart_end_time,omitempty"`
	UseUTCTime             bool   `json:"RestartTimeInUtc" schema:"use_utc_time"`
}

// NewService is a Service constructor
func NewService(c *restapi.RestClient) *Service {
	s := Service{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Service
	s.SetType = settype.Service.String()
	s.apiRead = "/Subscriptions/GetSubscription"
	s.apiCreate = "/Subscriptions/AddSubscription"
	s.apiDelete = "/Subscriptions/DeleteSubscription"
	s.apiUpdate = "/Subscriptions/UpdateSubscription"
	s.apiPermissions = "/Subscriptions/SetSubscriptionPermissions"

	return &s
}

// Read function fetches a Service from source, including attribute values. Returns error if any
func (o *Service) Read() error {
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

// Create function creates a new Service
func (o *Service) Create() (*restapi.StringResponse, error) {
	err := o.resolveIDs()
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

// Delete function deletes a Service
func (o *Service) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing Service
func (o *Service) Update() (*restapi.StringResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	err := o.resolveIDs()
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

// Query function returns a single Service object in map format
func (o *Service) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Subscriptions WHERE 1=1"
	if o.Name != "" {
		query += " AND WindowsServiceName='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns service ID by name
func (o *Service) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Service name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving service: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves service from tenant by name
func (o *Service) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of service %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a service by name
func (o *Service) DeleteByName() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of service %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *Service) resolveIDs() error {
	var err error

	// Resolve SystemID
	if o.SystemName != "" {
		system := NewSystem(o.client)
		system.Name = o.SystemName
		system.ComputerClass = computerclass.Windows.String()
		o.SystemID, err = system.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf(err.Error())
		}
	}

	if o.EnableManagement {
		// Resolve AdminAccountID
		if o.AdminAccountUPN != "" {
			// Breaks account if it is upn <username>@<domain>
			acctparts := strings.Split(o.AdminAccountUPN, "@")
			var acctname, acctdomain string
			acctname = acctparts[0]
			if len(acctparts) > 1 {
				acctdomain = acctparts[1]
			} else {
				return fmt.Errorf("AdminAccountUPN must be in <username>@<domain> format. But it is '%s'", o.AdminAccountUPN)
			}

			account := NewAccount(o.client)
			account.User = acctname
			//var resourceID string
			var err error
			account.ResourceType = resourcetype.Domain.String()
			account.ResourceName = acctdomain

			o.AdminAccountID, err = account.GetIDByName()
			if err != nil {
				logger.Errorf(err.Error())
				return fmt.Errorf(err.Error())
			}
		}

		// Resolve MultiplexedAccountID
		if o.MultiplexedAccountName != "" {
			mplex := NewMultiplexedAccount(o.client)
			mplex.Name = o.MultiplexedAccountName
			o.MultiplexedAccountID, err = mplex.GetIDByName()
			if err != nil {
				logger.Errorf(err.Error())
				return fmt.Errorf(err.Error())
			}
		}
	}

	return nil
}

/*
Get Service
https://developer.centrify.com/reference#post_subscriptions-getsubscription

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
			"ComputerID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": "",
			"WindowsServiceName": "TestWindowsService",
			"SecretID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_STAMP": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RestartService": true,
			"LogonCreds": "ad_admin (example.com)",
			"_TableName": "cpssubscriptions",
			"_PartitionKey": "XXXXXX",
			"IsActive": true,
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RestartTimeRestriction": false,
			"Type": "WindowsService",
			"Account": "Account for TestWindowsService",
			"PushCreds": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"LastUpdate": "/Date(1595301043898)/",
			"AccountID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Status": "OK",
			"Resource": "member1.example.com",
			"ACL": "true",
			"Mode": "Push",
			"CurrentMultiplexAccount": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RestartTimeInUtc": true
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Create Service
https://developer.centrify.com/reference#post_subscriptions-addsubscription

	Request body
	{
		"RestartTimeInUtc": false,
		"ComputerID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"AccountID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"PushCreds": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Resource": "member2",
		"Description": "TestWinService",
		"Type": "WindowsService",
		"WindowsServiceName": "testwinsrv",
		"IsActive": true,
		"LogonCreds": "ad_admin (example.com)",
		"Account": "test",
		"RestartService": true,
		"RestartTimeRestriction": true,
		"DaysOfWeek": "Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday",
		"RestartStartTime": "10:00",
		"RestartEndTime": "09:00",
		"Mode": "Push"
	}

	Responde Result
	{
		"success":true,
		"Result":"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Message":null,
		"MessageID":null,
		"Exception":null,
		"ErrorID":null,
		"ErrorCode":null,"
		IsSoftError":false,
		"InnerExceptions":null
	}

Update Service
https://developer.centrify.com/reference#post_subscriptions-updatesubscription

	Request body
	{
		"ComputerID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Description": "Test Windows Service",
		"WindowsServiceName": "TestWindowsService",
		"SecretID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"_STAMP": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"RestartService": true,
		"LogonCreds": "ad_admin (example.com)",
		"_TableName": "cpssubscriptions",
		"_PartitionKey": "XXXXX",
		"IsActive": true,
		"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"RestartTimeRestriction": false,
		"Type": "WindowsService",
		"Account": "Account for TestWindowsService",
		"PushCreds": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"AccountID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Status": "OK",
		"Resource": "member1.example.com",
		"ACL": "true",
		"Mode": "Push",
		"CurrentMultiplexAccount": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

Delete Service
https://developer.centrify.com/reference#post_subscriptions-deletesubscription

	Request body
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Responde Result
	{
		"success":true,
		"Result":null,
		"Message":null,
		"MessageID":null,
		"Exception":null,
		"ErrorID":null,
		"ErrorCode":null,
		"IsSoftError":false,
		"InnerExceptions":null
	}
*/
