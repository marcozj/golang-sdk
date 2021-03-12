package webcookie

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	log "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
	"golang.org/x/crypto/ssh/terminal"
)

// WebCookie represents a stateful web cookie client
type WebCookie struct {
	restapi.RestClient
	ClientID       string
	ClientSecret   string
	SkipCertVerify bool
	SessionID      string
	TenantID       string
}

func (c *WebCookie) startAuthentication() (*AuthResponse, error) {
	method := "/Security/StartAuthentication"
	args := make(map[string]interface{})
	//args["TenantId"] = c.TenantID
	args["User"] = c.ClientID
	args["Version"] = "1.0"

	body, err := c.postAndGetBody(method, args)
	if err != nil {
		return nil, err
	}

	response, err := NewAuthResponse(body)
	log.Debugf("StartAuthentication response: %+v\n", response)
	if err != nil {
		return nil, err
	}
	if response.Success == false {
		return nil, fmt.Errorf("Failed to initiate authentication: %+v", response.Message)
	}
	c.SessionID = response.Result.SessionID
	c.TenantID = response.Result.TenantID
	if c.SessionID == "" || c.TenantID == "" {
		return nil, fmt.Errorf("SessionId or TenantId is empty")
	}

	return response, nil

}

func (c *WebCookie) advanceAuthentication(authResp *AuthResponse) (string, error) {
	var token string
	challenges := authResp.Result.Challenges

	//var authMechs []AuthMechanism
	for i, challenge := range challenges {
		log.Debugf("Challenge number: %d\n", i+1)
		mechanisms := challenge.Mechanisms
		var authMech AuthMechanism
		if len(mechanisms) > 1 {
			fmt.Print("\n\n")
			// Display mechanisms
			for j, mechanism := range mechanisms {
				displayNum := j + 1
				fmt.Printf("%d. %s\n", displayNum, mechanism.PromptSelectMech)
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Please choose an authentication mechanism: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")
			choice, err := strconv.Atoi(input)
			choice = choice - 1
			if choice < 0 || choice > len(mechanisms)-1 || err != nil {
				return "", fmt.Errorf("Invalid choice")
			}
			authMech = mechanisms[choice]
			log.Debugf("selected mech: %+v\n", authMech)
		} else {
			// Only one mechanism so go ahead to ask for credential
			authMech = mechanisms[0]

		}
		var err error
		// Enter credential
		switch authMech.Name {
		case "UP", "SQ":
			token, err = c.doUPAuthentication(authMech)
			if err != nil {
				return "", fmt.Errorf("Password authentication failed: %+v", err)
			}
		case "OATH", "SMS", "EMAIL", "PF":
			token, err = c.doOOBAuthentication(authMech)
			if err != nil {
				return "", fmt.Errorf("Verificstion code authentication failed: %+v", err)
			}
		}

	}

	return token, nil
}

func (c *WebCookie) postAuthRequest(args map[string]interface{}) (string, error) {
	method := "/Security/AdvanceAuthentication"
	httpresp, err := c.postAndGetResp(method, args)
	if err != nil {
		return "", err
	}
	if httpresp.StatusCode != 200 {
		return "", fmt.Errorf("Bad http status code %v", httpresp.StatusCode)
	}

	body, _ := ioutil.ReadAll(httpresp.Body)
	resp, err := NewAdvanceAuthResponse(body)
	if err != nil {
		return "", fmt.Errorf("Error process respond body: %v", err)
	}
	log.Debugf("AdvanceAuthentication response: %+v\n", resp)

	if !resp.Success {
		return "", fmt.Errorf("Authentication failed: %s", resp.Message)
	}

	authResult := resp.Result["Summary"]
	if authResult != nil {
		switch authResult.(string) {
		case "LoginSuccess":
			// Get auth cookie
			cookie := httpresp.Cookies()
			//log.Debugf("Cookies: %+v\n", getCookieByName(cookie, ".ASPXAUTH"))
			return getCookieByName(cookie, ".ASPXAUTH"), nil
		case "StartNextChallenge":
			return "", nil
		case "OobPending":
			return "", nil
		default:
			return "", fmt.Errorf("%+v", resp.Result)
		}
	}
	return "", nil
}

func (c *WebCookie) doUPAuthentication(authMech AuthMechanism) (string, error) {
	args := make(map[string]interface{})
	args["TenantId"] = c.TenantID
	args["SessionId"] = c.SessionID
	args["MechanismId"] = authMech.MechanismID
	args["Action"] = "Answer"

	if authMech.Question != "" {
		// Security question prompt
		fmt.Printf("%s : ", authMech.Question)
	} else {
		// Password prompt
		fmt.Print("Enter Password: ")
	}
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := strings.TrimSpace(string(bytePassword))
	args["Answer"] = password

	log.Debugf("Performing password authentication with action: %s\n", args["Action"])
	cookie, err := c.postAuthRequest(args)
	if err != nil {
		return "", err
	}

	return cookie, nil
}

func (c *WebCookie) doOOBAuthentication(authMech AuthMechanism) (string, error) {
	// For SMS, Phone, OTP and Email, first trigger the sending of verification code
	args := make(map[string]interface{})
	args["TenantId"] = c.TenantID
	args["SessionId"] = c.SessionID
	args["MechanismId"] = authMech.MechanismID
	args["Action"] = "StartOOB"

	var cookie string
	var err error
	log.Debugf("Starting OOB authentication: %+v\n", args)
	cookie, err = c.postAuthRequest(args)
	if err != nil {
		return "", err
	}

	// After triggering verification code, prompt to enter code
	if cookie == "" {
		fmt.Print("Hit Enter if you have already authenticated out-of-bound or Enter Verification Code: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password := strings.TrimSpace(string(bytePassword))
		if password != "" {
			args["Action"] = "Answer"
			args["Answer"] = password
		} else {
			args["Action"] = "Poll"
		}
		log.Debugf("Performing OOB authentication with action: %s\n", args["Action"])
		cookie, err = c.postAuthRequest(args)
		if err != nil {
			return "", err
		}
	}
	return cookie, nil
}

func (c *WebCookie) postAndGetResp(method string, args map[string]interface{}) (*http.Response, error) {
	service := strings.TrimSuffix(c.Service, "/")
	method = strings.TrimPrefix(method, "/")
	postdata, _ := json.Marshal(args)
	postreq, err := http.NewRequest("POST", service+"/"+method, bytes.NewBuffer(postdata))

	if err != nil {
		return nil, err
	}

	postreq.Header.Add("Content-Type", "application/json")
	postreq.Header.Add("X-CENTRIFY-NATIVE-CLIENT", "Yes")
	postreq.Header.Add("X-CFY-SRC", restapi.SourceHeader)

	for k, v := range c.Headers {
		postreq.Header.Add(k, v)
	}

	httpresp, err := c.Client.Do(postreq)
	if err != nil {
		c.ResponseHeaders = nil
		return nil, err
	}

	// save response heasder
	c.ResponseHeaders = httpresp.Header

	return httpresp, nil
}

func (c *WebCookie) postAndGetBody(method string, args map[string]interface{}) ([]byte, error) {
	httpresp, err := c.postAndGetResp(method, args)
	if err != nil {
		return nil, err
	}
	defer httpresp.Body.Close()

	if httpresp.StatusCode == 200 {
		return ioutil.ReadAll(httpresp.Body)
	}

	body, _ := ioutil.ReadAll(httpresp.Body)
	return nil, fmt.Errorf("POST to %s failed with code %d, body: %s", method, httpresp.StatusCode, body)
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}
