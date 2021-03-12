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

	////////////////////////////////////
	// Sample code to create a domain //
	////////////////////////////////////
	obj := platform.NewDomain(client)
	obj.Name = "example.com" // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating domain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created domain '%s'\n", obj.Name)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "admin@centrify.com.207",
			PrincipalType: "User",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Edit, platform.Right.Delete, platform.Right.UnlockAccount, platform.Right.AddAccount},
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
		fmt.Printf("Error assign permissions to domain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to domain '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"Test Set"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding domain to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added domain %s to Sets '%+v'\n", obj.Name, sets)

	////////////////////////////////////
	// Sample code to update a domain //
	////////////////////////////////////
	obj = platform.NewDomain(client)
	obj.Name = "example.com" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test domain"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating domain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated domain '%s'\n", obj.Name)

	////////////////////////////////////
	// Sample code to delete a domain //
	////////////////////////////////////
	obj = platform.NewDomain(client)
	obj.Name = "example.com" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting domain: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted domain '%s'\n", obj.Name)
}
