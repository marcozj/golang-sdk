package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// User - Encapsulates a single user
type User struct {
	vaultObject
	apiUpdatePassword string

	Name                    string `json:"Name,omitempty" schema:"username,omitempty"`
	Mail                    string `json:"Mail,omitempty" schema:"email,omitempty"` // Email address
	DisplayName             string `json:"DisplayName,omitempty" schema:"displayname,omitempty"`
	Password                string `json:"Password,omitempty" schema:"password,omitempty"`
	ConfirmPassword         string `json:"confirmPassword,omitempty" schema:"confirm_password,omitempty"`
	PasswordNeverExpire     bool   `json:"PasswordNeverExpire,omitempty" schema:"password_never_expire,omitempty"`          // Password never expires
	ForcePasswordChangeNext bool   `json:"ForcePasswordChangeNext,omitempty" schema:"force_password_change_next,omitempty"` // Require password change at next login
	OauthClient             bool   `json:"OauthClient" schema:"oauth_client"`                                               // Is OAuth confidential client
	SendEmailInvite         bool   `json:"SendEmailInvite,omitempty" schema:"send_email_invite,omitempty"`                  // Send email invite for user profile setup
	OfficeNumber            string `json:"OfficeNumber,omitempty" schema:"office_number,omitempty"`
	HomeNumber              string `json:"HomeNumber,omitempty" schema:"home_number,omitempty"`
	MobileNumber            string `json:"MobileNumber,omitempty" schema:"mobile_number,omitempty"`
	//RedirectMFA             bool   `json:"jsutil-checkbox-2598-inputEl" schema:"redirect_mfa"` // Redirect multi factor authentication to a different user account
	RedirectMFAUserID string `json:"CmaRedirectedUserUuid,omitempty" schema:"redirect_mfa_user_id,omitempty"` // Redirect multi factor authentication to a different user account
	ReportsTo         string `json:"ReportsTo" schema:"manager_username"`                                     // Manager

	// Roles
	Roles []string `json:"Roles,omitempty" schema:"roles,omitempty"`
}

// NewUser is a user constructor
func NewUser(c *restapi.RestClient) *User {
	s := User{}
	s.client = c
	s.apiRead = "/UserMgmt/GetUserAttributes"
	s.apiCreate = "/CDirectoryService/CreateUser"
	s.apiDelete = "/UserMgmt/RemoveUser"
	s.apiUpdate = "/CDirectoryService/ChangeUser"
	s.apiUpdatePassword = "/UserMgmt/ResetUserPassword"

	return &s
}

// Read function fetches a user from source, including attribute values. Returns error if any
func (o *User) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID

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
	//logger.Debugf("Filled object: %+v", o)

	return nil
}

// Delete function deletes a user and returns a map that contains deletion result
func (o *User) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Create function creates a new user and returns a map that contains creation result
func (o *User) Create() (*restapi.StringResponse, error) {
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
	// Upon successful creation, assign ID
	o.ID = resp.Result

	return resp, nil
}

