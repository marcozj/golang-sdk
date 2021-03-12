package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Generic object
type vaultObject struct {
	// Rest client
	client *restapi.RestClient
	// Standard attributes
	ID               string            `json:"ID,omitempty" schema:"id,omitempty"`
	Name             string            `json:"Name,omitempty" schema:"name,omitempty"`
	Description      string            `json:"Description,omitempty" schema:"description,omitempty"`
	ValidPermissions map[string]string `json:"-"`

	// Sets
	SetType     string       `json:"-"`
	Sets        []string     `json:"Sets,omitempty" schema:"sets,omitempty"`
	Permissions []Permission `json:"-"`

	// API endpoints
	apiRead        string //`json:"-"` // Ignoring this JSON field
	apiCreate      string //`json:"-"` // Ignoring this JSON field
	apiDelete      string //`json:"-"` // Ignoring this JSON field
	apiUpdate      string //`json:"-"` // Ignoring this JSON field
	apiPermissions string
}

// Permission represents object permission
type Permission struct {
	PrincipalID   string   `json:"PrincipalId,omitempty" schema:"principal_id,omitempty"` // Uuid of the principal
	PrincipalName string   `json:"Principal,omitempty" schema:"principal_name,omitempty"` // User name or role name
	PrincipalType string   `json:"PType,omitempty" schema:"principal_type,omitempty"`     // Principal type: User, Role etc..
	Rights        string   `json:"Rights,omitempty" schema:"rights,omitempty"`            // Permissions: Grant,View,Edit,Delete or None to remove this item
	RightList     []string `json:"-"`
}

// ChallengeRules represents list of login rule set
type ChallengeRules struct {
	Enabled   bool            `json:"Enabled,omitempty" schema:"enabled,omitempty"`
	UniqueKey string          `json:"_UniqueKey,omitempty" schema:"unique_key,omitempty"`
	Type      string          `json:"_Type,omitempty" schema:"type,omitempty"`
	Rules     []ChallengeRule `json:"_Value,omitempty" schema:"rule,omitempty"`
}

// ChallengeRule represents a set of login rule
type ChallengeRule struct {
	ChallengeCondition []ChallengeCondition `json:"Conditions,omitempty" schema:"rule,omitempty"`
	AuthProfileID      string               `json:"ProfileId,omitempty" schema:"authentication_profile_id,omitempty"` // "-1" means Not Allowed
}

// AccessKey represents AWS access key
type AccessKey struct {
	ID              string `json:"ID,omitempty" schema:"id,omitempty"`
	AccessKeyID     string `json:"AccessKeyId,omitempty" schema:"access_key_id,omitempty"`
	SecretAccessKey string `json:"SecretAccessKey,omitempty" schema:"secret_access_key,omitempty"`
}

// ChallengeCondition represents a single challenge rule
type ChallengeCondition struct {
	Filter    string `json:"Prop,omitempty" schema:"filter,omitempty"`
	Condition string `json:"Op,omitempty" schema:"condition,omitempty"`
	Value     string `json:"Val,omitempty" schema:"value,omitempty"`
}

// deleteObjectBoolAPI a object and returns a map that contains deletion result
func (o *vaultObject) deleteObjectBoolAPI(idfield string) (*restapi.BoolResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var funcArg = make(map[string]interface{})
	if idfield == "" {
		funcArg["ID"] = o.ID
	} else {
		funcArg[idfield] = o.ID
	}

	resp, err := o.client.CallBoolAPI(o.apiDelete, funcArg)
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

// deleteObjectMapAPI a object and returns a map that contains deletion result
func (o *vaultObject) deleteObjectMapAPI(idfield string) (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var funcArg = make(map[string]interface{})
	if idfield == "" {
		funcArg["ID"] = o.ID
	} else {
		funcArg[idfield] = o.ID
	}

	reply, err := o.client.CallGenericMapAPI(o.apiDelete, funcArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, fmt.Errorf(reply.Message)
	}

	return reply, nil
}

// deleteObjectStringAPI a object and returns a map that contains deletion result
func (o *vaultObject) deleteObjectStringAPI(idfield string) (*restapi.StringResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var funcArg = make(map[string]interface{})
	if idfield == "" {
		funcArg["ID"] = o.ID
	} else {
		funcArg[idfield] = o.ID
	}

	reply, err := o.client.CallStringAPI(o.apiDelete, funcArg)
	if err != nil {
		return nil, err
	}

	if !reply.Success {
		return nil, fmt.Errorf(reply.Message)
	}

	return reply, nil
}

// SetPermissions sets permissions. isRemove indicates whether to remove all permissions instead of setting permissions
func (o *vaultObject) SetPermissions(isRemove bool) (*restapi.BaseAPIResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	var permissions []map[string]interface{}
	for _, v := range o.Permissions {
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
		logger.Debugf("Generated Map for SetPermissions(): %+v", queryArg)
		resp, err := o.client.CallBaseAPI(o.apiPermissions, queryArg)
		if err != nil {
			return nil, err
		}
		if !resp.Success {
			return nil, fmt.Errorf(fmt.Sprintf("%s %s", resp.Message, resp.Exception))
		}
		return resp, nil
	}
	return nil, nil
}

// FillStruct function fills a struct with map
func (o *vaultObject) FillStruct(m map[string]interface{}) error {
	logger.Debugf("Input map: %v", m)
	for k, v := range m {
		logger.Debugf("Map key %v, map value: %v", k, v)
		err := setField(o, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddToSetsByID add database into Sets
func (o *vaultObject) AddToSetsByID(setids []string) error {
	for _, v := range setids {
		setObj := NewManualSet(o.client)
		setObj.ID = v
		setObj.ObjectType = o.SetType
		resp, err := setObj.UpdateSetMembers([]string{o.ID}, "add")
		if err != nil || !resp.Success {
			return fmt.Errorf("Error adding %s to Sets: %v", o.Name, err)
		}
	}
	return nil
}

// AddToSetsByName add database into Sets
func (o *vaultObject) AddToSetsByName(sets []string) error {
	var ids []string
	for _, v := range sets {
		setObj := NewManualSet(o.client)
		setObj.Name = v
		setObj.ObjectType = o.SetType
		id, err := setObj.GetIDByName()
		if err != nil || id == "" {
			return fmt.Errorf("%s type Set %s doesn't exist", setObj.ObjectType, v)
		}
		ids = append(ids, id)
	}
	err := o.AddToSetsByID(ids)
	if err != nil {
		return err
	}
	return nil
}

/*
Set permissions (Collection/SetCollectionPermissions)
	Request body format
		{
			"Grants": [
				{
					"Principal": "LAB Infrastructure Admins",
					"PType": "Role",
					"Rights": "View,Edit",
					"PrincipalId": "5e104003_eeed_422f_9b45_bca14b61528d"
				},
				{
					"Principal": "LAB MFA Machines & Users",
					"PType": "Role",
					"Rights": "View",
					"PrincipalId": "d06fc516_8c9b_4f76_a08d_797ca6fd0a55"
				},
				{
					"Principal": "LAB Cloud Local Admins",
					"PType": "Role",
					"Rights": "View,Edit",
					"PrincipalId": "d958fad8_f90a_4c40_b986_f6fa31713bba"
				}
			],
			"ID": "e0f8aa18-270a-4bb7-82d9-23afd6a81861",
			"RowKey": "e0f8aa18-270a-4bb7-82d9-23afd6a81861",
			"PVID": "e0f8aa18-270a-4bb7-82d9-23afd6a81861"
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
*/
