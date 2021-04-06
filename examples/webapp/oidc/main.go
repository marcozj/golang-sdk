package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/directoryservice"
	"github.com/marcozj/golang-sdk/enum/webapp/accountmapping"
	"github.com/marcozj/golang-sdk/enum/webapp/oidc/applicationtemplate"
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
	// Sample code to create a webapp //
	////////////////////////////////////
	obj := platform.NewOidcWebApp(client)
	obj.Name = "Test OIDC WebApp"                           // Mandatory
	obj.TemplateName = applicationtemplate.Generic.String() // Mandatory
	obj.ApplicationID = "TestOIDCClient"                    // Mandatory No space is allowed
	obj.OAuthProfile = &platform.OidcProfile{
		ClientSecret:    "kljalksjdfla",
		Url:             "https://example.com",
		Redirects:       []string{"https://example.com", "https://test.com"},
		TokenLifetime:   "6:00:00",
		AllowRefresh:    true,
		RefreshLifetime: "200.00:00:00",
	}
	obj.UserNameStrategy = accountmapping.SharedAccount.String()
	obj.Username = "sharedaccount"
	// Assign workflow
	obj.WorkflowEnabled = true
	obj.WorkflowApproverList = []platform.WorkflowApprover{
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
	obj.Permissions = []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Run},
		},
	}
	obj.Sets = []string{"Test Web Apps"}

	err = obj.CreateComplete()
	if err != nil {
		fmt.Printf("Error creating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s '%s'\n", platform.GetVarType(obj), obj.Name)

	////////////////////////////////////
	// Sample code to update a webapp //
	////////////////////////////////////
	obj = platform.NewOidcWebApp(client)
	obj.Name = "Test OIDC WebApp" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test OIDC webapp"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated webapp '%s'\n", obj.Name)

	////////////////////////////////////
	// Sample code to delete a webapp //
	////////////////////////////////////
	obj = platform.NewOidcWebApp(client)
	obj.Name = "Test OIDC WebApp" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted webapp '%s'\n", obj.Name)

}
