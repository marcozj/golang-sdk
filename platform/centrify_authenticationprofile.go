package platform

import (
	"encoding/json"
	"fmt"
	"strings"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// AuthenticationProfile - Encapsulates a single Authentication Profile
type AuthenticationProfile struct {
	vaultObject
	UUID              string          `json:"Uuid,omitempty" schema:"uuid,omitempty"`
	DurationInMinutes int             `json:"DurationInMinutes" schema:"pass_through_duration"` // Challenge Pass-Through Duration. Can't omitempty because 0 mean no pass-through
	Challenges        []string        `json:"Challenges,omitempty" schema:"challenges,omitempty"`
	AdditionalData    *AdditionalData `json:"AdditionalData,omitempty" schema:"additional_data,omitempty"`
	NumberOfQuestions int             `json:"-"`
	Challenge1        []string        `json:"-"`
	Challenge2        []string        `json:"-"`
}

// AdditionalData for AuthenticationProfile
type AdditionalData struct {
	NumberOfQuestions int `json:"NumberOfQuestions" schema:"number_of_questions"` // Number of questions user must answer
}

// NewAuthenticationProfile is a AuthenticationProfile constructor
func NewAuthenticationProfile(c *restapi.RestClient) *AuthenticationProfile {
	s := AuthenticationProfile{}
	s.client = c
	s.apiRead = "/AuthProfile/GetProfile"
	s.apiCreate = "/AuthProfile/SaveProfile"
	s.apiDelete = "/AuthProfile/DeleteProfile"
	s.apiUpdate = "/AuthProfile/SaveProfile"

	return &s
}

type sliceAPIResponse struct {
	Success bool `json:"success"`
	Result  []interface{}
	Message string
}

// Read function fetches an authentication profile from source, including attribute values. Returns error if any
func (o *AuthenticationProfile) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["Uuid"] = o.ID

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}

	mapToStruct(o, resp.Result)
	o.expandChallenges()
	o.expandNumberOfQuestions()

	return nil
}

// Delete function deletes an authentication profile and returns a map that contains deletion result
func (o *AuthenticationProfile) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("uuid")
}

// Create function creates an authentication profile and returns a map that contains update result
func (o *AuthenticationProfile) Create() (*restapi.GenericMapResponse, error) {
	var queryArg = make(map[string]interface{})

	// Flatten chanllenges data first
	err := o.flattenChallenges()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	// Flatten NumberOfQuestions
	o.falttenNumberOfQuestions()
	settings, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	queryArg["settings"] = settings

	logger.Debugf("Generated Map for Create(): %+v", queryArg)

	resp, err := o.client.CallGenericMapAPI(o.apiCreate, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	// Assign ID after successful creation so that the same object can be used for subsequent update operation
	o.ID = resp.Result["Uuid"].(string)

	return resp, nil
}

// Update function updates an existing authentication profile and returns a map that contains update result
func (o *AuthenticationProfile) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})

	// Flatten chanllenges data first
	err := o.flattenChallenges()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	// Flatten NumberOfQuestions
	o.falttenNumberOfQuestions()
	settings, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	queryArg["settings"] = settings

	logger.Debugf("Generated Map for Update(): %+v", queryArg)

	resp, err := o.client.CallGenericMapAPI(o.apiUpdate, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	return resp, nil
}

// Query function returns a single authentication profile object
func (o *AuthenticationProfile) Query() (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	args := make(map[string]interface{})
	args["Caching"] = -1
	queryArg["Args"] = args

	// Attempt to read from an upstream API
	resp, err := o.client.CallRawAPI("/AuthProfile/GetProfileList", queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	reply := &sliceAPIResponse{}
	err = json.Unmarshal(resp, &reply)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, fmt.Errorf("Failed to unmarshal sliceAPIResponse from HTTP response: %v", err)
	}
	if !reply.Success {
		logger.Errorf(reply.Message)
		return nil, fmt.Errorf(reply.Message)
	}

	// This is the matched list of authentication profile. There should be only one
	var autheProfs []keyValue
	for _, v := range reply.Result {
		item := v.(map[string]interface{})
		if item["Name"] == o.Name {
			autheProfs = append(autheProfs, item)
		}
	}

	err = queryError(len(autheProfs))
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return autheProfs[0], nil
}

