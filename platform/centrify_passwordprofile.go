package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// PasswordProfile - Encapsulates a single Password Profile
type PasswordProfile struct {
	vaultObject
	ProfileFeature string `json:"ProfileFeature,omitempty" schema:"profile_feature,omitempty"`
	ProfileType    string `json:"ProfileType,omitempty" schema:"profile_type,omitempty"` // UserDefined, CheckPointGaia
	// password requirements
	MinimumPasswordLength              int    `json:"MinimumPasswordLength" schema:"minimum_password_length"`
	MaximumPasswordLength              int    `json:"MaximumPasswordLength" schema:"maximum_password_length"`
	AtLeastOneLowercase                bool   `json:"AtLeastOneLowercase" schema:"at_least_one_lowercase"`                                                    // At least one lower-case alpha character
	AtLeastOneUppercase                bool   `json:"AtLeastOneUppercase" schema:"at_least_one_uppercase"`                                                    // At least one upper-case alpha character
	AtLeastOneDigit                    bool   `json:"AtLeastOneDigit" schema:"at_least_one_digit"`                                                            // At least one digit
	ConsecutiveCharRepeatAllowed       bool   `json:"ConsecutiveCharRepeatAllowed,omitempty" schema:"no_consecutive_repeated_char,omitempty"`                 // No consecutive repeated characters
	AtLeastOneSpecial                  bool   `json:"AtLeastOneSpecial" schema:"at_least_one_special_char"`                                                   // At least one special character
	MaximumCharOccurrenceCount         int    `json:"MaximumCharOccurrenceCount,omitempty" schema:"maximum_char_occurrence_count,omitempty"`                  // Restrict number of character occurrences
	SpecialCharSet                     string `json:"SpecialCharSet,omitempty" schema:"special_charset,omitempty"`                                            // Special Characters
	FirstCharacterType                 string `json:"FirstCharacterType,omitempty" schema:"first_character_type,omitempty"`                                   // AlphaOnly or AlphaNumericOnly
	LastCharacterType                  string `json:"LastCharacterType,omitempty" schema:"last_character_type,omitempty"`                                     // AlphaOnly or AlphaNumericOnly
	MinimumAlphabeticCharacterCount    int    `json:"MinimumAlphabeticCharacterCount,omitempty" schema:"minimum_alphabetic_character_count,omitempty"`        // Min number of alpha characters
	MinimumNonAlphabeticCharacterCount int    `json:"MinimumNonAlphabeticCharacterCount,omitempty" schema:"minimum_non_alphabetic_character_count,omitempty"` // Min number of non-alpha characters
}

// NewPasswordProfile is a PasswordProfile constructor
func NewPasswordProfile(c *restapi.RestClient) *PasswordProfile {
	s := PasswordProfile{}
	s.client = c
	s.apiRead = "/ServerManage/GetPasswordProfiles"
	s.apiCreate = "/ServerManage/AddPasswordProfile"
	s.apiDelete = "/ServerManage/DeletePasswordProfile"
	s.apiUpdate = "/ServerManage/UpdatePasswordProfile"

	return &s
}

// Read function fetches an password profile from source, including attribute values. Returns error if any
func (o *PasswordProfile) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	args := make(map[string]interface{})
	args["Caching"] = -1
	queryArg["ProfileTypes"] = "All"
	queryArg["Args"] = args

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
	// Loop through respond results and grab the matched record
	var results = resp.Result["Results"].([]interface{})
	//logger.Debugf("Total returned password profile: %d", len(results))
	// This is the matched list of password profile. There should be only one really
	var pwdpfs []keyValue
	//logger.Debugf("Looking for password profile: %s", o.ID)
	for _, v := range results {
		item := v.(map[string]interface{})
		row := item["Row"].(map[string]interface{})
		//logger.Debugf("Checking name: %s, profiletype: %s", row["Name"], row["ProfileType"])
		if row["ID"] == o.ID {
			logger.Debugf("Found an item: %+v", row)
			// If ProfileType is defined, then compare it
			if o.ProfileType == "" {
				pwdpfs = append(pwdpfs, row)
			} else if row["ProfileType"] == o.ProfileType {
				pwdpfs = append(pwdpfs, row)
			}
		}
	}
	err = queryError(len(pwdpfs))
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	mapToStruct(o, pwdpfs[0])

	return nil
}

// Delete function deletes an password profile and returns a map that contains deletion result
func (o *PasswordProfile) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Create function creates an password profile and returns a map that contains update result
func (o *PasswordProfile) Create() (*restapi.StringResponse, error) {
	var queryArg = make(map[string]interface{})

	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	logger.Debugf("Generated Map for Create(): %+v", queryArg)

	resp, err := o.client.CallStringAPI(o.apiCreate, queryArg)
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
	o.ID = resp.Result

	return resp, nil
}

// Update function updates an existing password profile and returns a map that contains update result
func (o *PasswordProfile) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg, err := generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

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

