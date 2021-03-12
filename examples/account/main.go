package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/resourcetype"
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

	//////////////////////////////////////
	// Sample code to create an account //
	//////////////////////////////////////
	obj := platform.NewAccount(client)
	obj.User = "testaccount"                        // Mandatory
	obj.ResourceName = "centos1"                    // Mandatory
	obj.ResourceType = resourcetype.System.String() // Mandatory
	obj.Password = "kjlaxaHHJKJd*^234"              // Mandatory

	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created account '%s'\n", obj.User)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Checkout, platform.Right.Login, platform.Right.FileTransfer, platform.Right.Edit, platform.Right.Delete, platform.Right.UpdatePassword, platform.Right.WorkspaceLogin, platform.Right.RotatePassword},
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
	fmt.Printf("Assigned permissions %+v to account '%s'\n", obj.Permissions, obj.User)

	// Assign to Sets
	sets := []string{"Test Accounts"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding account to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added account %s to Sets '%+v'\n", obj.User, sets)

	//////////////////////////////////////
	// Sample code to update an account //
	//////////////////////////////////////
	obj = platform.NewAccount(client)
	obj.User = "testaccount"                        // Mandatory
	obj.ResourceName = "centos1"                    // Mandatory
	obj.ResourceType = resourcetype.System.String() // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test account - updated"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated account '%s'\n", obj.User)

	////////////////////////////////////////////////////
	// Sample code to checkout password of an account //
	////////////////////////////////////////////////////
	obj = platform.NewAccount(client)
	obj.User = "testaccount"                        // Mandatory
	obj.ResourceName = "centos1"                    // Mandatory
	obj.ResourceType = resourcetype.System.String() // Mandatory
	pwd, err := obj.CheckoutPassword(false)
	if err != nil {
		fmt.Printf("Error checkout password of account: %v\n", err)
		os.Exit(1)
	}
	if pwd != "" {
		fmt.Printf("Password for account %s in %s is: %s\n", obj.User, obj.ResourceName, pwd)
		if err != nil {
			fmt.Printf("Password checkin error: %+v\n", err)
		}
	}

	//////////////////////////////////////
	// Sample code to delete an account //
	//////////////////////////////////////
	obj = platform.NewAccount(client)
	obj.User = "testaccount"                        // Mandatory
	obj.ResourceName = "centos1"                    // Mandatory
	obj.ResourceType = resourcetype.System.String() // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting account: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted account '%s'\n", obj.User)

}
