package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/databaseclass"
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
	// Sample code to create a database //
	//////////////////////////////////////
	obj := platform.NewDatabase(client)
	obj.Name = "Test database"                           // Mandatory
	obj.DatabaseClass = databaseclass.SQLServer.String() // Mandatory
	obj.FQDN = "test.example.com"                        // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created database '%s'\n", obj.Name)

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
		fmt.Printf("Error assign permissions to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to database '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"LAB_Databases"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding database to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added database %s to Sets '%+v'\n", obj.Name, sets)

	//////////////////////////////////////
	// Sample code to update a database //
	//////////////////////////////////////
	obj = platform.NewDatabase(client)
	obj.Name = "Test database"                           // Mandatory
	obj.DatabaseClass = databaseclass.SQLServer.String() // Mandatory
	obj.FQDN = "test.example.com"                        // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test database - updated"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated database '%s'\n", obj.Name)

	//////////////////////////////////////
	// Sample code to delete a database //
	//////////////////////////////////////
	obj = platform.NewDatabase(client)
	obj.Name = "Test database"                           // Mandatory
	obj.DatabaseClass = databaseclass.SQLServer.String() // Mandatory
	obj.FQDN = "test.example.com"                        // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting database: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted database '%s'\n", obj.Name)
}
