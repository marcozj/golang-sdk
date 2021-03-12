package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/cloudprovidertype"
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

	////////////////////////////////////////////
	// Sample code to create a cloud provider //
	////////////////////////////////////////////
	obj := platform.NewCloudProvider(client)
	obj.Name = "Test AWS"                     // Mandatory
	obj.Type = cloudprovidertype.AWS.String() // Mandatory
	obj.CloudAccountID = "548715038142"       // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating cloud provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created cloud provider '%s'\n", obj.Name)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Edit, platform.Right.Delete},
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
		fmt.Printf("Error assign permissions to cloud provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to cloud provider '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"test"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding cloud provider to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added cloud provider %s to Sets '%+v'\n", obj.Name, sets)

	////////////////////////////////////////////
	// Sample code to update a cloud provider //
	////////////////////////////////////////////
	obj = platform.NewCloudProvider(client)
	obj.Name = "Test AWS"               // Mandatory
	obj.CloudAccountID = "548715038142" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test cloud provider"
	obj.EnableUnmanagedPasswordRotation = true
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating cloud provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated cloud provider '%s'\n", obj.Name)

	////////////////////////////////////////////
	// Sample code to delete a cloud provider //
	////////////////////////////////////////////
	obj = platform.NewCloudProvider(client)
	obj.Name = "Test AWS"               // Mandatory
	obj.CloudAccountID = "548715038142" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting cloud provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted cloud provider '%s'\n", obj.Name)

}
