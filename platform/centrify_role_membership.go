package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// RoleMembership - Encapsulates a single Generic RoleMembership
type RoleMembership struct {
	vaultObject
	// API endpoints
	apiGetRoleMembers string

	RoleID  string       `json:"Role,omitempty" schema:"role,omitempty"`
	Members []RoleMember `json:"Members,omitempty" schema:"member,omitempty"`
}

// NewRoleMembership is a role membership constructor
func NewRoleMembership(c *restapi.RestClient) *RoleMembership {
	s := RoleMembership{}
	s.client = c
	s.apiRead = "/SaasManage/GetRole"
	s.apiUpdate = "/Roles/UpdateRole"
	s.apiGetRoleMembers = "/SaasManage/GetRoleMembers"

	return &s
}

// Read function fetches a Role from source, including attribute values. Returns error if any
func (o *RoleMembership) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	if o.RoleID == "" {
		errormsg := fmt.Sprintf("Missing role ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["name"] = o.RoleID
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

	return nil
}

// UpdateRoleMembers adds or removes members into or from a role. Actions are 'Add' or 'Delete'. Types are 'Users', 'Roles', 'Groups'
func (o *RoleMembership) UpdateRoleMembers(members []RoleMember, action string) (*restapi.StringResponse, error) {
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

func (o *RoleMembership) getMembers() ([]RoleMember, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["name"] = o.RoleID
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
