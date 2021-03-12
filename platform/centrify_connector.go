package platform

import (
	"github.com/marcozj/golang-sdk/restapi"
)

// Connector - Encapsulates a single Connector
type Connector struct {
	vaultObject
}

// NewConnector is a Connector constructor
func NewConnector(c *restapi.RestClient) *Connector {
	s := Connector{}
	s.client = c

	return &s
}

// Read function fetches a ManaulSet from source, including attribute values. Returns error if any
func (o *Connector) Read() error {
	return nil
}

// Delete function deletes a Connector and returns a map that contains deletion result
func (o *Connector) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing Connector and returns a map that contains update result
func (o *Connector) Update() (*restapi.GenericMapResponse, error) {
	return nil, nil
}

// Query function returns a single Connector object in map format
func (o *Connector) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Proxy WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}
