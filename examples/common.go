package examples

import (
	"github.com/marcozj/golang-sdk/enum/authenticationtype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
	"github.com/marcozj/golang-sdk/utils"
)

// GetClient returns authenticated rest client
func GetClient() (*restapi.RestClient, error) {
	logger.SetLevel(logger.LevelDebug)
	logger.SetLogPath("centrifysdk.log")
	logger.EnableErrorStackTrace()

	// Initiate vault client
	vault := utils.VaultClient{}
	vault.AuthType = authenticationtype.OAuth2.String()
	vault.URL = "http://<tenantid>.my.centrify.net"
	//vault.AppID = "CentrifyCLI"
	//vault.Scope = "all"
	//vault.User = ""
	//vault.Password = ""
	vault.Token = ""

	// Authenticate and returns authenticated REST client
	client, err := vault.GetClient()
	if err != nil {
		return nil, err
	}

	return client, err
}
