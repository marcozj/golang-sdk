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

	/////////////////////////////////////////////////////
	// Sample code to create a Centrify Directory user //
	/////////////////////////////////////////////////////
	obj := platform.NewUser(client)
	obj.Name = "testuser@<suffix>"
	obj.Mail = "testuser@eexample.com"
	obj.DisplayName = "Test User"
	obj.Password = "xxxxxxxx"
	obj.ConfirmPassword = "xxxxxxxx"
	obj.PasswordNeverExpire = true
	obj.ForcePasswordChangeNext = false
	obj.Description = "test user created by SDK"
	obj.OfficeNumber = "12345678"
	obj.HomeNumber = "12345678"
	obj.MobileNumber = "12345678"

	obj.Roles = []string{"Role Name 1", "Role Name 2"}
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created user '%s'\n", obj.Name)

	// Update role membership
	err = obj.AddToRoles(obj.Roles)
	if err != nil {
		fmt.Printf("Error adding user %s to role %s. %v\n", obj.Name, obj.Roles, err)
	}
	fmt.Printf("Added user '%s' to role %s\n", obj.Name, obj.Roles)

	/////////////////////////////////////////////////////
	// Sample code to update a Centrify Directory user //
	/////////////////////////////////////////////////////
	obj = platform.NewUser(client)
	obj.Name = "testuser@<suffix>"
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.MobileNumber = "87654321"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated user '%s'\n", obj.Name)

	//////////////////////////////////////////////////////////////
	// Sample code to change a Centrify Directory user password //
	//////////////////////////////////////////////////////////////
	obj = platform.NewUser(client)
	obj.Name = "testuser@<suffix>"
	err = obj.ChangeUserPassword("xxxxxxxxx")
	if err != nil {
		fmt.Printf("Error changing user password: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Changed password for user '%s'\n", obj.Name)

	/////////////////////////////////////////////////////
	// Sample code to delete a Centrify Directory user //
	/////////////////////////////////////////////////////
	obj = platform.NewUser(client)
	obj.Name = "testuser@<suffix>"
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted user '%s'\n", obj.Name)

}
