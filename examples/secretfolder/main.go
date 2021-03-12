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

	///////////////////////////////////////////
	// Sample code to create a secret folder //
	///////////////////////////////////////////
	obj := platform.NewSecretFolder(client)
	obj.Name = "testfolder"             // Mandatory
	obj.ParentPath = "folder1\\folder2" // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating secret folder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created secret folder '%s'\n", obj.Name)

	// Assign folder permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Edit, platform.Right.Delete, platform.Right.Add},
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
		fmt.Printf("Error assign permissions to secret folder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to secret folder '%s'\n", obj.Permissions, obj.Name)

	// Assign member permissions
	memberPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Edit, platform.Right.Delete, platform.Right.RetrieveSecret},
		},
	}

	err = platform.ResolvePermissions(client, memberPermissions, obj.ValidMemberPermissions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Assign resolved permission
	obj.MemberPermissions = memberPermissions
	_, err = obj.SetMemberPermissions(false)
	if err != nil {
		fmt.Printf("Error assign member permissions to secret folder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned member permissions %+v to secret folder '%s'\n", obj.MemberPermissions, obj.Name)

	///////////////////////////////////////////
	// Sample code to update a secret folder //
	///////////////////////////////////////////
	obj = platform.NewSecretFolder(client)
	obj.Name = "testfolder"             // Mandatory
	obj.ParentPath = "folder1\\folder2" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test secret folder"

	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating secret folder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated secret folder '%s'\n", obj.Name)

	obj.NewParentPath = "folder1" // Move to another folder
	obj.MoveFolder()
	if err != nil {
		fmt.Printf("Error updating secret folder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Moved secret folder '%s'\n", obj.Name)

	///////////////////////////////////////////
	// Sample code to delete a secret folder //
	///////////////////////////////////////////
	obj = platform.NewSecretFolder(client)
	obj.Name = "testfolder"    // Mandatory
	obj.ParentPath = "folder1" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting secret folder: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted secret folder '%s'\n", obj.Name)

}
