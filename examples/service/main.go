package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/servicetype"
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

	///////////////////////////////////
	// Sample code to create service //
	///////////////////////////////////
	obj := platform.NewService(client)
	obj.Name = "TestWindowsService" // Mandatory
	obj.SystemName = "member1"      // Mandatory
	obj.ServiceType = servicetype.WindowsService.String()
	obj.EnableManagement = true
	obj.AdminAccountUPN = "ad_admin@demo.lab" // Mandatory if EnableManagement is true
	obj.MultiplexedAccountName = "Test"       // Mandatory if EnableManagement is true
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created service '%s'\n", obj.Name)

	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.Edit, platform.Right.Delete},
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
		fmt.Printf("Error assign permissions to service: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to service '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"Test Set"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding service to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added service %s to Sets '%+v'\n", obj.Name, sets)

	///////////////////////////////////
	// Sample code to update service //
	///////////////////////////////////
	obj = platform.NewService(client)
	obj.Name = "TestWindowsService" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test service"
	obj.RestartService = true
	obj.RestartTimeRestriction = true
	obj.DaysOfWeek = "Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday"
	obj.RestartStartTime = "09:00"
	obj.RestartEndTime = "10:00"
	obj.UseUTCTime = false
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating service: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated service '%s'\n", obj.Name)

	///////////////////////////////////
	// Sample code to delete service //
	///////////////////////////////////
	obj = platform.NewService(client)
	obj.Name = "TestWindowsService" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting service: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted service '%s'\n", obj.Name)

}
