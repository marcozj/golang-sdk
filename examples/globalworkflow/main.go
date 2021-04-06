package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/directoryservice"
	"github.com/marcozj/golang-sdk/enum/workflowtype"
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

	////////////////////////////////////////////
	// Sample code to update account workflow //
	////////////////////////////////////////////
	wfType := workflowtype.AccountWorkflow.String()
	//wfType := workflowtype.AgentAuthWorkflow.String()
	//wfType := workflowtype.SecretsWorkflow.String()
	//wfType := workflowtype.PrivilegeElevationWorkflow.String()
	obj, err := platform.NewGlobalWorkflow(client, wfType)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	obj.Settings.Enabled = true
	obj.Settings.ApproverList = []platform.WorkflowApprover{
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
			Name:             "System Administrator",
			Type:             "Role",
			DirectoryService: directoryservice.CentrifyDirectory.String(),
			DirectoryName:    "Centrify Directory",
		},
	}

	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating workflowt: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated workflow '%s'\n", obj.Type)

	////////////////////////////////////////////
	// Sample code to update account workflow //
	////////////////////////////////////////////
	obj, err = platform.NewGlobalWorkflow(client, wfType)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = obj.Delete()
	if err != nil {
		fmt.Printf("Error disabling workflowt: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Disabled workflowt '%s'\n", obj.Type)

}
