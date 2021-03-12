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

	//////////////////////////////////
	// Sample code to create a role //
	//////////////////////////////////
	obj := platform.NewRole(client)
	obj.Name = "Test role" // Mandatory
	obj.Description = "Test role created by SDK"
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating role: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created role '%s'\n", obj.Name)

	// Assign administrative rights
	obj.AdminRights = []string{"Privileged Access Service User", "Report Management"}
	_, err = obj.AssignAdminRights()
	if err != nil {
		fmt.Printf("Error updating role admin rights: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned admin rights %s\n", obj.AdminRights)

	//////////////////////////////////
	// Sample code to update a role //
	//////////////////////////////////
	obj = platform.NewRole(client)
	obj.Name = "Test role" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "Test role created by SDK ..."
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating role: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated role '%s'\n", obj.Name)

	//////////////////////////////////
	// Sample code to delete a role //
	//////////////////////////////////
	obj = platform.NewRole(client)
	obj.Name = "Test role"
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting role: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted role '%s'\n", obj.Name)
}
