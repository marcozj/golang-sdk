package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/computerclass"
	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// System - Encapsulates a single Generic System
type System struct {
	// System -> Settings menu related settings
	vaultObject
	apiGetChallenge                        string
	apiGetPrivilegeElevationChallenge      string
	apiAddToSets                           string
	apiGetAgentAuthWorkflowConfig          string
	apiGetPrivilegeElevationWorkflowConfig string
	//setTable        string

	FQDN          string `json:"FQDN,omitempty" schema:"fqdn,omitempty"`
	ComputerClass string `json:"ComputerClass,omitempty" schema:"computer_class,omitempty"` // Valid values are: Windows, Unix, CiscoIOS, CiscoNXOS, JuniperJunos, HpNonStopOS, IBMi, CheckPointGaia
	// PaloAltoNetworksPANOS, F5NetworksBIGIP, CiscoAsyncOS, VMwareVMkernel, GenericSsh, Customssh
	//SystemProfileId string `json:"SystemProfileId,omitempty" schema:"system_rofile_id,omitempty"` // For Customssh
	SessionType        string `json:"SessionType,omitempty" schema:"session_type,omitempty"`       // Valid values are: Rdp, Ssh
	ManagementMode     string `json:"ManagementMode,omitempty" schema:"management_mode,omitempty"` // Valid values are: RpcOverTcp, Smb, WinRMOverHttp, WinRMOverHttps, Disabled
	ManagementPort     int    `json:"ManagementPort,omitempty" schema:"management_port,omitempty"` // For Windows, F5, PAN-OS and VMKernel only
	Port               int    `json:"Port,omitempty" schema:"port,omitempty"`
	TimeZoneID         string `json:"TimeZoneID,omitempty" schema:"system_timezone,omitempty"` // System Time Zone
	UseMyAccount       bool   `json:"CertAuthEnable,omitempty" schema:"use_my_account,omitempty"`
	Status             string `json:"Status,omitempty" schema:"status,omitempty"`
	ProxyUser          string `json:"ProxyUser" schema:"proxyuser"` // To disable ProxyUser, it needs to be set to "" instead of omitting
	ProxyUserPassword  string `json:"ProxyUserPassword,omitempty" schema:"proxyuser_password,omitempty"`
	ProxyUserIsManaged bool   `json:"ProxyUserIsManaged" schema:"proxyuser_managed"` // ProxyUserIsManaged needs to be set instead of omitting otherwise update fails

	// System -> Policy menu related settings
	DefaultCheckoutTime              int             `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"`                                   // Checkout lifetime (minutes)
	AllowRemote                      bool            `json:"AllowRemote,omitempty" schema:"allow_remote_access,omitempty"`                                         // Allow access from a public network (web client only)
	AllowRdpClipboard                bool            `json:"AllowRdpClipboard,omitempty" schema:"allow_rdp_clipboard,omitempty"`                                   // Allow RDP client to sync local clipboard with remote session
	ChallengeRules                   *ChallengeRules `json:"LoginRules,omitempty" schema:"challenge_rule,omitempty"`                                               // System Login Challenge Rules
	LoginDefaultProfile              string          `json:"LoginDefaultProfile,omitempty" schema:"default_profile_id,omitempty"`                                  // Default System Login Profile (used if no conditions matched)
	PrivilegeElevationDefaultProfile string          `json:"PrivilegeElevationDefaultProfile,omitempty" schema:"privilege_elevation_default_profile_id,omitempty"` // Default Privilege Elevation Profile (used if no conditions matched)
	PrivilegeElevationRules          *ChallengeRules `json:"PrivilegeElevationRules,omitempty" schema:"privilege_elevation_rule,omitempty"`                        // Privilege Elevation Challenge Rules

	// System -> Advanced menu related settings
	AllowAutomaticLocalAccountMaintenance bool   `json:"AllowAutomaticLocalAccountMaintenance,omitempty" schema:"local_account_automatic_maintenance,omitempty"` // Local Account Automatic Maintenance
	AllowManualLocalAccountUnlock         bool   `json:"AllowManualLocalAccountUnlock,omitempty" schema:"local_account_manual_unlock,omitempty"`                 // Local Account Manual Unlock
	DomainID                              string `json:"DomainId,omitempty" schema:"domain_id,omitempty"`                                                        // Domain
	RemoveUserOnSessionEnd                bool   `json:"RemoveUserOnSessionEnd,omitempty" schema:"remove_user_on_session_end,omitempty"`
	AllowMultipleCheckouts                bool   `json:"AllowMultipleCheckouts,omitempty" schema:"allow_multiple_checkouts,omitempty"`                          // Allow multiple password checkouts for this system
	AllowPasswordRotation                 bool   `json:"AllowPasswordRotation,omitempty" schema:"enable_password_rotation,omitempty"`                           // Enable periodic password rotation
	PasswordRotateDuration                int    `json:"PasswordRotateDuration,omitempty" schema:"password_rotate_interval,omitempty"`                          // Password rotation interval (days)
	AllowPasswordRotationAfterCheckin     bool   `json:"AllowPasswordRotationAfterCheckin,omitempty" schema:"enable_password_rotation_after_checkin,omitempty"` // Enable password rotation after checkin
	MinimumPasswordAge                    int    `json:"MinimumPasswordAge,omitempty" schema:"minimum_password_age,omitempty"`                                  // Minimum Password Age (days)
	PasswordProfileID                     string `json:"PasswordProfileID,omitempty" schema:"password_profile_id,omitempty"`                                    // Password Complexity Profile
	AllowPasswordHistoryCleanUp           bool   `json:"AllowPasswordHistoryCleanUp,omitempty" schema:"enable_password_history_cleanup,omitempty"`              // Enable periodic password history cleanup
	PasswordHistoryCleanUpDuration        int    `json:"PasswordHistoryCleanUpDuration,omitempty" schema:"password_historycleanup_duration,omitempty"`          // Password history cleanup (days)

	AllowSSHKeysRotation       bool   `json:"AllowSshKeysRotation,omitempty" schema:"enable_sshkey_rotation,omitempty"`           // Enable periodic SSH key rotation
	SSHKeysRotateDuration      int    `json:"SshKeysRotateDuration,omitempty" schema:"sshkey_rotate_interval,omitempty"`          // SSH key rotation interval (days)
	MinimumSSHKeysAge          int    `json:"MinimumSshKeysAge,omitempty" schema:"minimum_sshkey_age,omitempty"`                  // Minimum SSH Key Age (days)
	SSHKeysGenerationAlgorithm string `json:"SshKeysGenerationAlgorithm,omitempty" schema:"sshkey_algorithm,omitempty"`           // SSH Key Generation Algorithm
	AllowSSHKeysCleanUp        bool   `json:"AllowSshKeysCleanUp,omitempty" schema:"enable_sshkey_history_cleanup,omitempty"`     // Enable periodic SSH key cleanup
	SSHKeysCleanUpDuration     int    `json:"SshKeysCleanUpDuration,omitempty" schema:"sshkey_historycleanup_duration,omitempty"` // SSH key cleanup (days)

	// Workflow
	AgentAuthWorkflowEnabled            bool               `json:"AgentAuthWorkflowEnabled,omitempty" schema:"agent_auth_workflow_enabled,omitempty"` // Enable Agent Auth Workflow
	AgentAuthWorkflowApprovers          []WorkflowApprover `json:"AgentAuthWorkflowApprovers,omitempty" schema:"agent_auth_workflow_approver,omitempty"`
	PrivilegeElevationWorkflowEnabled   bool               `json:"PrivilegeElevationWorkflowEnabled,omitempty" schema:"privilege_elevation_workflow_enabled,omitempty"` // Enable Privilege Elevation Request Workflow
	PrivilegeElevationWorkflowApprovers []WorkflowApprover `json:"PrivilegeElevationWorkflowApprovers,omitempty" schema:"privilege_elevation_workflow_approver,omitempty"`

	// System -> Zone Role Workflow menu related settings
	DomainOperationsEnabled      bool               `json:"DomainOperationsEnabled,omitempty" schema:"use_domainadmin_for_zonerole_workflow,omitempty"` // Use Domain Administrator Account for Zone Role Workflow operations
	ZoneRoleWorkflowEnabled      bool               `json:"ZoneRoleWorkflowEnabled,omitempty" schema:"enable_zonerole_workflow,omitempty"`              // Enable zone role requests for this system
	UseDomainWorkflowRoles       bool               `json:"UseDomainWorkflowRoles" schema:"use_domain_assignment_for_zoneroles"`                        // Assignable Zone Roles - Use domain assignments
	ZoneRoleWorkflowRoles        string             `json:"ZoneRoleWorkflowRoles,omitempty" schema:"assigned_zoneroles,omitempty"`                      // This is the actual attribute in string format
	ZoneRoleWorkflowRoleList     []ZoneRole         `json:"-" schema:"assigned_zonerole,omitempty"`                                                     // This is used in API call and tf file only
	UseDomainWorkflowApprovers   bool               `json:"UseDomainWorkflowApprovers" schema:"use_domain_assignment_for_zonerole_approvers"`           // Approver list - Use domain assignments
	ZoneRoleWorkflowApprovers    string             `json:"ZoneRoleWorkflowApprovers,omitempty" schema:"assigned_zonerole_approvers,omitempty"`         // This is the actual attribute in string format
	ZoneRoleWorkflowApproverList []WorkflowApprover `json:"-" schema:"assigned_zonerole_approver,omitempty"`                                            // This is used in tf file only

	// System -> Connectors menu related settings
	ProxyCollectionList string `json:"ProxyCollectionList,omitempty" schema:"connector_list,omitempty"` // List of Connectors used

	// Sets
	//Sets []string `json:"Sets,omitempty" schema:"sets,omitempty"`
}

