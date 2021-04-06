package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/directoryservice"
	"github.com/marcozj/golang-sdk/enum/secrettype"
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
	// Sample code to create a secret //
	////////////////////////////////////
	obj := platform.NewSecret(client)
	obj.SecretName = "Test secret"                   // Mandatory
	obj.SecretText = "kjfljakljklajsdlkfjklasjfdlkj" // Mandatory
	obj.Type = secrettype.Text.String()              // Mandatory
	obj.ParentPath = "folder1\\folder2"

	// Assign workflow
	obj.WorkflowEnabled = true
	obj.WorkflowApprovers = []platform.WorkflowApprover{
		{
			Type:            "Manager",
			OptionsSelector: true,
			NoManagerAction: "useBackup",
			BackupApprover: &platform.BackupApprover{
				Name:             "labadmin@demo.lab",
				Type:             "User",
				DirectoryService: directoryservice.ActiveDirectory.String(),
				DirectoryName:    "demo.lab",
			},
		},
		{
			Name:             "LAB Infrastructure Owners",
			Type:             "Role",
			DirectoryService: directoryservice.CentrifyDirectory.String(),
			DirectoryName:    "Centrify Directory",
		},
	}

	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created secret '%s'\n", obj.SecretName)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.RetrieveSecret, platform.Right.Edit, platform.Right.Delete},
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
		fmt.Printf("Error assign permissions to secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to secret '%s'\n", obj.Permissions, obj.SecretName)

	// Assign to Sets
	sets := []string{"Test Secrets"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding secret to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added secret %s to Sets '%+v'\n", obj.SecretName, sets)

	////////////////////////////////////
	// Sample code to update a secret //
	////////////////////////////////////
	obj = platform.NewSecret(client)
	obj.SecretName = "Test secret"      // Mandatory
	obj.ParentPath = "folder1\\folder2" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test secret"
	obj.NewParentPath = "folder1" // Move to another folder
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated secret '%s'\n", obj.SecretName)

	//////////////////////////////////////
	// Sample code to retrieve a secret //
	//////////////////////////////////////
	obj = platform.NewSecret(client)
	obj.SecretName = "Test secret" // Mandatory
	obj.ParentPath = "folder1"     // Mandatory
	var mykey string
	mykey, err = obj.CheckoutSecret()
	if err != nil {
		fmt.Printf("Error retrieve secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved secret '%s'\n", mykey)

	////////////////////////////////////
	// Sample code to delete a secret //
	////////////////////////////////////
	obj = platform.NewSecret(client)
	obj.SecretName = "Test secret" // Mandatory
	obj.ParentPath = "folder1"     // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted secret '%s'\n", obj.SecretName)

	///////////////////////////////////////////
	// Sample code to download a secret file //
	///////////////////////////////////////////
	obj = platform.NewSecret(client)
	obj.SecretName = "Alice"       // Mandatory
	obj.ParentPath = "Certificate" // Mandatory

	var secretfile string
	// Download secret file and save it in current directory
	// If input parameter is true, the downloaded file is save in user's home directory
	secretfile, err = obj.CheckoutSecretAndFile(false)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("Downloaded file %s\n", secretfile)

}
