package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/settype"
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
	// Sample code to create a manual set //
	////////////////////////////////////////
	obj, err := platform.NewManualSetWithType(client, settype.System.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//obj.SubObjectType = platform.SetSubtype.Desktop
	obj.Name = "Test Set" // Mandatory
	//obj.SubObjectType = "Desktop"  // If type is "Application", subtype to be set to either "Desktop" or "Web"
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating Set: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Set '%s'\n", obj.Name)

	//------------------//
	// Set permissions
	setPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Edit, platform.Right.Delete},
		},
		{
			PrincipalName: "admin@centrify.com.207",
			PrincipalType: "User",
			RightList:     []string{platform.Right.View, platform.Right.Edit},
		},
	}
	err = platform.ResolvePermissions(client, setPermissions, obj.ValidPermissions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Assign resolved permission
	obj.Permissions = setPermissions
	_, err = obj.SetPermissions(false)
	if err != nil {
		fmt.Printf("Error assign permissions to Set: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to Set '%s'\n", obj.Permissions, obj.Name)

	//------------------//
	// Set member permissions
	memberPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.ManageSession, platform.Right.AgentAuth, platform.Right.Delete},
		},
		{
			PrincipalName: "admin@centrify.com.207",
			PrincipalType: "User",
			RightList:     []string{platform.Right.View, platform.Right.ManageSession, platform.Right.Edit},
		},
	}
	err = platform.ResolvePermissions(client, memberPermissions, obj.ValidMemberPermissions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	obj.MemberPermissions = memberPermissions
	_, err = obj.SetMemberPermissions(false)
	if err != nil {
		fmt.Printf("Error assign member permissions to Set: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned member permissions %+v to Set '%s'\n", obj.MemberPermissions, obj.Name)

	/////////////////////////////////
	// Sample code to update a Set //
	/////////////////////////////////
	obj = platform.NewManualSet(client)
	obj.Name = "Test Set"     // Mandatory
	obj.ObjectType = "Server" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test set"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating Set: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated Set '%s'\n", obj.Name)

	/////////////////////////////////
	// Sample code to delete a Set //
	/////////////////////////////////
	obj = platform.NewManualSet(client)
	obj.Name = "Test Set"     // Mandatory
	obj.ObjectType = "Server" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting Set: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted Set '%s'\n", obj.Name)

}
