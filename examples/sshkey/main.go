package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/marcozj/golang-sdk/enum/keypairtype"
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
	// Sample code to create a SSHKey //
	////////////////////////////////////
	obj := platform.NewSSHKey(client)
	obj.Name = "Test key" // Mandatory
	key, err := ioutil.ReadFile("testkey.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	obj.PrivateKey = string(key) // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating sshkey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created sshkey '%s'\n", obj.Name)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Retrieve, platform.Right.Edit, platform.Right.Delete},
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
		fmt.Printf("Error assign permissions to sshkey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to sshkey '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"SSHKey Set"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding sshkey to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added sshkey %s to Sets '%+v'\n", obj.Name, sets)

	////////////////////////////////////
	// Sample code to update a sshkey //
	////////////////////////////////////
	obj = platform.NewSSHKey(client)
	obj.Name = "Test key" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test sshkey"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating sshkey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated sshkey '%s'\n", obj.Name)

	//////////////////////////////////////
	// Sample code to retrieve a sshkey //
	//////////////////////////////////////
	obj = platform.NewSSHKey(client)
	obj.Name = "Test key" // Mandatory
	obj.KeyPairType = keypairtype.PrivateKey.String()
	//obj.Passphrase = ""
	var mykey string
	mykey, err = obj.RetriveSSHKey()
	if err != nil {
		fmt.Printf("Error retrieve sshkey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved sshkey '%s'\n", mykey)

	////////////////////////////////////
	// Sample code to delete a sshkey //
	////////////////////////////////////
	obj = platform.NewSSHKey(client)
	obj.Name = "Test key" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting sshkey: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted sshkey '%s'\n", obj.Name)

}