// Query function returns a single password profile object
func (o *PasswordProfile) Query() (map[string]interface{}, error) {
	var queryArg = make(map[string]interface{})
	args := make(map[string]interface{})
	args["Caching"] = -1
	queryArg["ProfileTypes"] = "All"
	queryArg["Args"] = args

	// Attempt to read from an upstream API
	resp, err := o.client.CallGenericMapAPI(o.apiRead, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
	// Loop through respond results and grab the matched record
	var results = resp.Result["Results"].([]interface{})

	// This is the matched list of password profile. There should be only one really
	var pwdpfs []keyValue
	for _, v := range results {
		item := v.(map[string]interface{})
		row := item["Row"].(map[string]interface{})
		if row["Name"] == o.Name {
			logger.Debugf("Found an item: %+v", row)
			// If ProfileType is defined, then compare it
			if o.ProfileType == "" {
				pwdpfs = append(pwdpfs, row)
			} else if row["ProfileType"] == o.ProfileType {
				pwdpfs = append(pwdpfs, row)
			}
		}
	}

	err = queryError(len(pwdpfs))
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return pwdpfs[0], nil
}

// GetIDByName returns password profile ID by name
func (o *PasswordProfile) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("Password profile name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving password profile: %s", err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves password profile from tenant by name
func (o *PasswordProfile) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of password profile %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// DeleteByName deletes a password profile by name
func (o *PasswordProfile) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of password profile %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	return resp, nil
}

/*
	API to manage password profile

	Get password profile

		Request body format
		{
			"RRFormat": true,
			"ProfileTypes": "All",
			"Args": {
				"PageNumber": 1,
				"PageSize": 100000,
				"Limit": 100000,
				"SortBy": "",
				"direction": "False",
				"Caching": -1
			}
		}

		Respond result
		{
			"success": true,
			"Result": {
				"IsAggregate": false,
				"Count": 22,
				"Columns": [
					...}\
				],
				"FullCount": 22,
				"Results": [
					{
						"Entities": [
							{
								"Type": "PasswordProfile",
								"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
								"IsForeignKey": false
							}
						],
						"Row": {
							"LastCharacterType": "AnyChar",
							"_entitycontext": "*",
							"FirstCharacterType": "AlphaOnly",
							"SpecialCharSet": "!@#$%&()+,-./:;<=>?[\\]^_{|}~",
							"Name": "Check Point Gaia Profile",
							"ConsecutiveCharRepeatAllowed": false,
							"_metadata": {
								"Version": 1,
								"IndexingVersion": 1
							},
							"_PartitionKey": "centrify",
							"_Timestamp": "/Date(1596273224728)/",
							"AtLeastOneSpecial": true,
							"_encryptkeyid": "RootIndex:1",
							"AtLeastOneLowercase": true,
							"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"Description": "Default profile for Check Point Gaia systems",
							"MaximumPasswordLength": 32,
							"ProfileFeature": "Infrastructure",
							"AtLeastOneDigit": true,
							"_TableName": "pvpasswordprofile",
							"MinimumPasswordLength": 12,
							"AtLeastOneUppercase": true,
							"ProfileType": "CheckPointGaia"
						}
					},
				],
				"ReturnID": ""
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Create password profile

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Test PW Profile 1",
			"Description": "Test PW Profile 1",
			"MinimumPasswordLength": 12,
			"MaximumPasswordLength": 32,
			"AtLeastOneLowercase": true,
			"AtLeastOneUppercase": true,
			"AtLeastOneDigit": true,
			"ConsecutiveCharRepeatAllowed": true,
			"AtLeastOneSpecial": true,
			"jsutil-checkbox-8128-inputEl": true,
			"MaximumCharOccurrenceCount": 2,
			"SpecialCharSet": "!#$%&()*+,-./:;<=>?@[\\]^_{|}~",
			"jsutil-checkbox-8138-inputEl": true,
			"FirstCharacterType": "AlphaOnly",
			"jsutil-checkbox-8142-inputEl": true,
			"LastCharacterType": "AlphaNumericOnly",
			"MinimumAlphabeticCharacterCount": 1,
			"MinimumNonAlphabeticCharacterCount": 1
		}

		Respond result
		{
			"success": true,
			"Result": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Update password profile

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Name": "Test PW Profile",
			"Description": "testa afasdf jhjkhjkh",
			"MinimumPasswordLength": 12,
			"MaximumPasswordLength": 32,
			"AtLeastOneLowercase": true,
			"AtLeastOneUppercase": true,
			"AtLeastOneDigit": true,
			"ConsecutiveCharRepeatAllowed": true,
			"AtLeastOneSpecial": true,
			"jsutil-checkbox-4100-inputEl": false,
			"SpecialCharSet": "!#$%&()*+,-./:;<=>?@[\\]^_{|}~",
			"jsutil-checkbox-4110-inputEl": false,
			"jsutil-checkbox-4114-inputEl": false,
			"LastCharacterType": "AnyChar",
			"FirstCharacterType": "AnyChar"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete password profile

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": true,
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

*/
