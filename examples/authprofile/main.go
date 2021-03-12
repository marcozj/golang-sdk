package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/authmechanism"
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

	////////////////////////////////////////////////////
	// Sample code to create an authentication profile //
	////////////////////////////////////////////////////
	obj := platform.NewAuthenticationProfile(client)
	obj.Name = "Test MFA Profile" // Mandatory
	obj.Challenge1 = []string{authmechanism.Password.String(), authmechanism.MobileAuthenticator.String()}
	obj.Challenge2 = []string{authmechanism.OATH_OTP.String(), authmechanism.SecurityQuestions.String()}
	obj.DurationInMinutes = 30
	obj.NumberOfQuestions = 2
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating authentication profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created authentication profile '%s'\n", obj.Name)

	/////////////////////////////////////////////////////
	// Sample code to update an authentication profile //
	/////////////////////////////////////////////////////
	obj = platform.NewAuthenticationProfile(client)
	obj.Name = "Test MFA Profile" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Challenge1 = []string{authmechanism.EmailConfirmationCode.String(), authmechanism.FIDO2.String(), authmechanism.Radius.String()}
	obj.DurationInMinutes = 20
	obj.NumberOfQuestions = 1
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating authentication profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated authentication profile '%s'\n", obj.Name)

	/////////////////////////////////////////////////////
	// Sample code to delete an authentication profile //
	/////////////////////////////////////////////////////
	obj = platform.NewAuthenticationProfile(client)
	obj.Name = "Test MFA Profile" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting authentication profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted authentication profile '%s'\n", obj.Name)

}
