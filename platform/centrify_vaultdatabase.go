package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Database - Encapsulates a single Database
type Database struct {
	// Database -> Settings menu related settings
	vaultObject
	apiGetChallenge string
	apiAddToSets    string
	setTable        string

	FQDN                 string `json:"FQDN,omitempty" schema:"hostname,omitempty"`
	DatabaseClass        string `json:"DatabaseClass,omitempty" schema:"database_class,omitempty"` // Valid values are: SQLServer, Oracle, SAPAse
	Port                 int    `json:"Port,omitempty" schema:"port,omitempty"`
	InstanceName         string `json:"InstanceName,omitempty" schema:"instance_name,omitempty"` // MS SQL instance name
	ServiceName          string `json:"ServiceName,omitempty" schema:"service_name,omitempty"`   // Oracle database service name
	SkipReachabilityTest bool   `json:"SkipReachabilityTest,omitempty" schema:"skip_reachability_test,omitempty"`

	// Database -> Policy menu related settings
	DefaultCheckoutTime int `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)

	// Database -> Advanced menu related settings
	AllowMultipleCheckouts            bool   `json:"AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts for related accounts
	AllowPasswordRotation             bool   `json:"AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int    `json:"PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool   `json:"AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int    `json:"MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	PasswordProfileID                 string `json:"PasswordProfileID,omitempty" schema:"password_profile_id,omitempty"`                                    // Password Complexity Profile
	AllowPasswordHistoryCleanUp       bool   `json:"AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`              // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration    int    `json:"PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"`          // Password history cleanup (days)

	// Database -> Connectors menu related settings
	ProxyCollectionList string `json:"ProxyCollectionList,omitempty" schema:"connector_list,omitempty"` // List of Connectors used
}

// NewDatabase is a Database constructor
func NewDatabase(c *restapi.RestClient) *Database {
	s := Database{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Database
	s.SetType = settype.Database.String()
	s.apiRead = "/RedRock/query"
	s.apiCreate = "/ServerManage/AddDatabase"
	s.apiDelete = "/ServerManage/DeleteDatabase"
	s.apiUpdate = "/ServerManage/UpdateDatabase"
	s.apiGetChallenge = "/ServerManage/GetComputerChallenges"
	s.apiAddToSets = "/Collection/UpdateMembersCollection"
	s.apiPermissions = "/ServerManage/SetDatabasePermissions"
	s.setTable = "VaultDatabase"

	return &s
}

// Read function fetches a Database from source, including attribute values. Returns error if any
func (o *Database) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Script"] = "SELECT * FROM VaultDatabase WHERE VaultDatabase.ID = '" + o.ID + "'"
	queryArg["Args"] = subArgs

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
		logger.Debugf("Returning Database does not exist in tenant")
		return fmt.Errorf("Database does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return fmt.Errorf("There are more than one Database with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})
	//logger.Debugf("Input map: %+v", row)
	mapToStruct(o, row)

	//logger.Debugf("Filled object: %+v", o)

	return nil
}

// Create function creates a new Database and returns a map that contains creation result
func (o *Database) Create() (*restapi.StringResponse, error) {
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

// Delete function deletes a Database and returns a map that contains deletion result
func (o *Database) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing Database and returns a map that contains update result
func (o *Database) Update() (*restapi.GenericMapResponse, error) {
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

// Query function returns a single database object in map format
func (o *Database) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM VaultDatabase WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.FQDN != "" {
		query += " AND FQDN='" + o.FQDN + "'"
	}
	if o.DatabaseClass != "" {
		query += " AND DatabaseClass='" + o.DatabaseClass + "'"
	}
	if o.InstanceName != "" {
		query += " AND InstanceName='" + o.InstanceName + "'"
	}
	if o.ServiceName != "" {
		query += " AND ServiceName='" + o.ServiceName + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns database ID by name
func (o *Database) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Database name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving database: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves database from tenant by name
func (o *Database) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of database %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a database by name
func (o *Database) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of database %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
	API to manage system

	Fetch Database

		Request body format
		{
			"Script": "SELECT * FROM VaultDatabase WHERE VaultDatabase.ID = 'xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx'",
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
        "IsAggregate": false,
        "Count": 1,
        ...
        "FullCount": 1,
        "Results": [
            {
                "Entities": [
                    {
                        "Type": "VaultDatabase",
                        "Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Name": "Test DB",
                    "LastHealthCheck": "/Date(1597222428470)/",
                    "PasswordRotateInterval": null,
                    "HealthCheckInterval": null,
                    "ProxyCollectionList": null,
                    "PasswordHistoryCleanUpDuration": null,
                    "ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                    "MinimumPasswordAge": null,
                    "Description": "Test DB Instance 2",
                    "HealthStatus": "Unreachable",
                    "IPAddress": null,
                    "PasswordProfileID": null,
                    "InstanceName": "INSTANCE",
                    "FQDN": "127.0.0.1",
                    "DefaultCheckoutTime": null,
                    "AllowPasswordRotation": null,
                    "ReachableError": "_I18N_NoCloudConnectorsError",
                    "AllowPasswordHistoryCleanUp": null,
                    "AllowHealthCheck": null,
                    "ServiceName": null,
                    "LastState": "Unreachable",
                    "PasswordRotateDuration": null,
                    "HealthStatusError": "_I18N_NoCloudConnectorsError",
                    "Reachable": false,
                    "AllowMultipleCheckouts": null,
                    "_MatchFilter": null,
                    "DatabaseClass": "SQLServer",
                    "AllowPasswordRotationAfterCheckin": null,
                    "Port": 1433
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

	Create Database
	https://developer.centrify.com/reference#post_servermanage-adddatabase

		Request body format
		{
			"Name": "Test DB",
			"DatabaseClass": "SQLServer",
			"FQDN": "127.0.0.1",
			"Port": 1433,
			"InstanceName": "INSTANCE",
			"Description": "Test DB Instance",
			"SkipReachabilityTest": false
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

	Update Database
	https://developer.centrify.com/reference#post_servermanage-updatedatabase

		Request body format
		{
			"Name": "Test DB",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": "Test DB Instance 2",
			"HealthStatus": "Unreachable",
			"InstanceName": "INSTANCE",
			"FQDN": "127.0.0.1",
			"ReachableError": "_I18N_NoCloudConnectorsError",
			"LastState": "Unreachable",
			"HealthStatusError": "_I18N_NoCloudConnectorsError",
			"Reachable": false,
			"DatabaseClass": "SQLServer",
			"Port": 1433,
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Delete Database
	https://developer.centrify.com/reference#post_servermanage-deletedatabase

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
*/
