package utils

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/marcozj/golang-sdk/enum/authenticationtype"
	"golang.org/x/crypto/ssh/terminal"
)

// GetCmdParms parse command line argument
func (c *VaultClient) GetCmdParms() {
	// Common arguments
	authTypePtr := flag.String("auth", "oauth", "Authentication type <oauth|unpw|dmc>")
	urlPtr := flag.String("url", "", "Centrify tenant URL (Required)")
	skipCertPtr := flag.Bool("skipcert", false, "Ignore certification verification")
	debugPtr := flag.Bool("debug", false, "Trun on debug logging")

	// Other arguments
	appIDPtr := flag.String("appid", "", "OAuth2 application ID. Required if auth = oauth")
	scopePtr := flag.String("scope", "", "OAuth2 or DMC scope definition. Required if auth = oauth or dmc")
	tokenPtr := flag.String("token", "", "OAuth2 or DMC token. Optional if auth = oauth or dmc")
	usernamePtr := flag.String("user", "", "Authorized user to login to tenant. Required if auth = unpw. Optional if auth = oauth")
	passwordPtr := flag.String("password", "", "User password. You will be prompted to enter password if this isn't provided")

	prgname := os.Args[0]
	flag.Usage = func() {
		fmt.Printf("Usage: %s -auth oauth -url https://<tenant>.my.centrify.net -user user@company.com \n", prgname)
		flag.PrintDefaults()
	}

	flag.Parse()

	// Verify command argument length
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Verify authTypePtr value
	authChoices := map[string]bool{"oauth": true, "unpw": true, "dmc": true}
	if _, validChoice := authChoices[*authTypePtr]; !validChoice {
		flag.Usage()
		os.Exit(1)
	}
	// Check required argument that do not have default value
	if *urlPtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch strings.ToLower(*authTypePtr) {
	case authenticationtype.OAuth2.String():
		if *appIDPtr == "" || *scopePtr == "" {
			flag.Usage()
			os.Exit(1)
		}
		// Either token or username must be provided
		if *tokenPtr == "" && *usernamePtr == "" {
			flag.Usage()
			os.Exit(1)
		}
		// If password isn't provided, prompt for it
		if *passwordPtr == "" {
			fmt.Print("Enter Password: ")
			bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
			password := strings.TrimSpace(string(bytePassword))
			*passwordPtr = password
			fmt.Println()
		}
	case authenticationtype.UsernamePassword.String():
		if *urlPtr == "" || *usernamePtr == "" {
			flag.Usage()
			os.Exit(1)
		}
		// If password isn't provided, prompt for it
		if *passwordPtr == "" {
			fmt.Print("Enter Password: ")
			bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
			password := strings.TrimSpace(string(bytePassword))
			*passwordPtr = password
			fmt.Println()
		}
	case authenticationtype.DelegatedMachineCredential.String():
		if *tokenPtr == "" && *scopePtr == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	// Assign argument values to struct
	c.AuthType = *authTypePtr
	c.URL = *urlPtr
	c.AppID = *appIDPtr
	c.Scope = *scopePtr
	c.Token = *tokenPtr
	c.User = *usernamePtr
	c.Password = *passwordPtr
	c.Skipcert = *skipCertPtr
	c.Debug = *debugPtr
}
