package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/marcozj/golang-sdk/dmc"
)

func main() {
	scope := flag.String("scope", "", "Scope definition to be used for the machine credential")
	url := flag.String("url", "", "URL of the tenant")
	query := flag.String("query", "", "Query string")
	flag.Parse()
	if *url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *scope == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	rpc := dmc.NewLRPC2()
	//rpc := dmc.NewWinLRPC2()
	token, err := rpc.GetToken(*scope)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(token)

	call := dmc.DMC{}
	call.Service = *url
	call.Scope = *scope

	client, err := call.GetClient()
	if err != nil {
		fmt.Println(err)
	}

	var queryArg = make(map[string]interface{})
	queryArg["Script"] = *query
	//queryArg["Args"] = subArgs

	resp, err := client.CallGenericMapAPI("/RedRock/query", queryArg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Result)

}
