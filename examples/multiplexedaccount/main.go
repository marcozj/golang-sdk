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

	/////////////////////////////////////////////
	// Sample code to create multiplex account //
	/////////////////////////////////////////////
	obj := platform.NewMultiplexedAccount(client)
	obj.Name = "Test multiplex account"        // Mandatory
	obj.RealAccount1UPN = "test_svc1@demo.lab" // Mandatory
	obj.RealAccount2UPN = "test_svc2@demo.lab" // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating multiplex account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created multiplex account '%s'\n", obj.Name)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.Edit, platform.Right.Delete},
		},
	}

	err = platform.ResolvePermissions(client, myPermissions, obj.ValidPermissions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Assign resolved permission
	obj.Permissions = myPermissions
	_, err = obj.SetPermissions(false)
	if err != nil {
		fmt.Printf("Error assign permissions to account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to account '%s'\n", obj.Permissions, obj.Name)

	/////////////////////////////////////////////
	// Sample code to update multiplex account //
	/////////////////////////////////////////////
	obj = platform.NewMultiplexedAccount(client)
	obj.Name = "Test multiplex account" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test multiplex account"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating multiplex account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated multiplex account '%s'\n", obj.Name)

	/////////////////////////////////////////////
	// Sample code to delete multiplex account //
	/////////////////////////////////////////////
	obj = platform.NewMultiplexedAccount(client)
	obj.Name = "Test multiplex account" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting multiplex account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted multiplex account '%s'\n", obj.Name)

}