type AgentAuthWorkflowConfig struct {
	AgentAuthWorkflowEnabled   bool
	AgentAuthWorkflowApprovers []WorkflowApprover
}

type PrivilegeElevationWorkflowConfig struct {
	PrivilegeElevationWorkflowEnabled   bool
	PrivilegeElevationWorkflowApprovers []WorkflowApprover
}

// NewSystem is a System constructor
func NewSystem(c *restapi.RestClient) *System {
	s := System{}
	s.ValidPermissions = ValidPermissionMap.System
	s.SetType = settype.System.String()
	s.client = c
	s.apiRead = "/RedRock/query"
	s.apiCreate = "/ServerManage/AddResource"
	s.apiDelete = "/ServerManage/DeleteResource"
	s.apiUpdate = "/ServerManage/UpdateResource"
	s.apiGetChallenge = "/ServerManage/GetComputerChallenges"
	s.apiGetPrivilegeElevationChallenge = "/PrivilegeElevation/GetChallenges"
	s.apiAddToSets = "/Collection/UpdateMembersCollection"
	s.apiPermissions = "/ServerManage/SetResourcePermissions"
	s.apiGetAgentAuthWorkflowConfig = "/ServerManage/GetAgentAuthWorkflowConfig"
	s.apiGetPrivilegeElevationWorkflowConfig = "/ServerManage/GetPrivilegeElevationWorkflowConfig"
	//s.setTable = "Server"

	return &s
}