// GetIDByName returns authentication profile ID by name
func (o *AuthenticationProfile) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Authentication profile name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		return "", fmt.Errorf("Error retrieving authentication profile: %s", err)
	}
	o.ID = result["Uuid"].(string)

	return o.ID, nil
}

// GetByName retrieves authentication profile from tenant by name
func (o *AuthenticationProfile) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return fmt.Errorf("Failed to find ID of authentication profile %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a authentication profile by name
func (o *AuthenticationProfile) DeleteByName() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return nil, fmt.Errorf("Failed to find ID of authentication profile %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

// flattenChallenges converts Challenge1 and Challenge2 to Challenges for creation and update
func (o *AuthenticationProfile) flattenChallenges() error {
	if o.Challenges == nil || len(o.Challenges) == 0 {
		// If this is called from Create method, Challenges should be empty
		if o.Challenge1 != nil && len(o.Challenge1) > 0 {
			o.Challenges = append(o.Challenges, FlattenSliceToString(o.Challenge1))
		} else {
			return fmt.Errorf("Missing first challenges")
		}
		if o.Challenge2 != nil && len(o.Challenge2) > 0 {
			o.Challenges = append(o.Challenges, FlattenSliceToString(o.Challenge2))
		}
	} else {
		// If this is called from Update method, there is already values in Challenges
		var oldchallenge1, oldchallenge2 string
		for i, v := range o.Challenges {
			if i == 0 {
				oldchallenge1 = v
			}
			if i == 1 {
				oldchallenge2 = v
			}
		}
		// empty Challenges
		o.Challenges = nil
		// Add updated changes
		if o.Challenge1 == nil {
			o.Challenges = append(o.Challenges, oldchallenge1)
		} else {
			o.Challenges = append(o.Challenges, FlattenSliceToString(o.Challenge1))
		}
		if o.Challenge2 == nil {
			o.Challenges = append(o.Challenges, oldchallenge2)
		} else {
			o.Challenges = append(o.Challenges, FlattenSliceToString(o.Challenge2))
		}
	}

	return nil
}

// expandChallenges fills Challenge1 & Challenge2 attributes from Challenges attribute
func (o *AuthenticationProfile) expandChallenges() {
	for i, v := range o.Challenges {
		if i == 0 {
			o.Challenge1 = strings.Split(v, ",")
		}
		if i == 1 {
			o.Challenge2 = strings.Split(v, ",")
		}
	}
}

// falttenNumberOfQuestions converts NumberOfQuestions to AdditionalData.NumberOfQuestions
func (o *AuthenticationProfile) falttenNumberOfQuestions() {
	if o.NumberOfQuestions > 0 {
		a := AdditionalData{}
		a.NumberOfQuestions = o.NumberOfQuestions
		o.AdditionalData = &a
	}
}

func (o *AuthenticationProfile) expandNumberOfQuestions() {
	if o.AdditionalData.NumberOfQuestions > 0 {
		o.NumberOfQuestions = o.AdditionalData.NumberOfQuestions
	}
}

/*
	API to manage password profile

	Get authentication profile

		Request body format
		{
			"uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"Uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Name": "LAB Step-up Authentication Profile",
				"DurationInMinutes": 0,
				"Challenges": [
					"OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
				],
				"AdditionalData": {
					"NumberOfQuestions": 1
				}
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}


	Create authentication profile

		Request body format
		{
			"settings": {
				"Name": "test",
				"Challenges": [
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ",
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
				],
				"DurationInMinutes": 30,
				"AdditionalData": {
					"NumberOfQuestions": 1
				}
			}
		}

		Respond result
		{
			"success": true,
			"Result": {
				"Uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Name": "test",
				"DurationInMinutes": 30,
				"Challenges": [
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ",
					"UP,OTP,PF,SMS,EMAIL,OATH,RADIUS,U2F,SQ"
				],
				"AdditionalData": {
					"NumberOfQuestions": 1
				}
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update authentication profile

		Request body format


		Respond result

	Delete authentication profile

		Request body format
		{
			"uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": null,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

*/
