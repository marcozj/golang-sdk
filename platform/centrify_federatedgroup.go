package platform

import (
	"fmt"

	"github.com/marcozj/golang-sdk/enum/directoryservice"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// FederatedGroup - Encapsulates a single Federated Group
type FederatedGroup struct {
	// Rest client
	client *restapi.RestClient
	// Standard attributes
	ID   string `json:"InternalName,omitempty" schema:"id,omitempty"`
	Name string `json:"SystemName,omitempty" schema:"name,omitempty"`

	// API endpoints
	apiRead   string //`json:"-"` // Ignoring this JSON field
	apiCreate string //`json:"-"` // Ignoring this JSON field
	// There isn't API to delete federated group. This API here is for deleting global group mapping instead.
	apiDelete string //`json:"-"` // Ignoring this JSON field
}

// NewFederatedGroup is a FederatedGroup constructor
func NewFederatedGroup(c *restapi.RestClient) *FederatedGroup {
	s := FederatedGroup{}
	s.client = c
	s.apiRead = "/UserMgmt/DirectoryServiceQuery"
	s.apiCreate = "/Federation/AddGlobalGroupAssertionMapping"
	s.apiDelete = "/Federation/DeleteGlobalGroupAssertionMapping"

	return &s
}

// Read function fetches a FederatedGroup from source, including attribute values. Returns error if any
func (o *FederatedGroup) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	// Get Federated Directory Service
	ds := NewDirectoryServices(o.client)
	theds, err := ds.GetByName(directoryservice.FederatedDirectory.String(), "Federated Directory Service")
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	// Get federated group
	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	queryArg["directoryServices"] = []string{theds.ID}
	queryArg["group"] = "{\"InternalName\":{\"_like\":\"" + o.ID + "\"}}"
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

	rs := resp.Result["Group"].(map[string]interface{})
	var results = rs["Results"].([]interface{})
	if len(results) != 1 {
		return fmt.Errorf("Query didn't return exactly 1 group (found %d)", len(results))
	}
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		mapToStruct(o, row)
	}

	return nil
}

// Create function creates a new FederatedGroup and returns a map that contains creation result
func (o *FederatedGroup) Create() (string, error) {
	/*
		// There isn't API call to create federated group so use global group mapping API instead
		var queryArg = make(map[string]interface{})
		queryArg["TempAttributeNameFromSDK_DoNotUse"] = o.Name
		resp, err := o.client.CallStringAPI(o.apiCreate, queryArg)
		if err != nil {
			logger.Errorf(err.Error())
			return "", err
		}
		if !resp.Success {
			errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
			logger.Errorf(errmsg)
			return "", fmt.Errorf(errmsg)
		}
	*/
	// After successful creation of global group mapping, get group ID
	fedgrp := NewFederatedGroup(o.client)
	fedgrp.Name = o.Name
	id, err := fedgrp.GetIDByName()
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	/*
		// Since we create a global group mapping in order to create the federated group, we need to delete the mapping
		resp, err = o.client.CallStringAPI(o.apiDelete, queryArg)
		if err != nil {
			logger.Errorf(err.Error())
			return "", err
		}
		if !resp.Success {
			errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
			logger.Errorf(errmsg)
			return "", fmt.Errorf(errmsg)
		}
	*/
	return id, nil
}

func (o *FederatedGroup) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("group name must be provided")
	}

	// Get Federated Directory Service
	ds := NewDirectoryServices(o.client)
	theds, err := ds.GetByName(directoryservice.FederatedDirectory.String(), "Federated Directory Service")
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}

	// Get federated group
	var queryArg = make(map[string]interface{})
	queryArg["Args"] = subArgs
	queryArg["directoryServices"] = []string{theds.ID}
	queryArg["group"] = "{\"SystemName\":{\"_like\":\"" + o.Name + "\"}}"
	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return "", fmt.Errorf(errmsg)
	}

	rs := resp.Result["Group"].(map[string]interface{})
	var results = rs["Results"].([]interface{})
	if len(results) == 0 {
		return "", fmt.Errorf("Query returns 0 federated group")
	}
	// There could be more than one groups returned because query uses "like" operator
	var row map[string]interface{}
	for _, result := range results {
		row = result.(map[string]interface{})["Row"].(map[string]interface{})
		if row["SystemName"] == o.Name {
			o.ID = row["InternalName"].(string)
			return o.ID, nil
		}
	}

	return "", fmt.Errorf("unknown problem getting ID of federated group %s", o.Name)
}

func (o *FederatedGroup) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return fmt.Errorf("Failed to find ID of federated group %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}
