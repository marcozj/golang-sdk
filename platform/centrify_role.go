package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Role - Encapsulates a single Generic Role
type Role struct {
	vaultObject
	// API endpoints
	apiGetRights      string
	apiAssignRights   string
	apiUnassignRights string
	apiGetRoleMembers string

	// Users
	//Users       []string     `json:"Users,omitempty" schema:"users,omitempty"`
	Members     []RoleMember `json:"Members,omitempty" schema:"member,omitempty"`
	AdminRights []string     `json:"AdminRights,omitempty" schema:"adminrights,omitempty"`
}

// RoleMember - Encapsulates a single role member
type RoleMember struct {
	MemberName string `json:"Name,omitempty" schema:"name,omitempty"`
	MemberID   string `json:"Guid,omitempty" schema:"id,omitempty"`
	MemberType string `json:"Type,omitempty" schema:"type,omitempty"`
}

// NewRole is a role constructor
func NewRole(c *restapi.RestClient) *Role {
	s := Role{}
	s.client = c
	s.apiRead = "/SaasManage/GetRole"
	s.apiCreate = "/SaasManage/StoreRole"
	s.apiDelete = "/SaasManage/DeleteRole"
	s.apiUpdate = "/Roles/UpdateRole"
	s.apiGetRights = "/Core/GetAssignedAdministrativeRights"
	s.apiAssignRights = "/saasManage/AssignSuperRights"
	s.apiUnassignRights = "/saasManage/UnAssignSuperRights"
	s.apiGetRoleMembers = "/SaasManage/GetRoleMembers"

	return &s
}

// Read function fetches a Role from source, including attribute values. Returns error if any
func (o *Role) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["name"] = o.ID
	queryArg["suppressPrincipalsList"] = true
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

	mapToStruct(o, resp.Result)

	// Get role members
	members, err := o.getMembers()
	if err != nil {
		return err
	}
	o.Members = members

	// Get admin rights
	rights, err := o.GetAdminRights()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	var r []string
	for k := range rights {
		r = append(r, k)
	}
	o.AdminRights = r

	//logger.Debugf("Filled object: %+v", o)

	return nil
}

// Create function creates a new role and returns a map that contains creation result
func (o *Role) Create() (*restapi.GenericMapResponse, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.Name
	if o.Description != "" {
		queryArg["Description"] = o.Description
	}

	resp, err := o.client.CallGenericMapAPI(o.apiCreate, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
	// Upon successful creation, assign ID
	o.ID = resp.Result["_RowKey"].(string)

	return resp, nil
}

// Delete function deletes a role and returns a map that contains deletion result
func (o *Role) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("Name")
}

// Update function updates a existing role and returns a map that contains update result
func (o *Role) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.ID
	if o.Name != "" {
		queryArg["NewName"] = o.Name
	}
	if o.Description != "" {
		queryArg["Description"] = o.Description
	}

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

