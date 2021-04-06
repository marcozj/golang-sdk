package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/webapp/oauth/applicationtemplate"
	"github.com/marcozj/golang-sdk/enum/webapp/oauth/clientidtype"
	"github.com/marcozj/golang-sdk/enum/webapp/oauth/tokentype"
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
	obj := platform.NewOauthWebApp(client)
	obj.Name = "Test OAuth Server WebApp"                        // Mandatory
	obj.TemplateName = applicationtemplate.OAuth2Server.String() // Mandatory

	_, err = obj.Create()
	if err != nil {
		fmt.Printf("Error creating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s '%s'\n", platform.GetVarType(obj), obj.Name)

	// Many of parameters can only be set in update process rather than initinal creation
	obj.ApplicationID = "TestOAuthClient" // No space is allowed
	obj.OAuthProfile = &platform.OAuthProfile{
		ClientIDType: int(clientidtype.Confidential),
		//Issuer: "",
		Audience: "MyAudience",
		//AllowedClients:       []string{"client1", "client2"},
		MustBeOauthClient:    true,
		Redirects:            []string{"https://example.com", "https://test.com"},
		TokenType:            tokentype.JwtRS256.String(),
		AllowedAuth:          "AuthorizationCode,Implicit",
		TokenLifetime:        "6:00:00",
		AllowRefresh:         true,
		RefreshLifetime:      "200.00:00:00",
		ConfirmAuthorization: true,
		AllowScopeSelect:     true,
		KnownScopes: []platform.OAuthScope{
			{
				Name:        "cli",
				Description: "Used for CLI call",
			},
			{
				Name:        "aapm",
				Description: "Used for AAPM calls",
			},
		},
	}
	obj.Script = "// custom setClaim\nsetClaim('Department', LoginUser.Get('department'));"

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
	obj = platform.NewOauthWebApp(client)
	obj.Name = "Test OAuth Server WebApp" // Mandatory
	err = obj.GetByName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set the atributes that are to be updated
	obj.Description = "This is a test OAuth server webapp"
	_, err = obj.Update()
	if err != nil {
		fmt.Printf("Error updating webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated webapp '%s'\n", obj.Name)

	////////////////////////////////////
	// Sample code to delete a webapp //
	////////////////////////////////////
	obj = platform.NewOauthWebApp(client)
	obj.Name = "Test OAuth Server WebApp" // Mandatory
	_, err = obj.DeleteByName()
	if err != nil {
		fmt.Printf("Error deleting webapp: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Deleted webapp '%s'\n", obj.Name)

}
