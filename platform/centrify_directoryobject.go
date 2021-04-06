package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/restapi"
)

// DirectoryObjects -
type DirectoryObjects struct {
	client            *restapi.RestClient
	ObjectType        string // Either user or group
	QueryName         string
	DirectoryServices []string          `json:"DirectoryServices,omitempty" schema:"directory_services,omitempty"`
	DirectoryObjects  []DirectoryObject `json:"DirectoryObjects,omitempty" schema:"directory_object,omitempty"`

	apiRead string
}

// DirectoryObject -
type DirectoryObject struct {
	ID                string `json:"InternalName,omitempty" schema:"id,omitempty"`
	RoleID            string `json:"_ID,omitempty" schema:"roleid,omitempty"` // this is only for Centrify Directory role
	Name              string `json:"Name,omitempty" schema:"name,omitempty"`
	SystemName        string `json:"SystemName,omitempty" schema:"system_name,omitempty"`
	DisplayName       string `json:"DisplayName,omitempty" schema:"display_name,omitempty"`
	DistinguishedName string `json:"DistinguishedName,omitempty" schema:"distinguished_name,omitempty"`
	ObjectType        string `json:"ObjectType,omitempty" schema:"object_type,omitempty"`
	Forest            string `json:"Forest,omitempty" schema:"forest,omitempty"`
}

// NewDirectoryObjects is a DirectoryObjects constructor
func NewDirectoryObjects(c *restapi.RestClient) *DirectoryObjects {
	s := DirectoryObjects{}
	s.client = c
	s.apiRead = "/UserMgmt/DirectoryServiceQuery"

	return &s
}

// Read function fetches directory objects from source
func (o *DirectoryObjects) Read() error {
	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	queryArg["directoryServices"] = o.DirectoryServices
	switch o.ObjectType {
	case "User":
		queryArg["user"] = "{\"_and\":[{\"SystemName\":{\"_like\":\"" + o.QueryName + "\"}},{\"ObjectType\":\"User\"}]}"
	case "Group":
		queryArg["group"] = "{\"SystemName\":{\"_like\":\"" + o.QueryName + "\"}}"
	case "Role":
		queryArg["roles"] = "{\"Name\":{\"_like\":\"" + o.QueryName + "\"}}"
	}

	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}

	var rs map[string]interface{}
	switch o.ObjectType {
	case "User":
		rs = resp.Result["User"].(map[string]interface{})
	case "Group":
		rs = resp.Result["Group"].(map[string]interface{})
	case "Role":
		rs = resp.Result["roles"].(map[string]interface{})
	}
	var results = rs["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		obj := &DirectoryObject{}
		mapToStruct(obj, row)
		o.DirectoryObjects = append(o.DirectoryObjects, *obj)
	}

	return nil
}

func (o *DirectoryObjects) GetByName(objType string, name string, dir DirectoryService) (*DirectoryObject, error) {
	o.QueryName = name
	o.ObjectType = objType
	o.DirectoryServices = []string{dir.ID}
	err := o.Read()
	if err != nil {
		return nil, fmt.Errorf("error retrieving directory services: %s", err)
	}

	if len(o.DirectoryObjects) == 0 {
		return nil, fmt.Errorf("query returns 0 object for directory object %s", name)
	}
	if len(o.DirectoryObjects) > 1 {
		return nil, fmt.Errorf("search directory object %s, but returns too many objects (found %d, expected 1)", name, len(o.DirectoryObjects))
	}

	return &o.DirectoryObjects[0], nil
}

/*
https://developer.centrify.com/reference#post_usermgmt-directoryservicequery

Request body
{
    "user": "{\"_and\":[{\"_or\":[{\"DisplayName\":{\"_like\":\"LAB\"}},{\"givenName\":{\"_like\":\"LAB\"}},{\"sn\":{\"_like\":\"LAB\"}},{\"SystemName\":{\"_like\":\"LAB\"}}]},{\"ObjectType\":\"user\"}]}",
    "directoryServices": [
        "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
        "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
        "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
    ],
    "group": "{\"_or\":[{\"DisplayName\":{\"_like\":\"LAB\"}},{\"SystemName\":{\"_like\":\"LAB\"}}]}",
    "roles": "{\"_or\":[{\"_ID\":{\"_like\":\"LAB\"}},{\"Name\":{\"_like\":\"LAB\"}}]}",
    "Args": {
        "PageNumber": 1,
        "PageSize": 100000,
        "Limit": 100000,
        "SortBy": "",
        "direction": "False",
        "Caching": -1
    }
}

Respond result
{
    "success": true,
    "Result": {
        "User": {
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
                            "Type": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                            "Key": "User",
                            "IsForeignKey": false
                        }
                    ],
                    "Row": {
                        "Description": null,
                        "DisplayName": "Centrify AD Admin",
                        "ObjectType": "User",
                        "DistinguishedName": "CN=Centrify AD Admin,OU=Lab Service Accounts,DC=example,DC=com",
                        "DirectoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                        "SystemName": "admin@example.com",
                        "ServiceInstance": "AdProxy_example.com",
                        "Locked": false,
                        "InternalName": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
                        "StatusEnum": "Created",
                        "ServiceInstanceLocalized": "Active Directory (example.com)",
                        "ServiceType": "AdProxy",
                        "Forest": "example.com",
                        "EMail": "admin@example.com",
                        "Status": "Not Invited",
                        "Enabled": true
                    }
                }
            ],
            "ReturnID": ""
        }
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