// Update function updates a existing user and returns a map that contains update result
func (o *User) Update() (*restapi.GenericMapResponse, error) {
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

// ChangePassword function changes user's password
func (o *User) ChangePassword() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	// Mandatory attributes
	queryArg["ID"] = o.ID
	queryArg["newPassword"] = o.Password

	resp, err := o.client.CallBoolAPI(o.apiUpdatePassword, queryArg)
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

// ChangeUserPassword sets new password for a user
func (o *User) ChangeUserPassword(pw string) error {
	// If ID isn't define, find it using username
	if o.ID == "" {
		if o.Name == "" {
			return fmt.Errorf("Missing name for the user object")
		}
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of user %s. %v", o.Name, err)
		}
	}
	o.Password = pw
	_, err := o.ChangePassword()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

// Query function returns a single user object in map format
func (o *User) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM User WHERE 1=1"
	if o.Name != "" {
		query += " AND Username='" + o.Name + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns user ID by name
func (o *User) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("User name must be provided")
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving user '%s': %s", o.Name, err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves user from tenant by name
func (o *User) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return fmt.Errorf("Failed to find ID of user %s. %v", o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// AddToRoles adds user to list of role
func (o *User) AddToRoles(roles []string) error {
	if len(roles) > 0 {
		for _, v := range roles {
			role := NewRole(o.client)
			role.Name = v
			id, err := role.GetIDByName()
			if err != nil {
				return fmt.Errorf("Failed to find ID of role %s. %v", v, err)
			}
			role.ID = id
			resp, err := role.UpdateMembers([]string{o.ID}, "Add", "Users")
			if err != nil || !resp.Success {
				return fmt.Errorf("Error adding user to role: %v", err)
			}
		}
	}
	return nil
}

// DeleteByName deletes a Centrify Directory user by username
func (o *User) DeleteByName() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			return nil, fmt.Errorf("Failed to find ID of user %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
	API to manage user
	https://centrify-dev.readme.io/docs/create-and-manage-cloud-directory-users-_new

	Fetch user
	https://developer.centrify.com/reference#post_cdirectoryservice-exemptuserfrommfa

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"DirectoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"Description": "Test user",
				"ForcePasswordChangeNext": "True",
				"directoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"DisplayName": "Test User",
				"PictureUri": "/UserMgmt/GetUserPicture?id=xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"CloudState": "None",
				"InEverybodyRole": false,
				"OauthClient": false,
				"MobileNumber": "+00 00000000",
				"LastPasswordChangeDate": "/Date(-62135596800000)/",
				"CreateDate": "/Date(1597323830750)/",
				"OfficeNumber": "+00 00000000",
				"CmaRedirectedUser": "admin@examp.com",
				"SubjectToCloudLocks": true,
				"Alias": "example.com",
				"HomeNumber": "+00 00000000",
				"ReportsTo": "admin@examp.com",
				"Name": "testuser@example.com",
				"PreferredCulture": "",
				"Version": "1",
				"Mail": "testuser@example.com",
				"CmaRedirectedUserUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"Uuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"State": "None"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Create user
	https://developer.centrify.com/reference#post_cdirectoryservice-createuser

		Step 1: Create user
		Request body format
		{
			"LoginName": "testuser",
			"Mail": "testuser@centrify.lab",
			"DisplayName": "Test User",
			"Password": "xxxxxxxxx",
			"confirmPassword": "xxxxxxxxx",
			"PasswordNeverExpire": false,
			"ForcePasswordChangeNext": true,
			"InEverybodyRole": false,
			"OauthClient": false,
			"SendEmailInvite": true,
			"Description": "Test user",
			"OfficeNumber": "+00 00000000",
			"HomeNumber": "+00 00000000",
			"MobileNumber": "+00 00000000",
			"fileName": "centrify_logo.jpg",
			"ID": "",
			"state": "None",
			"jsutil-checkbox-2598-inputEl": true,
			"CmaRedirectedUserUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"jsutil-text-2601-inputEl": "admin@examp.com",
			"ReportsTo": "admin@examp.com",
			"Name": "testuser@example.com"
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

		Step 2: Set User Picture

	Update user
	https://developer.centrify.com/reference#post_cdirectoryservice-changeuser

		Request body format
		{
			"LoginName": "testuser",
			"Mail": "testuser@centrify.lab",
			"DisplayName": "Test User test",
			"CloudState": false,
			"PasswordNeverExpire": false,
			"ForcePasswordChangeNext": true,
			"InEverybodyRole": true,
			"OauthClient": false,
			"SendEmailInvite": false,
			"Description": "Test user",
			"OfficeNumber": "+00 00000000",
			"HomeNumber": "+00 00000000",
			"MobileNumber": "+00 00000000",
			"fileName": "",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"state": "None",
			"jsutil-checkbox-2914-inputEl": true,
			"CmaRedirectedUserUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"jsutil-text-2917-inputEl": "admin@examp.com",
			"ReportsTo": "admin@examp.com",
			"Name": "testuser@example.com",
			"Password": "undefined"
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

	Delete user
	https://developer.centrify.com/reference#post_usermgmt-removeuser

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Change password

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"newPassword": "xxxxxx"
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

	Assign rights to role
	https://developer.centrify.com/docs/manage-rolesnew#assigning-rights-to-the-role
	https://developer.centrify.com/reference#post_core-getassignedadministrativerights

		Request body format


		Respond result

*/