// GetAdminRights function fetches admin rights that are assigned to a role
// and returns a map. The map key is admin right name and map value is path of the json file
func (o *Role) GetAdminRights() (map[string]interface{}, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["role"] = o.ID

	rights := make(map[string]interface{})
	resp, err := o.client.CallGenericMapAPI(o.apiGetRights, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	if resp.Success {
		// Results is an array of map[string]interface{}
		results := resp.Result["Results"].([]interface{})
		for _, v := range results {
			resultItem := v.(map[string]interface{})
			row := resultItem["Row"].(map[string]interface{})
			//rights = append(rights, row["Description"].(string))
			rights[row["Description"].(string)] = row["Path"]

		}
	} else {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	return rights, nil
}

// AssignAdminRights function adds admin rights to a role. The rights parameter is a slice of admin right name
// It returns a map that contains call result
func (o *Role) AssignAdminRights() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArgs []map[string]interface{}

	// Only run if admin rights isn't empty
	if o.AdminRights != nil && len(o.AdminRights) > 0 {
		// Fetch full list of admin rights from tenant so that we can know all corresponding path of json files
		allRights, err := getAllAdminRights(o.client)

		// Convert o.AdminRights from list to following format:
		// Role: xxxxx
		// Path: xxxxx
		for _, v := range o.AdminRights {
			queryArg := make(map[string]interface{})
			queryArg["Role"] = o.ID
			queryArg["Path"] = allRights[v]
			queryArgs = append(queryArgs, queryArg)
		}

		resp, err := o.client.CallGenericMapListAPI(o.apiAssignRights, queryArgs)
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

// RemoveAdminRights function removes existing admin rights from a role.
// The rights parameter is a map. The map key is admin right name and map value is path of the json file
func (o *Role) RemoveAdminRights(rights map[string]interface{}) (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArgs []map[string]interface{}
	for _, v := range rights {
		queryArg := make(map[string]interface{})
		queryArg["Role"] = o.ID
		queryArg["Path"] = v
		queryArgs = append(queryArgs, queryArg)
	}

	resp, err := o.client.CallGenericMapListAPI(o.apiUnassignRights, queryArgs)
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

// Helper function to retrieve all avaiable admin rights from source.
// Return a map. The map key is admin right name, and map value is path of the json file
func getAllAdminRights(client *restapi.RestClient) (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Script"] = "@/lib/get_superrights.js(excludeRight:'')"
	queryArg["Args"] = subArgs

	rights := make(map[string]interface{})
	resp, err := client.CallGenericMapAPI("/Redrock/Query", queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	if resp.Success {
		results := resp.Result["Results"].([]interface{})
		for _, v := range results {
			resultItem := v.(map[string]interface{})
			row := resultItem["Row"].(map[string]interface{})
			rights[row["Description"].(string)] = row["Path"]
		}
	} else {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
	logger.Debugf("List of all admin rights: %v", rights)

	return rights, err
}

// UpdateMembers adds or removes members into or from a role. Actions are 'Add' or 'Delete'. Types are 'Users', 'Roles', 'Groups'
func (o *Role) UpdateMembers(ids []string, action string, membertype string) (*restapi.StringResponse, error) {
	var resp *restapi.StringResponse
	var err error
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	if action != "Add" && action != "Delete" {
		return nil, fmt.Errorf("Update role action must be either 'Add' or 'Delete'")
	}
	if membertype != "Users" && membertype != "Roles" && membertype != "Groups" {
		return nil, fmt.Errorf("Role member type must be either 'Users', 'Roles', or 'Groups'")
	}

	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.ID // This is the ID of role
	var actionArg = make(map[string]interface{})
	actionArg[action] = ids
	queryArg[membertype] = actionArg
	resp, err = o.client.CallStringAPI(o.apiUpdate, queryArg)
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

// UpdateRoleMembers adds or removes members into or from a role. Actions are 'Add' or 'Delete'. Types are 'Users', 'Roles', 'Groups'
func (o *Role) UpdateRoleMembers(members []RoleMember, action string) (*restapi.StringResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	if action != "Add" && action != "Delete" {
		return nil, fmt.Errorf("Update role action must be either 'Add' or 'Delete'")
	}
	var queryArg = make(map[string]interface{})
	queryArg["Name"] = o.ID
	var roleids, groupids, userids []string

	for _, member := range members {
		switch member.MemberType {
		case "Role":
			roleids = append(roleids, member.MemberID)
		case "Group":
			groupids = append(groupids, member.MemberID)
		case "User":
			userids = append(userids, member.MemberID)
		}
	}
	if len(roleids) > 0 {
		var actionArg = make(map[string]interface{})
		actionArg[action] = roleids
		queryArg["Roles"] = actionArg
	}
	if len(groupids) > 0 {
		var actionArg = make(map[string]interface{})
		actionArg[action] = groupids
		queryArg["Groups"] = actionArg
	}
	if len(userids) > 0 {
		var actionArg = make(map[string]interface{})
		actionArg[action] = userids
		queryArg["Users"] = actionArg
	}

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

func (o *Role) getMembers() ([]RoleMember, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["name"] = o.ID
	queryArg["Args"] = subArgs

	var members []RoleMember
	resp, err := o.client.CallGenericMapAPI(o.apiGetRoleMembers, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	if resp.Success {
		results := resp.Result["Results"].([]interface{})
		for _, v := range results {
			var member RoleMember
			resultItem := v.(map[string]interface{})
			row := resultItem["Row"].(map[string]interface{})
			member.MemberID = row["Guid"].(string)
			member.MemberName = row["Name"].(string)
			member.MemberType = row["Type"].(string)
			members = append(members, member)
		}
	} else {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	return members, nil
}

// Query function returns a single role object in map format
func (o *Role) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Role WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns role ID by name
func (o *Role) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Role name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving role: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves role from tenant by name
func (o *Role) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of role %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a role by name
func (o *Role) DeleteByName() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of role %s. %v", o.Name, err)
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
	API to manage role
	https://developer.centrify.com/docs/manage-rolesnew

	Fetch role
	https://developer.centrify.com/reference#post_saasmanage-getrole

		Request body format
		{
			"name": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"suppressPrincipalsList": true,
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
				"ReadOnly": false,
				"Description": "AD accounts that are granted local administrator access to non-domain joined machines.",
				"DirectoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"RoleType": "PrincipalList",
				"Name": "LAB Cloud Local Admins",
				"Uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Create role
	https://developer.centrify.com/reference#post_saasmanage-storerole

		Step 1: Create role
		Request body format
		{
			"Name": "role name",
			"Description": "role description"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

		Step 2: Assign admin rights

	Update role
	https://developer.centrify.com/reference#post_roles-updaterole-1

		Request body format
		{
			"Name": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"NewName": "afdsaf",
			"Description": "sdafsdasd afasdf"
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

	Delete role
	https://developer.centrify.com/reference#post_saasmanage-deleterole

		Request body format
		{
			"Name": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Assign rights to role
	https://developer.centrify.com/docs/manage-rolesnew#assigning-rights-to-the-role
	https://developer.centrify.com/reference#post_core-getassignedadministrativerights

		Request body format
		[{
			"Role": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Path": "/lib/rights/cssintegration.json"
		}, {
			"Role": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Path": "/lib/rights/fedman.json"
		}]

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

	Updates role members

	Request body format
	{
		"Users": {
			"Delete": [
				"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			]
		},
		"Roles": {
			"Add": [
				"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			]
		},
		"Groups": {
			"Add": [
				"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			]
		},
		"Name": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Description": "AD accounts who can login to non-domain joined machines but without any privileges."
	}
*/
