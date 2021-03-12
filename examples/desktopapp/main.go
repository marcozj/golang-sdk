package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/desktopapp/applicationtemplate"
	"github.com/marcozj/golang-sdk/enum/desktopapp/cmdparamtype"
	"github.com/marcozj/golang-sdk/enum/desktopapp/logincredential"
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

	////////////////////////////////////////
	// Sample code to create a desktopapp //
	////////////////////////////////////////
	obj := platform.NewDesktopApp(client)
	obj.Name = "Test DesktopApp"                            // Mandatory
	obj.TemplateName = applicationtemplate.Generic.String() // Mandatory
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating desktopapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s '%s'\n", platform.GetVarType(obj), obj.Name)

	// Many of parameters can only be set in update process rather than initinal creation
	obj.DesktopAppProgramName = "test_app" // Mandatory
	obj.DesktopAppRunHostName = "member2"  // Mandatory
	obj.DesktopAppRunAccountType = logincredential.SharedAccount.String()
	obj.DesktopAppRunAccountName = "shared_account@demo.lab"
	//obj.DesktopAppRunAccountName = "clocal_account"
	obj.DesktopAppCmdline = "-S {database.FQDN}\\{database.InstanceName} -U {account.User} -P {account.Password}"
	appParms := []platform.DesktopAppParam{
		{
			ParamName:        "database",
			ParamType:        cmdparamtype.Database.String(),
			TargetObjectName: "SQL-CENTRIFYSUITE",
		},
		{
			ParamName:          "account",
			ParamType:          cmdparamtype.Account.String(),
			TargetObjectName:   "sa",
			TargetResourceName: "SQL-CENTRIFYSUITE",
			TargetResourceType: "database",
		},
	}
	obj.DesktopAppParams = appParms

	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating desktopapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated desktopapp '%s'\n", obj.Name)

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
		fmt.Printf("Error assign permissions to desktopapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to desktopapp '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"Test Desktop Apps"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding desktopapp to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added desktopapp %s to Sets '%+v'\n", obj.Name, sets)

	////////////////////////////////////////
	// Sample code to update a desktopapp //
	////////////////////////////////////////
	obj = platform.NewDesktopApp(client)
	obj.Name = "Test DesktopApp" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test desktopapp"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating desktopapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated desktopapp '%s'\n", obj.Name)

	////////////////////////////////////////
	// Sample code to delete a desktopapp //
	////////////////////////////////////////
	obj = platform.NewDesktopApp(client)
	obj.Name = "Test DesktopApp" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting desktopapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted desktopapp '%s'\n", obj.Name)

}
