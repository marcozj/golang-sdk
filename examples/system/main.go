package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/computerclass"
	"github.com/marcozj/golang-sdk/enum/directoryservice"
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
	// Sample code to create a system //
	////////////////////////////////////
	obj := platform.NewSystem(client)
	obj.Name = "Test System"                           // Mandatory
	obj.ComputerClass = computerclass.Windows.String() // Mandatory
	obj.FQDN = "testsystem.example.test"               // Mandatory
	obj.Description = "This is a test system"
	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating system: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created system '%s'\n", obj.Name)

	// Assign default logn profile. It can only be done after system has been created
	authProfile := platform.NewAuthenticationProfile(client)
	authProfile.Name = "Default Other Login Profile"
	authID, err := authProfile.GetIDByName()
	if err != nil {
		fmt.Printf("Error retrieving authenticaiton profile %s", authProfile.Name)
		os.Exit(1)
	}
	obj.LoginDefaultProfile = authID
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating system: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated system '%s'\n", obj.Name)

	//------------------//
	// Assign permissions
	myPermissions := []platform.Permission{
		{
			PrincipalName: "System Administrator",
			PrincipalType: "Role",
			RightList:     []string{platform.Right.Grant, platform.Right.View, platform.Right.Edit, platform.Right.Delete},
		},
		{
			PrincipalName: "admin@centrify.com.207",
			PrincipalType: "User",
			RightList:     []string{platform.Right.View, platform.Right.Edit, platform.Right.ManageSession, platform.Right.RequestZoneRole},
		},
	}
	obj.ResolveValidPermissions() // This is needed for system because Windows/Unix type of system has different set of permissions than others
	err = platform.ResolvePermissions(client, myPermissions, obj.ValidPermissions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Assign resolved permission
	obj.Permissions = myPermissions
	_, err = obj.SetPermissions(false)
	if err != nil {
		fmt.Printf("Error assign permissions to system: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Assigned permissions %+v to system '%s'\n", obj.Permissions, obj.Name)

	// Assign to Sets
	sets := []string{"Custom Systems", "LAB Systems"}
	err = obj.AddToSetsByName(sets)
	if err != nil {
		fmt.Printf("Error adding system to Sets %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added system %s to Sets '%+v'\n", obj.Name, sets)

	////////////////////////////////////
	// Sample code to update a system //
	////////////////////////////////////
	obj = platform.NewSystem(client)
	obj.Name = "Test System"             // Mandatory
	obj.ComputerClass = "Windows"        // Mandatory
	obj.FQDN = "testsystem.example.test" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test system - updated"

	// Update Agent Auth and Elevation workflows. Only for Windows/Unix that has Centrify Client deployed
	obj.AgentAuthWorkflowEnabled = true
	obj.AgentAuthWorkflowApprovers = []platform.WorkflowApprover{
		{
			Name:             "System Administrator",
			Type:             "Role",
			DirectoryService: directoryservice.CentrifyDirectory.String(),
			DirectoryName:    "Centrify Directory",
		},
	}
	obj.PrivilegeElevationWorkflowEnabled = true
	obj.PrivilegeElevationWorkflowApprovers = []platform.WorkflowApprover{
		{
			Name:             "System Administrator",
			Type:             "Role",
			DirectoryService: directoryservice.CentrifyDirectory.String(),
			DirectoryName:    "Centrify Directory",
		},
	}

	// If this is a Zone joined machine, retrieve domain id and assign DomainID attribute
	domain := platform.NewDomain(client)
	domain.Name = "demo.lab"
	domainid, err := domain.GetIDByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	obj.DomainID = domainid

	// Enable zone workflow
	obj.ZoneRoleWorkflowEnabled = true
	obj.DomainOperationsEnabled = true
	// Assign Zone Roles
	obj.UseDomainWorkflowRoles = false // This must be false in order to override zone roles
	obj.ZoneRoleWorkflowRoleList = []platform.ZoneRole{
		{Name: "cfyw-Windows Super Admin/Windows Zone"}, // zone role name is in format of "<zone role name>/<zone name>"
		{Name: "cfyw-Windows System Admin/Windows Zone"},
	}
	// Assign Zone Roles approver list
	obj.UseDomainWorkflowApprovers = false // This must be false in order to override approver list
	obj.ZoneRoleWorkflowApproverList = []platform.WorkflowApprover{
		{
			Name:             "System Administrator",
			Type:             "Role",
			DirectoryService: directoryservice.CentrifyDirectory.String(),
			DirectoryName:    "Centrify Directory",
		},
	}

	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating system: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated system '%s'\n", obj.Name)

	////////////////////////////////////
	// Sample code to delete a system //
	////////////////////////////////////
	obj = platform.NewSystem(client)
	obj.Name = "Test System"             // Mandatory
	obj.ComputerClass = "Windows"        // Mandatory
	obj.FQDN = "testsystem.example.test" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting system: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted system '%s'\n", obj.Name)

}
