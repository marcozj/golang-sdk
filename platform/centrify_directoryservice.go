package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/restapi"
)

// DirectoryServices - Encapsulates Directory Services
type DirectoryServices struct {
	client      *restapi.RestClient
	DirServices []DirectoryService `json:"DirServices,omitempty" schema:"directory_service,omitempty"`

	apiRead string
}

// DirectoryService represents directory service
type DirectoryService struct {
	ID               string `json:"directoryServiceUuid,omitempty" schema:"id,omitempty"`
	Name             string `json:"Name,omitempty" schema:"name,omitempty"`
	Description      string `json:"Description,omitempty" schema:"description,omitempty"`
	DisplayName      string `json:"DisplayName,omitempty" schema:"displayName,omitempty"`
	DisplayNameShort string `json:"DisplayNameShort,omitempty" schema:"short_name,omitempty"`
	Service          string `json:"Service,omitempty" schema:"service,omitempty"`
	Status           string `json:"Status,omitempty" schema:"status,omitempty"`
	Config           string `json:"Config,omitempty" schema:"config,omitempty"`
	Forest           string `json:"Forest,omitempty" schema:"forest,omitempty"`
}

// NewDirectoryServices is a DirectoryServices constructor
func NewDirectoryServices(c *restapi.RestClient) *DirectoryServices {
	s := DirectoryServices{}
	s.client = c
	s.apiRead = "/core/GetDirectoryServices"

	return &s
}

// GetDirectorServices etches a DirectorServices from source and returns list of map
func (o *DirectoryServices) GetDirectorServices() ([]map[string]interface{}, error) {
	var dirs []map[string]interface{}

	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf(resp.Message)
	}

	var results = resp.Result["Results"].([]interface{})
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		dirs = append(dirs, row)
	}

	return dirs, nil
}

// Read function fetches a DirectorServices from source
func (o *DirectoryServices) Read() error {
	dirs, err := o.GetDirectorServices()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		obj := &DirectoryService{}
		mapToStruct(obj, dir)
		o.DirServices = append(o.DirServices, *obj)
	}
	//logger.Debugf("Filled object: %+v", o)

	return nil
}

/*
https://developer.centrify.com/reference#post_core-getdirectoryservices

{
    "success": true,
    "Result": {
        "IsAggregate": false,
        "Count": 6,
        ......
        "FullCount": 6,
        "Results": [
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "CDS",
                    "DisplayName": "Centrify Centrify Directory",
                    "Tenant": "ABC0751",
                    "Name": "CDS",
                    "Status": "Active",
                    "Config": "Centrify Directory",
                    "StatusDisplay": "Online",
                    "Everybody": true,
                    "Description": "Centrify Directory",
                    "DisplayNameShort": "Centrify Directory",
                    "directoryServiceUuid": "09B9A9B0-6CE8-465F-AB03-65766D33B05E"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "AdProxy",
                    "DisplayName": "Active Directory: demo.lab",
                    "Tenant": "ABC0751",
                    "Name": "AdProxy_demo.lab",
                    "Status": "Active",
                    "Config": "demo.lab",
                    "StatusDisplay": "Online",
                    "Everybody": true,
                    "Description": "Active Directory",
                    "DisplayNameShort": "AD: demo.lab",
                    "directoryServiceUuid": "e09def65-17c3-0f40-c475-a6ee8825611f",
                    "Forest": "demo.lab"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "AdProxy",
                    "DisplayName": "Active Directory: centrifylab.aws",
                    "Tenant": "ABC0751",
                    "Name": "AdProxy_centrifylab.aws",
                    "Status": "Inactive",
                    "Config": "centrifylab.aws",
                    "StatusDisplay": "Offline",
                    "Everybody": true,
                    "Description": "Active Directory",
                    "DisplayNameShort": "AD: centrifylab.aws",
                    "directoryServiceUuid": "e3344ea3-429c-c901-8854-1a72ea1404b9",
                    "Forest": "centrifylab.aws"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "GDS",
                    "DisplayName": "Google Directory: democorp.club",
                    "Tenant": "ABC0751",
                    "Name": "GDS_c0262k9n4",
                    "Status": "Active",
                    "Config": "democorp.club",
                    "StatusDisplay": "Online",
                    "Everybody": true,
                    "Description": "Google Directory",
                    "DisplayNameShort": "GDS: democorp.club",
                    "directoryServiceUuid": "b3c443a6-1186-4ff1-99f0-8ad8bd6104f8"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "AdProxy",
                    "DisplayName": "Active Directory: centrify.lab",
                    "Tenant": "ABC0751",
                    "Name": "AdProxy_centrify.lab",
                    "Status": "Inactive",
                    "Config": "centrify.lab",
                    "StatusDisplay": "Offline",
                    "Everybody": true,
                    "Description": "Active Directory",
                    "DisplayNameShort": "AD: centrify.lab",
                    "directoryServiceUuid": "d6126e54-88fc-3e9e-8681-34f51b856a9f",
                    "Forest": "centrify.lab"
                }
            },
            {
                "Entities": [
                    {
                        "Type": "DirectoryServices",
                        "Key": "??",
                        "IsForeignKey": false
                    }
                ],
                "Row": {
                    "Service": "FDS",
                    "DisplayName": "Federated Directory Service",
                    "Tenant": "ABC0751",
                    "Name": "FDS",
                    "Status": "Active",
                    "Config": "Federated Directory Service",
                    "StatusDisplay": "Online",
                    "Everybody": true,
                    "Description": "Federated Directory",
                    "DisplayNameShort": "FDS",
                    "directoryServiceUuid": "C30B30B1-0B46-49AC-8D99-F6279EED7999"
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
*/
