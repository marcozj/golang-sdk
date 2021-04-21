package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// PolicyLinks - Encapsulates policy links
type PolicyLinks struct {
	Plinks []PolicyLink `json:"Plinks,omitempty" schema:"policy_order,omitempty"`

	apiRead   string
	apiUpdate string
	client    *restapi.RestClient
}

// PolicyLink - encapsulates policy
type PolicyLink struct {
	ID              string   `json:"ID,omitempty" schema:"id,omitempty"`
	Description     string   `json:"Description,omitempty" schema:"description,omitempty"`
	EnableCompliant bool     `json:"EnableCompliant,omitempty" schema:"enable_compliant,omitempty"`
	LinkType        string   `json:"LinkType,omitempty" schema:"link_type,omitempty"` // Global, Role, Collection, Inactive
	PolicySet       string   `json:"PolicySet,omitempty" schema:"policy_set,omitempty"`
	Params          []string `json:"Params,omitempty" schema:"policy_assignment,omitempty"` // Policy assignment to role or set
}

// NewPolicyLinks is a policy link constructor
func NewPolicyLinks(c *restapi.RestClient) *PolicyLinks {
	s := PolicyLinks{}
	//s.Plinks = []PolicyLink{}
	s.client = c
	s.apiRead = "/Policy/GetNicePlinks"
	s.apiUpdate = "/Policy/setPlinksv2"

	return &s
}

// GetPlinks fetches PolicyLinks from Centrify tenant and return in map format
func (o *PolicyLinks) GetPlinks() ([]map[string]interface{}, string, error) {
	var plinks []map[string]interface{}

	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, "", err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, "", fmt.Errorf(errmsg)
	}

	var rev = resp.Result["RevStamp"].(string)
	var results = resp.Result["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		plinks = append(plinks, row)
	}

	return plinks, rev, nil
}

// Read function fetches a PolicyLinks from source
func (o *PolicyLinks) Read() error {
	plinks, _, err := o.GetPlinks()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	links := []PolicyLink{}
	for _, plink := range plinks {
		obj := &PolicyLink{}
		mapToStruct(obj, plink)
		links = append(links, *obj)
	}
	o.Plinks = links

	return nil
}

// Update function updates an existing PolicyLinks and returns a map that contains update result
func (o *PolicyLinks) Update() (*restapi.GenericMapResponse, error) {
	oldplinks, rev, err := o.GetPlinks()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	// Only change plinks order, not insert or delete any from the list
	if len(o.Plinks) != len(oldplinks) {
		return nil, fmt.Errorf("There are %d defined polices but there are %d existing policies in the tenant", len(o.Plinks), len(oldplinks))
	}

	var newplinks []map[string]interface{}
	for _, v := range o.Plinks {
		found := findItem("ID", v.ID, oldplinks)
		if found == nil {
			// Can't find a matched ID in tenant plinks, return error
			return nil, fmt.Errorf("Policy %s not found in policy list", v.ID)
		}
		newplinks = append(newplinks, found)
	}

	var queryArg = make(map[string]interface{})
	queryArg["Plinks"] = newplinks
	queryArg["RevStamp"] = rev

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

func findItem(key string, value string, items []map[string]interface{}) map[string]interface{} {
	for _, v := range items {
		if v[key] == value {
			return v
		}
	}
	return nil
}

/*
Get PLinks
{
    "success": true,
    "Result": {
        "Columns": [
            ...
        ],
        "RevStamp": "637336119080000000",
        "Count": 2,
        "Results": [
            {
                "Entities": [
                    {
                        "Type": "PolicyLink",
                        "Key": "/Policy/Default Policy",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Params": [],
                    "ID": "/Policy/Default Policy",
                    "EnableCompliant": true,
                    "I18NDescriptionTag": "_I18N_DefaultGlobalPolicyDescriptionTag",
                    "Description": "Default Policy Settings.",
                    "LinkType": "Inactive",
                    "PolicySet": "/Policy/Default Policy"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "PolicyLink",
                        "Key": "/Policy/LAB Deny Login Policy",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Params": [],
                    "ID": "/Policy/LAB Deny Login Policy",
                    "EnableCompliant": true,
                    "Description": "Catch all policy that denies users who don't have PAS role from logging in. This policy must be placed at the bottom of policy list.",
                    "LinkType": "Global",
                    "PolicySet": "/Policy/LAB Deny Login Policy"
                }
            }
        ],
        "FullCount": 2,
        "ReturnID": "",
        "IsAggregate": false
    },
    "Message": null,
    "MessageID": null,
    "Exception": null,
    "ErrorID": null,
    "ErrorCode": null,
    "IsSoftError": false,
    "InnerExceptions": null
}

Set PLinks
https://developer.centrify.com/reference#post_policy-setplinksv2

	Request body format
	{
		"Plinks": [
			{
				"Params": [
					"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
				],
				"ID": "/Policy/Invited Users",
				"EnableCompliant": true,
				"I18NDescriptionTag": "_I18N_InviteUser_Policy_Description",
				"Description": "This policy is created as part of inviting users action. It allows invited users to enroll device through invitation links sent to them.",
				"LinkType": "Role",
				"PolicySet": "/Policy/Invited Users",
				"Name": "Invited Users"
			},
			{
				"Params": [],
				"ID": "/Policy/PolicySet_1",
				"EnableCompliant": true,
				"Description": "",
				"LinkType": "Global",
				"PolicySet": "/Policy/PolicySet_1",
				"Name": "PolicySet_1"
			},
			{
				"Params": [],
				"ID": "/Policy/Default Policy",
				"EnableCompliant": true,
				"I18NDescriptionTag": "_I18N_DefaultGlobalPolicyDescriptionTag",
				"Description": "Default Policy Settings.",
				"LinkType": "Inactive",
				"PolicySet": "/Policy/Default Policy",
				"Name": "Default Policy"
			}
		],
		"RevStamp": "637336119080000000"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"RevStamp": "637334855480000000"
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
