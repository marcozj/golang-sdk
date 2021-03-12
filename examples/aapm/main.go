package main

import (
	"fmt"
	"os"

	"github.com/marcozj/golang-sdk/enum/keypairtype"
	"github.com/marcozj/golang-sdk/enum/resourcetype"
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

	/////////////////////////////////////////
	// Sample to checkout account password //
	/////////////////////////////////////////
	// Construct account object
	acct1 := platform.NewAccount(client)
	acct1.User = "dbadmin"
	acct1.ResourceName = "MySQL (Demo Lab)"
	acct1.ResourceType = resourcetype.System.String()

	// Checkout password
	pw, err := acct1.CheckoutPassword(false)
	if pw == "" && err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if pw != "" {
		fmt.Printf("Password for account %s in %s is: %s\n", acct1.User, acct1.ResourceName, pw)
		if err != nil {
			fmt.Printf("Password checkin error: %+v\n", err)
		}
	}

	///////////////////////////////////////////////
	// Sample to retrieve SSH key for an account //
	///////////////////////////////////////////////
	// Construct account object
	acct2 := platform.NewAccount(client)
	acct2.User = "testuser"
	acct2.ResourceName = "centos1"
	acct2.ResourceType = resourcetype.System.String()
	// Retrieve SSH key
	thiskey, err := acct2.RetrieveSSHKey(keypairtype.PrivateKey.String(), "")
	if thiskey == "" && err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if thiskey != "" {
		fmt.Printf("SSH key for account %s in %s is: %s\n", acct2.User, acct2.ResourceName, thiskey)
	}

	/////////////////////////////////////////
	// Sample to retrieve SSH key directly //
	/////////////////////////////////////////
	key := platform.NewSSHKey(client)
	key.Name = "testkey"
	key.KeyPairType = keypairtype.PrivateKey.String()
	thatkey, err := key.RetriveSSHKey()
	if thatkey == "" && err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if thatkey != "" {
		fmt.Printf("SSH key for %s is: %s\n", key.Name, thatkey)
	}

	///////////////////////////////
	// Sample to retrieve secret //
	///////////////////////////////
	// Construct secret object
	secret := platform.NewSecret(client)
	secret.SecretName = "testsecret2"
	secret.ParentPath = "folder1\\folder2"
	// Retrieve secret text
	secrettext, err := secret.CheckoutSecret()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Secret text for %s in %s is: %s\n", secret.SecretName, secret.ParentPath, secrettext)

	///////////////////////////////////////
	// Sample to retrieve IAM access key //
	///////////////////////////////////////
	// Construct IAM account object
	acct3 := platform.NewAccount(client)
	acct3.User = "testiam1"
	acct3.ResourceName = "My AWS"
	acct3.ResourceType = resourcetype.CloudProvider.String()
	accesskeyid := "XXXXXXXXXXX"
	secretkey, err := acct3.RetrieveAccessKey(accesskeyid)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Secret key for access key id %s in %s is: %s\n", accesskeyid, acct3.ResourceName, secretkey)
}
