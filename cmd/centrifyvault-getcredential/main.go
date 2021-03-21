package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/marcozj/golang-sdk/enum/resourcetype"
	"github.com/marcozj/golang-sdk/platform"
	"github.com/marcozj/golang-sdk/utils"
)

// CliParameters is data structure for commandline parameters
type CliParameters struct {
	CredentialPath string
	SaveToHome     bool
}

type vaultObject struct {
	resourceType string
	resourceName string
	parentPath   string
	secretName   string
	accesskeyID  string
}

func main() {
	//logger.SetLevel(logger.LevelDebug)
	//logfile := os.Args[0] + ".log"
	//logger.SetLogPath(logfile)

	pars := &CliParameters{}
	vault := &utils.VaultClient{}
	getCmdParms(vault, pars)

	// Construct vault object from credential path
	vo, err := getVaultObject(pars.CredentialPath)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// Authenticate and returns authenticated REST client
	client, err := vault.GetClient()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	switch strings.ToLower(vo.resourceType) {
	case resourcetype.System.String(), resourcetype.Database.String(), resourcetype.Domain.String():
		acct := platform.NewAccount(client)
		acct.User = vo.secretName
		acct.ResourceName = vo.resourceName
		acct.ResourceType = vo.resourceType
		// Checkout password
		pw, err := acct.CheckoutPassword(false)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Print(pw)
	case resourcetype.CloudProvider.String():
		acct := platform.NewAccount(client)
		acct.User = vo.secretName
		acct.ResourceName = vo.resourceName
		acct.ResourceType = resourcetype.CloudProvider.String()
		secretkey, err := acct.RetrieveAccessKey(vo.accesskeyID)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Print(secretkey)
	case "secret":
		secret := platform.NewSecret(client)
		secret.SecretName = vo.secretName
		secret.ParentPath = vo.parentPath
		// Retrieve secret text
		secrettext, err := secret.CheckoutSecretAndFile(pars.SaveToHome)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Print(secrettext)
	}
}

func getVaultObject(credPath string) (*vaultObject, error) {
	var vo vaultObject
	credparts := strings.Split(credPath, "/")
	splitLength := len(credparts)
	vo.resourceType = credparts[0]
	switch vo.resourceType {
	case resourcetype.System.String(), resourcetype.Database.String(), resourcetype.Domain.String():
		// Handle vaulted account for system, database and domain
		// Minimumlly must be at least "system/systemname/accountname"
		if splitLength > 2 {
			vo.resourceName = credparts[1]
			vo.secretName = credparts[2]
		}
		if vo.resourceName == "" || vo.secretName == "" {
			return nil, fmt.Errorf("invalid credential path %s", credPath)
		}
	case resourcetype.CloudProvider.String():
		// Handle AWS IAM account access key
		// Credential path format should be "cloudprovider/My AWS/iamaccount/accesskeyid"
		if splitLength > 3 {
			vo.resourceName = credparts[1]
			vo.secretName = credparts[2]
			vo.accesskeyID = credparts[3]
		}
		if vo.resourceName == "" || vo.secretName == "" || vo.accesskeyID == "" {
			return nil, fmt.Errorf("invalid credential path %s", credPath)
		}
	case "secret":
		// Handle secret
		if splitLength > 1 {
			// Minimumlly must be at least "secret/secretname"
			// Or "secret/folderpath/secretname"
			// or "secret/folder1\folder2/secretname"
			vo.secretName = credparts[splitLength-1]
			// Extract only the path from original split
			if splitLength > 2 {
				for i := 1; i <= splitLength-2; i++ {
					if vo.parentPath != "" {
						// if it is not the first level of folder, add "\". Double "\\" is to escape "\"
						// In Golang, it takes single "\" Script:SELECT * FROM DataVault WHERE 1=1 AND SecretName='testsecret2' AND ParentPath='folder1\folder2'
						// In Postman, it takes double "\\" Script:SELECT * FROM DataVault WHERE 1=1 AND SecretName='testsecret2' AND ParentPath='folder1\\folder2'
						vo.parentPath = vo.parentPath + "\\"
					}
					vo.parentPath = vo.parentPath + credparts[i]
				}
			}
			if vo.secretName == "" {
				return nil, fmt.Errorf("invalid credential path %s", credPath)
			}
		}
	default:
		return nil, fmt.Errorf("invalid resource type")
	}

	return &vo, nil
}