// Read function fetches a system from source, including attribute values. Returns error if any
func (o *System) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Script"] = "SELECT * FROM Server WHERE Server.ID = '" + o.ID + "'"
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
		return fmt.Errorf("System does not exist in tenant")
	} else if len(results) > 1 {
		// this should never happen
		return fmt.Errorf("There are more than one system with the same ID in tenant")
	}
	var result = results[0].(map[string]interface{})
	// Populate vaultObject struct with map from response
	var row = result["Row"].(map[string]interface{})
	mapToStruct(o, row)

	// Get system login profile information
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	if v, ok := resp.Result["LoginDefaultProfile"]; ok {
		o.LoginDefaultProfile = v.(string)
	}
	// Fill login rules
	if v, ok := resp.Result["LoginRules"]; ok {
		challengerules := &ChallengeRules{}
		mapToStruct(challengerules, v.(map[string]interface{}))
		o.ChallengeRules = challengerules
	}

	var args = make(map[string]interface{})
	args["ID"] = o.ID
	// Get Privilege Elevation Challenge profile information
	resp, err = o.client.CallGenericMapAPI(o.apiGetPrivilegeElevationChallenge, args)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	if v, ok := resp.Result["PrivilegeElevationDefaultProfile"]; ok {
		o.PrivilegeElevationDefaultProfile = v.(string)
	}
	// Fill login rules
	if v, ok := resp.Result["PrivilegeElevationRules"]; ok {
		challengerules := &ChallengeRules{}
		mapToStruct(challengerules, v.(map[string]interface{}))
		o.PrivilegeElevationRules = challengerules
	}

	// Get AgentAuth workflow approvers
	resp, err = o.client.CallGenericMapAPI(o.apiGetAgentAuthWorkflowConfig, args)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	// Fill AgentAuthWorkflowApprovers
	aawfconfig := &AgentAuthWorkflowConfig{}
	mapToStruct(aawfconfig, resp.Result)
	o.AgentAuthWorkflowEnabled = aawfconfig.AgentAuthWorkflowEnabled
	o.AgentAuthWorkflowApprovers = aawfconfig.AgentAuthWorkflowApprovers

	// Get privilege elevation workflow approvers
	resp, err = o.client.CallGenericMapAPI(o.apiGetPrivilegeElevationWorkflowConfig, args)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	// Fill PrivilegeElevationWorkflowApprovers
	pewfconfig := &PrivilegeElevationWorkflowConfig{}
	mapToStruct(pewfconfig, resp.Result)
	o.PrivilegeElevationWorkflowEnabled = pewfconfig.PrivilegeElevationWorkflowEnabled
	o.PrivilegeElevationWorkflowApprovers = pewfconfig.PrivilegeElevationWorkflowApprovers

	return nil
}

