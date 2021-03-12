package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/examples"
	"github.com/marcozj/golang-sdk/platform"
)

func main() {
	// Authenticate and returns authenticated REST client
	client, err := examples.GetClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	///////////////////////////////////////
	// Sample code to add group mappings //
	///////////////////////////////////////
	obj := platform.NewGroupMappings(client)
	obj.Mappings = []platform.GroupMapping{
		platform.GroupMapping{
			AttributeValue: "Test 1",
			GroupName:      "Okta PAS Admin",
		},
		platform.GroupMapping{
			AttributeValue: "Test 2",
			GroupName:      "Azure PAS Admin",
		},
	}
	err = obj.Create()
	if err != nil {
		fmt.Printf("Error adding group mappings: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added group mappings '%+v'\n", obj.Mappings)

	//////////////////////////////////////////
	// Sample code to delete group mappings //
	//////////////////////////////////////////
	obj = platform.NewGroupMappings(client)
	obj.Mappings = []platform.GroupMapping{
		platform.GroupMapping{
			AttributeValue: "Test 1",
			GroupName:      "Okta PAS Admin",
		},
		platform.GroupMapping{
			AttributeValue: "Test 2",
			GroupName:      "Azure PAS Admin",
		},
	}
	err = obj.Delete()
	if err != nil {
		fmt.Printf("Error deleting group mappings: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted group mappings %+v'\n", obj.Mappings)

	/////////////////////////////////////////////////////
	// Sample code to retrieve existing group mappings //
	/////////////////////////////////////////////////////
	obj = platform.NewGroupMappings(client)
	err = obj.Read()
	if err != nil {
		fmt.Printf("Error reading group mappings: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read group mappings %+v'\n", obj.Mappings)

}
