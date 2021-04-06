package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Domain - Encapsulates a single Domain
type Domain struct {
	vaultObject
	apiSetAdminAccount string
	apiCanDelete       string

	VerifyDomain bool   `json:"VerifyDomain,omitempty" schema:"verify,omitempty"`
	ParentID     string `json:"ParentID,omitempty" schema:"parent_id,omitempty"`
	ForestID     string `json:"ForestID,omitempty" schema:"forest_id,omitempty"`
	// Policy menu related settings
	DefaultCheckoutTime int `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	// Advanced menu -> Administrative Account Settings
	AdminAccountID           string `json:"Administrator,omitempty" schema:"administrative_account_id,omitempty"`
	AdministratorDisplayName string `json:"AdministratorDisplayName,omitempty" schema:"administrator_display_name,omitempty"`
	//AdminAccountDomain           string `json:"AdminAccountDomain,omitempty" schema:"administrative_account_domain,omitempty"`
	AdminAccountPassword         string `json:"AdminAccountPassword,omitempty" schema:"administrative_account_password,omitempty"`
	AdminAccountName             string `json:"AdminAccountName,omitempty" schema:"administrative_account_name,omitempty"`
	AutoDomainAccountMaintenance bool   `json:"AllowAutomaticAccountMaintenance" schema:"auto_domain_account_maintenance"`     // Enable Automatic Domain Account Maintenance
	AutoLocalAccountMaintenance  bool   `json:"AllowAutomaticLocalAccountMaintenance" schema:"auto_local_account_maintenance"` // Enable Automatic Local Account Maintenance
	ManualDomainAccountUnlock    bool   `json:"AllowManualAccountUnlock" schema:"manual_domain_account_unlock"`                // Enable Manual Domain Account Unlock
	ManualLocalAccountUnlock     bool   `json:"AllowManualLocalAccountUnlock" schema:"manual_local_account_unlock"`            // Enable Manual Local Account Unlock
	ProvisioningAdminID          string `json:"ProvisioningAdminID,omitempty" schema:"provisioning_admin_id,omitempty"`        // An administrative account to provision the reconciliation account on Unix systems. (must be managed)
	ReconciliationAccountName    string `json:"ReconciliationAccountName,omitempty" schema:"reconciliation_account_name,omitempty"`
	// Advanced -> Security Settings
	AllowMultipleCheckouts            bool   `json:"AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts per AD account added for this domain
	AllowPasswordRotation             bool   `json:"AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration            int    `json:"PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin bool   `json:"AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                int    `json:"MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	PasswordProfileID                 string `json:"PasswordProfileID,omitempty" schema:"password_profile_id,omitempty"`                                    // Password Complexity Profile
	// Advanced -> Maintenance Settings
	AllowPasswordHistoryCleanUp    bool `json:"AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`     // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration int  `json:"PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"` // Password history cleanup (days)
	// Advanced -> Domain/Zone Tasks
	AllowRefreshZoneJoined           bool `json:"AllowRefreshZoneJoined,omitempty" schema:"enable_zone_joined_check,omitempty"`             // Enable periodic domain/zone joined check
	RefreshZoneJoinedIntervalMinutes int  `json:"RefreshZoneJoinedIntervalMinutes,omitempty" schema:"zone_joined_check_interval,omitempty"` // Domain/zone joined check interval (minutes)
	AllowZoneRoleCleanup             bool `json:"AllowZoneRoleCleanup,omitempty" schema:"enable_zonerole_cleanup,omitempty"`                // Enable periodic removal of expired zone role assignments
	ZoneRoleCleanupIntervalHours     int  `json:"ZoneRoleCleanupIntervalHours,omitempty" schema:"zonerole_cleanup_interval,omitempty"`      // Expired zone role assignment removal interval (hours)
	// Zone Role Workflow
	ZoneRoleWorkflowEnabled      bool               `json:"ZoneRoleWorkflowEnabled" schema:"zonerole_workflow_enabled"`                         // Enable zone role requests for systems in this domain
	ZoneRoleWorkflowRoles        string             `json:"ZoneRoleWorkflowRoles,omitempty" schema:"assigned_zoneroles,omitempty"`              // Assignable zone roles
	ZoneRoleWorkflowRoleList     []ZoneRole         `json:"-" schema:"assigned_zonerole,omitempty"`                                             // This is used in tf file only
	ZoneRoleWorkflowApprovers    string             `json:"ZoneRoleWorkflowApprovers,omitempty" schema:"assigned_zonerole_approvers,omitempty"` // This is the actual attribute in string format
	ZoneRoleWorkflowApproverList []WorkflowApprover `json:"-,omitempty" schema:"assigned_zonerole_approver,omitempty"`                          // This is used in tf file only
	// System -> Connectors menu related settings
	ProxyCollectionList string `json:"ProxyCollectionList,omitempty" schema:"connector_list,omitempty"` // List of Connectors used
	// Sets
	//Sets []string `json:"Sets,omitempty" schema:"sets,omitempty"`
}