// Create function creates a new system
func (o *System) Create() (*restapi.StringResponse, error) {
	err := o.processWorkflow()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	err = o.processZoneRoleWorkflow()
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
	// Special handling of system login profile
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

// Delete function deletes a system and returns a map that contains deletion result
func (o *System) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing system and returns a map that contains update result
func (o *System) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	err := o.processWorkflow()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	err = o.processZoneRoleWorkflow()
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
	// Special handling of system login profile
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

// ValidateZoneWorkflow checks if domain_id is set if use_domainadmin_for_zonerole_workflow is true
func (o *System) ValidateZoneWorkflow() error {
	// Before Zone workflow can be enabled, domain must be set
	if o.DomainOperationsEnabled && o.DomainID == "" {
		return fmt.Errorf("domain_id must be set in order to enable use_domainadmin_for_zonerole_workflow")
	}
	return nil
}

// Query function returns a single System object in map format
func (o *System) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Server WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.FQDN != "" {
		query += " AND FQDN='" + o.FQDN + "'"
	}
	if o.ComputerClass != "" {
		query += " AND ComputerClass='" + o.ComputerClass + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns system ID by name
func (o *System) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("System name must be provided")
	}
	if o.ComputerClass == "" {
		return "", fmt.Errorf("Computer class must be provided")
	}
	//if o.FQDN == "" {
	//	return "", fmt.Errorf("FQDN must be provided")
	//}

	result, err := o.Query()
	if err != nil {
		errormsg := fmt.Sprintf("Failed to retrieve system '%s' with type '%s'. %s", o.Name, o.ComputerClass, err)
		logger.Errorf(errormsg)
		return "", fmt.Errorf(errormsg)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves system from tenant by name
func (o *System) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf(err.Error())
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a system by name
func (o *System) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf(err.Error())
		}
	}
	resp, err := o.Delete()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

