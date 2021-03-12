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

	//////////////////////////////////////////////
	// Sample code to create a password profile //
	//////////////////////////////////////////////
	obj := platform.NewPasswordProfile(client)
	obj.Name = "Test Password Profile" // Mandatory
	obj.MinimumPasswordLength = 9      // Mandatory
	obj.MaximumPasswordLength = 16     // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating password profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created password profile '%s'\n", obj.Name)

	//////////////////////////////////////////////
	// Sample code to update a password profile //
	//////////////////////////////////////////////
	obj = platform.NewPasswordProfile(client)
	obj.Name = "Test Password Profile" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.AtLeastOneSpecial = true
	obj.SpecialCharSet = "!$%&()*+,-./:;<=>?[\\]^_{|}~"
	obj.MaximumCharOccurrenceCount = 2
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating password profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated password profile '%s'\n", obj.Name)

	//////////////////////////////////////////////
	// Sample code to delete a password profile //
	//////////////////////////////////////////////
	obj = platform.NewPasswordProfile(client)
	obj.Name = "Test Password Profile" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting password profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted password profile '%s'\n", obj.Name)

}