// NewDomain is a Domain constructor
func NewDomain(c *restapi.RestClient) *Domain {
	s := Domain{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Domain
	s.SetType = settype.Domain.String()
	s.apiRead = "/RedRock/query"
	s.apiCreate = "/ServerManage/AddDomain"
	s.apiDelete = "/ServerManage/DeleteDomain"
	s.apiUpdate = "/ServerManage/UpdateDomain"
	s.apiPermissions = "/ServerManage/SetDomainPermissions"
	s.apiSetAdminAccount = "/ServerManage/SetAdministrativeAccounts"
	s.apiCanDelete = "/ServerManage/CanDeleteDomain"

	return &s
}

// Read function fetches a Domain from source, including attribute values. Returns error if any
func (o *Domain) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Script"] = "SELECT * FROM VaultDomain WHERE ID = '" + o.ID + "'"
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
		return fmt.Errorf("Domain does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return fmt.Errorf("There are more than one domains with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})
	//logger.Debugf("Input map: %+v", row)
	mapToStruct(o, row)

	return nil
}

// Create function creates a new Domain and returns a map that contains creation result
func (o *Domain) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.Name
	queryArg["VerifyDomain"] = o.VerifyDomain

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

// Delete function deletes a Domain and returns a map that contains deletion result
func (o *Domain) Delete() (*restapi.BoolResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	// Check if the domain can be deleted
	resp, err := o.client.CallGenericMapAPI(o.apiCanDelete, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
	if resp.Result["can"].(bool) {
		return o.deleteObjectBoolAPI("")
	}

	logger.Debugf("Domain cannot be deleted: %+v", resp.Result["why"])
	return nil, nil
}

// Update function updates an existing Domain and returns a map that contains update result
func (o *Domain) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	err := o.processZoneRoleWorkflow()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
	if err != nil {
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
func (o *Domain) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM VaultDomain WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// SetAdminAccount sets domain administrative account
func (o *Domain) SetAdminAccount() error {
	if o.Name == "" {
		return fmt.Errorf("Domain name is empty")
	}
	var queryArg = make(map[string]interface{})
	queryArg["Domains"] = []string{o.Name}
	if o.AdminAccountID != "" {
		if o.AdminAccountPassword != "" {
			// This is a non-vaulted AD account
			queryArg["User"] = o.AdminAccountName
			queryArg["UserUuid"] = o.AdminAccountID
			queryArg["Password"] = o.AdminAccountPassword
		} else {
			// This is a vaulted account
			queryArg["PVID"] = o.AdminAccountID
		}
		if o.AdminAccountName != "" {
			//queryArg["User"] = o.AdminAccountName
		}
	}

	resp, err := o.client.CallGenericMapAPI(o.apiSetAdminAccount, queryArg)
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

// GetIDByName returns domain ID by name
func (o *Domain) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Domain name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving domain: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves domain from tenant by name
func (o *Domain) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of domain %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a domain by name
func (o *Domain) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of domain %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

func (o *Domain) processZoneRoleWorkflow() error {
	// Due to historical reason, ZoneRoleWorkflowRoles and ZoneRoleWorkflowApprovers attributes are not in json format rather they are in string so need to perform conversion
	if o.ZoneRoleWorkflowEnabled {
		if o.ZoneRoleWorkflowRoleList != nil {
			// Resolve zone role attributes using provided zone role name
			err := resolveZoneRoles(o.client, o.ZoneRoleWorkflowRoleList, o.ID)
			if err != nil {
				return err
			}
			// Convert zone roles from struct to string
			o.ZoneRoleWorkflowRoles = FlattenZoneRoles(o.ZoneRoleWorkflowRoleList)
		}

		// Resolve guid of each approver
		if o.ZoneRoleWorkflowApproverList != nil {
			err := ResolveWorkflowApprovers(o.client, o.ZoneRoleWorkflowApproverList)
			if err != nil {
				return err
			}
			// Convert approvers from struct to string so that it can be assigned to the actual attribute used for privision.
			o.ZoneRoleWorkflowApprovers = FlattenWorkflowApprovers(o.ZoneRoleWorkflowApproverList)
			//logger.Debugf("Converted approvers: %+v", o.WorkflowApprovers)
		}
	}

	return nil
}

/*
	Fetch Domain

		Request body format


		Respond result
		{
			"success": true,
			"Result": {
				"IsAggregate": false,
				"Count": 1,
				"Columns": [
					...
				],
				"FullCount": 1,
				"Results": [
					{
						"Entities": [
							{
								"Type": "VaultDomain",
								"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
								"IsForeignKey": false
							}
						],
						"Row": {
							"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"AllowMultipleCheckouts": null,
							"DefaultCheckoutTime": null,
							"MinimumPasswordAge": null,
							"AllowManualLocalAccountUnlock": true,
							"ZoneRoleWorkflowRoles": "[{\"Name\":\"cfyl-SSH Management/Unix Zone\",\"Unix\":true,\"ZoneDn\":\"CN=Unix Zone,CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\",\"Description\":\"Manage SSHD daemon and configuration\",\"ZoneCanonicalName\":\"centrifylab.aws/Centrify/Zones/Global/Unix Zone\",\"ParentZoneDn\":\"CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\"},{\"Name\":\"cfyl-Super System Admin/Unix Zone\",\"Unix\":true,\"ZoneDn\":\"CN=Unix Zone,CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\",\"Description\":\"Super system admin with root access\",\"ZoneCanonicalName\":\"centrifylab.aws/Centrify/Zones/Global/Unix Zone\",\"ParentZoneDn\":\"CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\"},{\"Name\":\"cfyw-Windows System Admin/Windows Zone\",\"ZoneDn\":\"CN=Windows Zone,CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\",\"Description\":\"Windows system admin for managing local configurations\",\"ZoneCanonicalName\":\"centrifylab.aws/Centrify/Zones/Global/Windows Zone\",\"Windows\":true,\"ParentZoneDn\":\"CN=Global,CN=Zones,OU=Centrify,DC=centrifylab,DC=aws\"}]",
							"AllowAutomaticAccountUnlock": null,
							"ForestID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"PasswordHistoryCleanUpDuration": null,
							"AllowHealthCheck": null,
							"HealthStatusError": "Success",
							"UniqueId": null,
							"Reachable": true,
							"LastState": "OK",
							"AllowZoneRoleCleanup": true,
							"ZoneRoleWorkflowApproversList": "[{\"BackupApprover\":{\"PType\":\"Role\",\"ObjectType\":\"Role\",\"Name\":\"LAB Infrastructure Owners\",\"Guid\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"_ID\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1583273484\",\"Principal\":\"LAB Infrastructure Owners\",\"Description\":\"Approver who can approve access to access lab systems.\",\"RoleType\":\"PrincipalList\",\"ReadOnly\":false,\"Type\":\"Role\",\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\"},\"NoManagerAction\":\"useBackup\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1583273484\",\"OptionsSelector\":true,\"Type\":\"Manager\"}]",
							"LastZoneRoleCleanup": "/Date(1595383548740)/",
							"Administrator": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"SyncFromConnector": true,
							"PasswordProfileID": null,
							"PasswordRotateDuration": null,
							"Description": "",
							"ProxyCollectionList": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"AdministratorDisplayName": "admin (example.com)",
							"AllowPasswordRotationAfterCheckin": null,
							"AllowAutomaticAccountMaintenance": true,
							"Name": "example.com",
							"AllowRefreshZoneJoined": true,
							"ZoneRoleWorkflowApprovers": "[{\"BackupApprover\":{\"PType\":\"Role\",\"ObjectType\":\"Role\",\"Name\":\"LAB Infrastructure Owners\",\"Guid\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"_ID\":\"9e6022c7_028d_47a8_aecb_aa952201c221\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1583273484\",\"Description\":\"Approver who can approve access to access lab systems.\",\"Principal\":\"LAB Infrastructure Owners\",\"RoleType\":\"PrincipalList\",\"ReadOnly\":false,\"DirectoryServiceUuid\":\"09B9A9B0-6CE8-465F-AB03-65766D33B05E\",\"Type\":\"Role\"},\"NoManagerAction\":\"useBackup\",\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1583273484\",\"OptionsSelector\":true,\"Type\":\"Manager\"}]",
							"AllowAutomaticLocalAccountMaintenance": true,
							"LastHealthCheck": "/Date(1595383490412)/",
							"ZoneRoleWorkflowEnabled": true,
							"ReachableError": "Success",
							"HealthCheckInterval": null,
							"ParentID": null,
							"PasswordRotateInterval": null,
							"ZoneRoleCleanupIntervalHours": 6,
							"AllowPasswordHistoryCleanUp": null,
							"RefreshZoneJoinedIntervalMinutes": 1440,
							"DiscoveredTime": null,
							"_MatchFilter": null,
							"LastRefreshZoneJoined": "/Date(1598764755975)/",
							"HealthStatus": "OK",
							"AllowManualAccountUnlock": true,
							"AllowPasswordRotation": null
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

	Create Domain
	https://developer.centrify.com/reference#post_servermanage-adddomain

		Request body format
		{
			"Name": "example.com",
			"VerifyDomain": false
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

	Update Domain
	https://developer.centrify.com/reference#post_servermanage-updatedomain

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"AllowMultipleCheckouts": true,
			"DefaultCheckoutTime": 30,
			"MinimumPasswordAge": 98,
			"AllowManualLocalAccountUnlock": true,
			"PasswordHistoryCleanUpDuration": 100,
			"HealthStatusError": "UnknownError",
			"Reachable": false,
			"LastState": "Unreachable",
			"AllowZoneRoleCleanup": true,
			"PasswordProfileID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"PasswordRotateDuration": 90,
			"Description": "example domain",
			"AdministratorDisplayName": "example.om\\admin",
			"AllowPasswordRotationAfterCheckin": true,
			"AllowAutomaticAccountMaintenance": true,
			"Name": "example.com",
			"AllowRefreshZoneJoined": true,
			"AllowAutomaticLocalAccountMaintenance": true,
			"ReachableError": "UnknownError",
			"ZoneRoleCleanupIntervalHours": 6,
			"AllowPasswordHistoryCleanUp": true,
			"RefreshZoneJoinedIntervalMinutes": 1440,
			"HealthStatus": "Unreachable",
			"AllowManualAccountUnlock": true,
			"AllowPasswordRotation": true,
			"newAdministrator": {
				"User": "example.com\\admin",
				"PVID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Domains": [
					"example.com"
				]
			},
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

	Delete Domain
	https://developer.centrify.com/reference#post_servermanage-deletedomain

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

	Set Administrative Account

	Request body format
	{
		"User": "example.com\\admin",
		"PVID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Domains": [
			"example.com"
		]
	}
	or
	{
		"User": "Administrator@example.com",
		"UserUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Password": "xxxxxx",
		"Domains": [
			"example.com"
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

	CanDelete Domain?

	Respond result
	{
		"success": true,
		"Result": {
			"can": true,
			"why": null
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}
*/