// ResolveValidPermissions assign valid permissions according to computer class
func (o *System) ResolveValidPermissions() {
	if o.ComputerClass == computerclass.Windows.String() || o.ComputerClass == computerclass.Unix.String() {
		o.ValidPermissions = ValidPermissionMap.WinNix
	} else {
		o.ValidPermissions = ValidPermissionMap.System
	}
}

func (o *System) processWorkflow() error {
	// Resolve guid of each approver
	if o.AgentAuthWorkflowEnabled && o.AgentAuthWorkflowApprovers != nil {
		err := ResolveWorkflowApprovers(o.client, o.AgentAuthWorkflowApprovers)
		if err != nil {
			return err
		}
	}
	if o.PrivilegeElevationWorkflowEnabled && o.PrivilegeElevationWorkflowApprovers != nil {
		err := ResolveWorkflowApprovers(o.client, o.PrivilegeElevationWorkflowApprovers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *System) processZoneRoleWorkflow() error {
	// Due to historical reason, ZoneRoleWorkflowRoles and ZoneRoleWorkflowApprovers attributes are not in json format rather they are in string so need to perform conversion
	if o.ZoneRoleWorkflowEnabled {
		if !o.UseDomainWorkflowRoles && o.ZoneRoleWorkflowRoleList != nil {
			// Resolve zone role attributes using provided zone role name
			err := resolveZoneRoles(o.client, o.ZoneRoleWorkflowRoleList, o.DomainID)
			if err != nil {
				return err
			}
			// Convert zone roles from struct to string
			o.ZoneRoleWorkflowRoles = FlattenZoneRoles(o.ZoneRoleWorkflowRoleList)
		}

		// Resolve guid of each approver
		if !o.UseDomainWorkflowApprovers && o.ZoneRoleWorkflowApproverList != nil {
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
	API to manage system
	https://centrify-dev.readme.io/docs/add-resourcesnew

	Fetch System

		Request body format
		{
			"Script": "SELECT * FROM Server WHERE Server.ID = 'xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx'",
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
								"Type": "Server",
								"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
								"IsForeignKey": false
							}
						],
						"Row": {
							"AllowPasswordRotationAfterCheckin": null,
							"ProxyUser": "admin@example.com",
							"UseDomainWorkflowApprovers": null,
							"JoinedDate": null,
							"CanonicalName": null,
							"AllowPasswordRotation": null,
							"ZoneRoleWorkflowEnabled": null,
							"DisableNla": null,
							"SshKeyFolder": null,
							"AllowRemote": null,
							"Description": "Windows system 1",
							"DiscoveryRefreshTime": null,
							"ProxyUserIsManaged": false,
							"NumGoodAccounts": -1,
							"JoinedBy": null,
							"Port": null,
							"UniqueId": null,
							"ZoneStatus": null,
							"Rights": "ManageSession, Edit, Delete, Grant, AgentAuth, RequestZoneRole, View, AddAccount, UnlockAccount, OfflineRescue",
							"NumBrokenAccounts": -1,
							"DomainId": null,
							"AgentVersion": null,
							"AdministrativeAccountID": null,
							"PasswordProfileID": null,
							"HealthStatus": "Unreachable",
							"OperatingSystemServicePack": null,
							"MinimumPasswordAge": null,
							"SystemProfileId": null,
							"FQDN": "192.168.2.3",
							"ZoneRoleWorkflowApprovers": null,
							"DomainName": null,
							"AllowPasswordHistoryCleanUp": null,
							"ZoneRoleWorkflowApproversList": null,
							"Reachable": false,
							"DiscoveredTime": null,
							"DefaultHome": null,
							"CertAuthEnable": false,
							"PostalAddress": null,
							"IsFavorite": false,
							"ProxyCollectionList": null,
							"ManagementMode": null,
							"OperatingSystem": null,
							"LastState": "Unreachable",
							"ServiceAccountID": null,
							"ActiveSessions": 0,
							"ShowCpsOnMobile": null,
							"ZoneJoined": null,
							"PasswordRotateDuration": null,
							"TimeZoneID": null,
							"AllowAutomaticLocalAccountMaintenance": false,
							"ManagementPort": null,
							"HealthCheckInterval": null,
							"PasswordHistoryCleanUpDuration": null,
							"UserID": null,
							"ZoneRoleWorkflowRoles": null,
							"SessionType": "Rdp",
							"CredentialKmipMode": null,
							"ComputerClass": "Windows",
							"Name": "Windows 01",
							"DomainOperationsEnabled": false,
							"AllowHealthCheck": null,
							"IPAddress": null,
							"Joined": null,
							"DefaultShell": null,
							"_MatchFilter": null,
							"AllowManualLocalAccountUnlock": false,
							"HealthStatusError": "_I18N_NoCloudConnectorsError",
							"ActiveCheckouts": 0,
							"DefaultCheckoutTime": null,
							"UseDomainWorkflowRoles": null,
							"ReachableError": "_I18N_NoCloudConnectorsError",
							"DiscoveryAccountId": null,
							"ProxyUserKmipId": null,
							"Accounts": null,
							"ComputerClassDisplayName": "Windows",
							"LastHealthCheck": "/Date(1596262439557)/",
							"AdministrativeAccountEnabled": null,
							"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"AllowMultipleCheckouts": null,
							"DistinguishedName": null
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

	Create System
	https://developer.centrify.com/reference#post_servermanage-addresource

		Request body format
		{
			"ComputerClass": "Windows",
			"FQDN": "127.0.0.1",
			"Name": "Windows Server 01",
			"Description": "Windows Server 01",
			"SessionType": "Rdp"
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

	Update System
	https://developer.centrify.com/reference#post_servermanage-updateresource

		Request body format
		{
			"ProxyUser": "admin@example.com",
			"Description": "Windows Server 01 Test",
			"ProxyUserIsManaged": false,
			"Rights": "ManageSession, Edit, Delete, Grant, AgentAuth, RequestZoneRole, View, AddAccount, UnlockAccount, OfflineRescue",
			"HealthStatus": "Unreachable",
			"FQDN": "127.0.0.1",
			"Reachable": false,
			"CertAuthEnable": false,
			"IsFavorite": false,
			"ProxyCollectionList": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx,xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"LastState": "Unreachable",
			"ActiveSessions": 0,
			"SessionType": "Rdp",
			"ComputerClass": "Windows",
			"Name": "Windows Server 01",
			"HealthStatusError": "_I18N_NoCloudConnectorsError",
			"ActiveCheckouts": 0,
			"ReachableError": "_I18N_NoCloudConnectorsError",
			"ComputerClassDisplayName": "Windows",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"SessionCount": 0,
			"DisableNla": false,
			"ProxyUserPassword": "xxxxxxxx",
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": true,
			"resetProxyUserPassword": false,
			"AdministrativeAccountID": null
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

	Delete System
	https://developer.centrify.com/reference#post_servermanage-deleteresource

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

		Create Login rules
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"LoginRules": {
				"_Type": "RowSet",
				"Enabled": true,
				"_UniqueKey": "Condition",
				"_Value": [
					{
						"Conditions": [
							{
								"Prop": "IpAddress",
								"Op": "OpInCorpIpRange"
							},
							{
								"Prop": "DayOfWeek",
								"Op": "OpIsDayOfWeek",
								"Val": "L,1,2,3,4,5"
							}
						],
						"ProfileId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
					},
					{
						"Conditions": [
							{
								"Prop": "DeviceOs",
								"Op": "OpEqual",
								"Val": "WindowsMobile"
							}
						],
						"ProfileId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
					}
				]
			},
			"SessionCount": 0,
		}

		Fetch Login rules
		{
			"success": true,
			"Result": {
				"LoginRules": {
					"_UniqueKey": "Condition",
					"_Value": [
						{
							"Conditions": [
								{
									"Prop": "IpAddress",
									"Op": "OpInCorpIpRange"
								}
							],
							"ProfileId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
						}
					],
					"Enabled": true,
					"_Type": "RowSet"
				},
				"LoginDefaultProfile": ""
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
