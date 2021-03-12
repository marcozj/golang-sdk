package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/authenticationtype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/platform"
	"github.com/marcozj/golang-sdk/utils"
)

func main() {
	logger.SetLevel(logger.LevelDebug)
	logger.SetLogPath("centrifysdk.log")

	// Initiate client using OAuth client credential
	vault := utils.VaultClient{}
	vault.AuthType = authenticationtype.OAuth2.String()
	vault.URL = "http://tenantid.my.centrify.net"
	vault.AppID = ""
	vault.Scope = ""
	vault.User = ""
	vault.Password = ""

	/*
		// Initiate client using OAuth token
		vault := utils.VaultClient{}
		vault.AuthType = authenticationtype.OAuth2.String()
		vault.URL = "http://tenantid.my.centrify.net"
		vault.Scope = ""
		vault.Token = ""

		// Initiate client using DMC
		vault := utils.VaultClient{}
		vault.AuthType = authenticationtype.DelegatedMachineCredential.String()
		vault.URL = "http://tenantid.my.centrify.net"
		vault.Scope = ""
	*/

	// Altenate to provide attributes to establish vault client, use command line parameter to provide them
	//vault.GetCmdParms()

	// Authenticate and returns authenticated REST client
	client, err := vault.GetClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Sample usage of using authenticated REST client to query users
	sql := "Select ID, Username from User"
	var args map[string]interface{}
	args["Caching"] = -1
	//args["PageSize"] = 10000
	//args["Limit"] = 10000
	results, err := platform.RedRockQuery(client, sql, args)
	if err != nil {
		fmt.Printf("\nFailed to query: %v\n", err)
	} else {
		for _, v := range results {
			var resultItem = v.(map[string]interface{})
			var row = resultItem["Row"].(map[string]interface{})
			fmt.Printf("%+v\n", row)
		}
	}
}
