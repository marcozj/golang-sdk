package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/directoryservice"
	"github.com/marcozj/golang-sdk/enum/webapp/saml/applicationtemplate"
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
	obj := platform.NewSamlWebApp(client)
	obj.Name = "Test SAML WebApp"                        // Mandatory
	obj.TemplateName = applicationtemplate.SAML.String() // Mandatory

	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s '%s'\n", platform.GetVarType(obj), obj.Name)

	// Many of parameters can only be set in update process rather than initinal creation
	obj.SpConfigMethod = 1 // 0 indicates manual configruation, 1 indicates metadata
	obj.SpMetadataUrl = "https://nexus.microsoftonline-p.com/federationmetadata/saml20/federationmetadata.xml"
	/*
		metaxml, err := ioutil.ReadFile("sp_meta.xml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		obj.SpMetadataXml = string(metaxml)
	*/
	/*
		obj.Audience = "urn:federation:MicrosoftOnline"
		obj.ACS_Url = "https://login.microsoftonline.com/login.srf"
		obj.RecipientSameAsAcsUrl = false
		obj.Recipient = "https://login.microsoftonline.com/login.srf"
		obj.WantAssertionsSigned = true
		obj.NameIDFormat = "emailAddress"
		obj.SpSingleLogoutUrl = "https://login.microsoftonline.com/logout.srf"
		obj.RelayState = "state"
		obj.AuthnContextClass = "X509"
	*/
	obj.SamlAttributes = []platform.SamlAttribute{
		{
			Name:  "attribute1",
			Value: "value1",
		},
		{
			Name:  "attribute2",
			Value: "value2",
		},
	}
	obj.UserNameStrategy = "ADAttribute"
	obj.Username = "userprincipalname"

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

	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated webapp '%s'\n", obj.Name)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Run},
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
		fmt.Printf("Error assign permissions to webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to webapp '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"Test Web Apps"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding webapp to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added webapp %s to Sets '%+v'\n", obj.Name, sets)

	////////////////////////////////////
	// Sample code to update a webapp //
	////////////////////////////////////
	obj = platform.NewSamlWebApp(client)
	obj.Name = "Test SAML WebApp" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test SAML webapp"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("IdP meta url: %s\n", obj.IdpMetadataUrl)
	fmt.Printf("Updated webapp '%s'\n", obj.Name)

	////////////////////////////////////
	// Sample code to delete a webapp //
	////////////////////////////////////
	obj = platform.NewSamlWebApp(client)
	obj.Name = "Test SAML WebApp" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted webapp '%s'\n", obj.Name)

}
