package platform

import (
	"fmt"
	"strings"

	"github.com/marcozj/golang-sdk/enum/resourcetype"
	"github.com/marcozj/golang-sdk/enum/settype"
	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

// Account - Encapsulates a single generic Account
type Account struct {
	vaultObject
	// For password checkout and direct SDK call purpose
	ResourceType string `json:"-"`
	ResourceName string `json:"-"`
	// Account specific APIs
	apiUpdatePassword    string
	apiCheckDelete       string
	apiGetChallenge      string
	apiCheckoutPassword  string
	apiCheckinPassword   string
	apiSetAdminAccount   string
	apiGetAccessKeys     string
	apiRetrieveAccessKey string
	apiDeleteAccessKey   string
	apiAddAccessKey      string
	apiVerifyAccessKey   string

	// Settings menu
	User            string `json:"User,omitempty" schema:"name,omitempty"` // User Name
	Password        string `json:"Password,omitempty" schema:"password,omitempty"`
	Host            string `json:"Host,omitempty" schema:"host_id,omitempty"`
	SSHKeyID        string `json:"SshKeyId,omitempty" schema:"sshkey_id,omitempty"`
	DomainID        string `json:"DomainID,omitempty" schema:"domain_id,omitempty"`
	DatabaseID      string `json:"DatabaseID,omitempty" schema:"database_id,omitempty"`
	CredentialType  string `json:"CredentialType,omitempty" schema:"credential_type,omitempty"` // Password or SshKey
	CredentialID    string `json:"CredentialId,omitempty" schema:"credential_id,omitempty"`
	CloudProviderID string `json:"CloudProviderId,omitempty" schema:"cloudprovider_id,omitempty"`
	IsRootAccount   bool   `json:"IsRootAccount,omitempty" schema:"is_root_account,omitempty"`

	// Policy menu
	UseWheel                       bool            `json:"UseWheel,omitempty" schema:"use_proxy_account,omitempty"` // Use proxy account
	IsManaged                      bool            `json:"IsManaged,omitempty" schema:"managed,omitempty"`          // manage this credential
	Description                    string          `json:"Description,omitempty" schema:"description,omitempty"`
	Status                         string          `json:"Status,omitempty" schema:"status,omitempty"`
	DefaultCheckoutTime            int             `json:"DefaultCheckoutTime,omitempty" schema:"checkout_lifetime,omitempty"` // Checkout lifetime (minutes)
	PasswordCheckoutDefaultProfile string          `json:"PasswordCheckoutDefaultProfile" schema:"default_profile_id"`         // Default Password Checkout Profile (used if no conditions matched)
	ChallengeRules                 *ChallengeRules `json:"PasswordCheckoutRules,omitempty" schema:"challenge_rule,omitempty"`

	IsAdminAccount                     bool            `json:"IsAdminAccount,omitempty" schema:"is_admin_account,omitempty"`
	AccessKeys                         []AccessKey     `json:"AccessKeys,omitempty" schema:"access_key,omitempty"`
	AccessSecretCheckoutDefaultProfile string          `json:"AccessSecretCheckoutDefaultProfile,omitempty" schema:"access_secret_checkout_default_profile_id,omitempty"`
	AccessSecretCheckoutRules          *ChallengeRules `json:"AccessSecretCheckoutRules,omitempty" schema:"access_secret_checkout_rule,omitempty"`
}

// NewAccount is Account constructor
func NewAccount(c *restapi.RestClient) *Account {
	s := Account{}
	s.client = c
	s.ValidPermissions = ValidPermissionMap.Account
	s.SetType = settype.Account.String()
	s.apiRead = "/ServerManage/GetAllAccountInformation"
	s.apiCreate = "/ServerManage/AddAccount"
	s.apiDelete = "/ServerManage/DeleteAccount"
	s.apiUpdate = "/ServerManage/UpdateAccount"
	s.apiUpdatePassword = "/ServerManage/UpdatePassword"
	s.apiCheckDelete = "/ServerManage/CanDeleteAccount"
	s.apiGetChallenge = "/ServerManage/GetAccountChallenges"
	s.apiCheckoutPassword = "/ServerManage/CheckoutPassword"
	s.apiCheckinPassword = "/ServerManage/CheckinPassword"
	s.apiPermissions = "/ServerManage/SetAccountPermissions"
	s.apiSetAdminAccount = "/ServerManage/SetAdministrativeAccounts"
	s.apiGetAccessKeys = "/Aws/GetAccessKeys"
	s.apiRetrieveAccessKey = "/Aws/RetrieveAccessKey"
	s.apiDeleteAccessKey = "/Aws/DeleteAccessKey"
	s.apiAddAccessKey = "/Aws/AddAccessKey"
	s.apiVerifyAccessKey = "/Aws/VerifyAccessKeyForUserAccount"

	return &s
}

// Read function fetches a Account from source, including attribute values. Returns error if any
func (o *Account) Read() error {
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

	var va = resp.Result["VaultAccount"].(map[string]interface{})
	var row = va["Row"].(map[string]interface{})

	mapToStruct(o, row)

	// Get password checkout profile information
	resp, err = o.client.CallGenericMapAPI(o.apiGetChallenge, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}
	if v, ok := resp.Result["PasswordCheckoutDefaultProfile"]; ok {
		o.PasswordCheckoutDefaultProfile = v.(string)
	}

	// Fill challenge rules
	if v, ok := resp.Result["PasswordCheckoutRules"]; ok {
		challengerules := &ChallengeRules{}
		mapToStruct(challengerules, v.(map[string]interface{}))
		o.ChallengeRules = challengerules
	}

	// Fill AWS Access key
	if o.CloudProviderID != "" {
		keys, err := o.GetAccessKeys()
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}
		for _, key := range keys {
			o.AccessKeys = append(o.AccessKeys, key)
		}
	}

	return nil
}

// Create function creates a new Account and returns a map that contains creation result
func (o *Account) Create() (*restapi.StringResponse, error) {
	// Resolve host id if not provided
	err := o.resolveHostID()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	// Special handling of password checkout profile
	//if queryArg["PasswordCheckoutDefaultProfile"] != "" {
	queryArg["updateChallenges"] = true
	//}

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

// Delete function deletes a Account and returns a map that contains deletion result
func (o *Account) Delete() (*restapi.BoolResponse, error) {
	return o.deleteObjectBoolAPI("")
}

// Update function updates an existing Account and returns a map that contains update result
func (o *Account) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	// Resolve host id if not provided
	err := o.resolveHostID()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	var queryArg = make(map[string]interface{})
	queryArg, err = generateRequestMap(o)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	// Special handling of password checkout profile
	//if queryArg["PasswordCheckoutDefaultProfile"] != "" {
	queryArg["updateChallenges"] = true
	//}

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

// ChangePassword function updates an existing Account password and returns a map that contains update result
func (o *Account) ChangePassword() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	// Mandatory attributes
	queryArg["ID"] = o.ID
	queryArg["Password"] = o.Password

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

// ValidateCredentialType checks credential type matches password or sshkey setting
func (o *Account) ValidateCredentialType() error {
	if o.CredentialType == "Password" && o.Password == "" {
		return fmt.Errorf("Credential type is 'Password' but password isn't set")
	}
	if o.CredentialType == "SshKey" && o.SSHKeyID == "" {
		return fmt.Errorf("Credential type is 'SSHKey' but Sshkey_id isn't set")
	}
	return nil
}

// Query function returns a single Account object in map format
func (o *Account) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM VaultAccount WHERE 1=1"
	if o.User != "" {
		query += " AND User='" + o.User + "'"
	}
	if o.Host != "" {
		query += " AND Host='" + o.Host + "'"
	}
	if o.DatabaseID != "" {
		query += " AND DatabaseID='" + o.DatabaseID + "'"
	}
	if o.DomainID != "" {
		query += " AND DomainID='" + o.DomainID + "'"
	}
	if o.CloudProviderID != "" {
		query += " AND CloudProviderId='" + o.CloudProviderID + "'"
	}

	return queryVaultObject(o.client, query)
}

// CheckoutPassword checks out account password from vault
func (o *Account) checkoutPassword() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Description"] = "Checkout by golang SDK"

	resp, err := o.client.CallGenericMapAPI(o.apiCheckoutPassword, queryArg)
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

// CheckoutPassword checks out account password from vault
// Returns actual password, coid or error
func (o *Account) CheckoutPassword(checkin bool) (string, error) {
	// To checkout account password, we must know its ID
	// In order to know the ID of the account, we must know username + Host/DatabaseID/DomainID
	if o.ID == "" {
		// if ID is unknown, try to find out using User, ResourceType and ResourceName
		_, err := o.getResourceID()
		if err != nil {
			logger.Errorf(err.Error())
			return "", err
		}
		acctresult, err := o.Query()
		if err != nil {
			logger.Errorf(err.Error())
			return "", fmt.Errorf("Error retrieving account object: %s", err)
		}
		o.ID = acctresult["ID"].(string)
	}
	// Check again if ID is known
	if o.ID == "" {
		return "", fmt.Errorf("Missing ID for account %s in %s with type %s", o.User, o.ResourceName, o.ResourceType)
	}

	// Checking out password
	reply, err := o.checkoutPassword()
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}

	if pw, ok := reply.Result["Password"]; ok {
		if checkin {
			coid := reply.Result["COID"]
			if coid != nil {
				result, err := o.CheckinPassword(coid.(string))
				if err != nil {
					logger.Errorf(err.Error())
					return pw.(string), err
				}
				if !result.Success {
					return pw.(string), fmt.Errorf(result.Message)
				}
			} else {
				return pw.(string), fmt.Errorf("No COID returned from checkout")
			}
		}
		return pw.(string), nil
	}
	return "", fmt.Errorf("Password checkout call doesn't contain password")
}

// CheckinPassword checks in an checked out account password
func (o *Account) CheckinPassword(coid string) (*restapi.BoolResponse, error) {
	if coid == "" {
		errormsg := fmt.Sprintf("Missing COID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = coid

	resp, err := o.client.CallBoolAPI(o.apiCheckinPassword, queryArg)
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

// RetrieveSSHKey retrieves SSH key from an account
func (o *Account) RetrieveSSHKey(keytype string, passphrase string) (string, error) {
	if o.ID == "" {
		// if ID is unknown, try to find out using User, ResourceType and ResourceName
		_, err := o.getResourceID()
		if err != nil {
			logger.Errorf(err.Error())
			return "", err
		}
		acctresult, err := o.Query()
		if err != nil {
			return "", fmt.Errorf("Error retrieving account object: %s", err)
		}
		o.ID = acctresult["ID"].(string)
		o.CredentialID = acctresult["CredentialId"].(string)
	}
	// Check again if ID is known
	if o.ID == "" {
		return "", fmt.Errorf("Missing ID for account %s in %s with type %s", o.User, o.ResourceName, o.ResourceType)
	}
	if o.CredentialID == "" {
		return "", fmt.Errorf("SSH Key ID not found for account %s", o.User)
	}

	sshkey := NewSSHKey(o.client)
	sshkey.ID = o.CredentialID
	sshkey.KeyPairType = keytype
	sshkey.Passphrase = passphrase
	sshkey.KeyFormat = "PEM"
	thekey, err := sshkey.RetriveSSHKey()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieve sshkey. %v", err)
	}

	return thekey, nil
}

func (o *Account) getResourceID() (string, error) {
	var resourceID string
	if o.Host != "" {
		resourceID = o.Host
	} else if o.DomainID != "" {
		resourceID = o.DomainID
	} else if o.DatabaseID != "" {
		resourceID = o.DatabaseID
	} else if o.CloudProviderID != "" {
		resourceID = o.CloudProviderID
	}

	if resourceID == "" {
		if o.User != "" && o.ResourceType != "" && o.ResourceName != "" {
			// Get resource ID
			switch strings.ToLower(o.ResourceType) {
			case resourcetype.System.String():
				resource := NewSystem(o.client)
				resource.Name = o.ResourceName
				result, err := resource.Query()
				if err != nil {
					logger.Errorf(err.Error())
					return "", fmt.Errorf("Error retrieving system object: %s", err)
				}
				resourceID = result["ID"].(string)
				o.Host = resourceID
			case resourcetype.Database.String():
				resource := NewDatabase(o.client)
				resource.Name = o.ResourceName
				result, err := resource.Query()
				if err != nil {
					logger.Errorf(err.Error())
					return "", fmt.Errorf("Error retrieving database object: %s", err)
				}
				resourceID = result["ID"].(string)
				o.DatabaseID = resourceID
			case resourcetype.Domain.String():
				resource := NewDomain(o.client)
				resource.Name = o.ResourceName
				result, err := resource.Query()
				if err != nil {
					logger.Errorf(err.Error())
					return "", fmt.Errorf("Error retrieving domain object: %s", err)
				}
				resourceID = result["ID"].(string)
				o.DomainID = resourceID
			case resourcetype.CloudProvider.String():
				resource := NewCloudProvider(o.client)
				resource.Name = o.ResourceName
				result, err := resource.Query()
				if err != nil {
					logger.Errorf(err.Error())
					return "", fmt.Errorf("Error retrieving domain object: %s", err)
				}
				resourceID = result["ID"].(string)
				o.CloudProviderID = resourceID
			default:
				return "", fmt.Errorf("Invalid resource type: %s", o.ResourceType)
			}
		} else {
			return "", fmt.Errorf("Missing required attributes User: %s, ResourceType: %s, ResourceName: %s", o.User, o.ResourceType, o.ResourceName)
		}
	}
	return resourceID, nil
}

// ResolveValidPermissions resolves valid permission according to account type
func (o *Account) ResolveValidPermissions() {
	if o.Host != "" {
		o.ValidPermissions = ValidPermissionMap.Account
	} else if o.DomainID != "" {
		o.ValidPermissions = ValidPermissionMap.DomainAccount
	} else if o.DatabaseID != "" {
		o.ValidPermissions = ValidPermissionMap.DBAccount
	} else if o.CloudProviderID != "" {
		o.ValidPermissions = ValidPermissionMap.CloudAccount
	}
}

// SetAdminAccount set this account as admin account
func (o *Account) SetAdminAccount(enable bool) error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	if enable {
		queryArg["PVID"] = o.ID
	}
	queryArg["Systems"] = []string{o.Host}
	resp, err := o.client.CallGenericMapAPI(o.apiSetAdminAccount, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}

	return nil
}

// VerifyAccessKey verifies that access key is valid against AWS
func (o *Account) VerifyAccessKey(key AccessKey) error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["AccountId"] = o.ID
	queryArg["CloudProviderId"] = o.CloudProviderID
	queryArg["User"] = o.User
	queryArg["UserName"] = o.User
	queryArg["AccessKeyId"] = key.AccessKeyID
	queryArg["SecretAccessKey"] = key.SecretAccessKey

	resp, err := o.client.CallStringAPI(o.apiVerifyAccessKey, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}

	return nil
}

// AddAccessKey adds access key into this account
func (o *Account) AddAccessKey(key AccessKey) error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["AccountId"] = o.ID
	queryArg["User"] = o.Name
	queryArg["AccessKeyId"] = key.AccessKeyID
	queryArg["SecretAccessKey"] = key.SecretAccessKey

	resp, err := o.client.CallStringAPI(o.apiAddAccessKey, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}

	return nil
}

// SafeAddAccessKey verifies then adds access key
func (o *Account) SafeAddAccessKey(key AccessKey) error {
	err := o.VerifyAccessKey(key)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	err = o.AddAccessKey(key)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

// GetAccessKeys get all access key entries
func (o *Account) GetAccessKeys() ([]AccessKey, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["ID"] = o.ID
	queryArg["Args"] = subArgs

	resp, err := o.client.CallSliceAPI(o.apiGetAccessKeys, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}

	keys := []AccessKey{}
	logger.Debugf("Get key response: %+v", resp)
	for _, p := range resp.Result {
		key := &AccessKey{}
		mapToStruct(key, p.(map[string]interface{}))
		logger.Debugf("Filled key: %+v", key)
		keys = append(keys, *key)
	}

	return keys, nil
}

// DeleteAccessKey deletes an IAM access key
func (o *Account) DeleteAccessKey(id string) error {
	if id == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}

	var queryArg = make(map[string]interface{})
	queryArg["AccessKeyRowkey"] = id
	resp, err := o.client.CallStringAPI(o.apiDeleteAccessKey, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return fmt.Errorf(errmsg)
	}
	return nil
}

// RetrieveAccessKey retrieves secret access key
func (o *Account) RetrieveAccessKey(accessKeyID string) (string, error) {
	if accessKeyID == "" {
		errormsg := fmt.Sprintf("Missing Access Key ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return "", fmt.Errorf(errormsg)
	}

	// To retrieve secret key, we must know acccount ID
	// In order to know the ID of the account, we must know username + CloudProviderID
	if o.ID == "" {
		// if ID is unknown, try to find out using User, ResourceType and ResourceName
		_, err := o.getResourceID()
		if err != nil {
			logger.Errorf(err.Error())
			return "", err
		}
		acctresult, err := o.Query()
		if err != nil {
			logger.Errorf(err.Error())
			return "", fmt.Errorf("Error retrieving account object: %s", err)
		}
		o.ID = acctresult["ID"].(string)
	}
	// Check again if ID is known
	if o.ID == "" {
		return "", fmt.Errorf("Missing ID for account %s in %s with type %s", o.User, o.ResourceName, o.ResourceType)
	}

	keys, err := o.GetAccessKeys()
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	var id string
	for _, key := range keys {
		if accessKeyID == key.AccessKeyID {
			id = key.ID
		}
	}

	var queryArg = make(map[string]interface{})
	queryArg["AccessKeyRowkey"] = id
	resp, err := o.client.CallGenericMapAPI(o.apiRetrieveAccessKey, queryArg)
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	if !resp.Success {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return "", fmt.Errorf(errmsg)
	}

	return resp.Result["SecretAccessKey"].(string), nil
}

// GetIDByName returns vault object ID by name
func (o *Account) GetIDByName() (string, error) {
	if o.User == "" {
		return "", fmt.Errorf("%s name must be provided", GetVarType(o))
	}

	_, err := o.getResourceID()
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("Error retrieving %s %s: %s", GetVarType(o), o.User, err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves vault object from tenant by name
func (o *Account) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("Failed to find ID of %s %s. %v", GetVarType(o), o.User, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a DesktopApp by name
func (o *Account) DeleteByName() (*restapi.BoolResponse, error) {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("Failed to find ID of DesktopApp %s. %v", o.Name, err)
		}
	}
	resp, err := o.Delete()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *Account) resolveHostID() error {
	if o.Host == "" && o.DomainID == "" && o.DatabaseID == "" && o.CloudProviderID == "" {
		_, err := o.getResourceID()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	API to manage vault account

	Read Account

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"VaultCheckout": null,
				"Workflow": {
					"WorkflowApprover": ...
					"WorkflowApprovers": ...
					"WorkflowEnabled": true,
					"PasswordResettable": false,
					"AccountUnlockable": false
				},
				"VaultAccount": {
					"Row": {
						"IsQuickStartAccount": null,
						"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
						"ActiveSessions": 0,
						"Host": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
						"LastChange": "/Date(1588075662725)/",
						"DefaultCheckoutTime": null,
						"EffectiveWorkflowEnabled": null,
						"DueBack": null,
						"ProxyUser": "admin",
						"CanManage": false,
						"Mode": null,
						"Status": "",
						"User": "adminuser",
						"PasswordResetRetryCount": 0,
						"DomainID": null,
						"WorkflowApproversList": null,
						"KmipId": null,
						"PasswordResetLastError": "",
						"AccountRights": "0000000000000000000000000000000000000000000000000000011101111011",
						"WorkflowApprover": ...
						"NeedsPasswordReset": "NotNeeded",
						"HealthError": "OK",
						"SessionType": "Ssh",
						"MissingPassword": false,
						"CredentialId": null,
						"MPParent": null,
						"IsPrivileged": null,
						"WorkflowApprovers": ...
						"Description": "",
						"FQDN": "192.168.18.53",
						"IsManaged": true,
						"EffectiveWorkflowApprover": null,
						"Name": "FortiManager",
						"Rights": "Login, Naked, Manage, Owner, Delete, UpdatePassword, RotatePassword, View, FileTransfer",
						"OwnerId": null,
						"UseWheel": false,
						"WorkflowDefaultOptions": null,
						"DeviceID": null,
						"LastHealthCheck": "/Date(1588075662725)/",
						"DesktopApps": [],
						"CredentialType": "Password",
						"ComputerClass": "CustomSsh",
						"ActiveCheckouts": 0,
						"WorkflowEnabled": null,
						"DiscoveredTime": null,
						"_MatchFilter": null,
						"Healthy": "OK",
						"OwnerName": null,
						"EffectiveWorkflowApprovers": null,
						"UserDisplayName": "adminuser (FortiManager)",
						"IsFavorite": false,
						"DatabaseID": null,
						"NeedPassword": false
					},
					"Entities": [
						{
							"Type": "Server",
							"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"IsForeignKey": true
						},
						{
							"Type": "VaultAccount",
							"Key": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
							"IsForeignKey": false
						}
					]
				},
				"RelatedResource": {
					"Row": {
						...
					},
					"Entities": [
						...
					]
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

	Create Account
	https://developer.centrify.com/reference#post_servermanage-addaccount

		Request body format
		{
			"User": "testaccount",
			"CredentialType": "Password",
			"IsManaged": false,
			"Password": "xxxxxxxxxx",
			"UseWheel": false,
			"Description": "afdasdfasfd",
			"Host": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}
		or
		{
			"User": "account2",
			"CredentialType": "SshKey",
			"IsManaged": false,
			"undefined": "false",
			"SshKeyName": "my_ami_key",
			"SshKeyId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Host": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"Description": ""
		}
		or
		{
			"User": "sa",
			"Password": "xxxxxxxxx",
			"UseWheel": false,
			"Description": "sa account",
			"IsManaged": false,
			"DatabaseID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

	Update Account
	https://developer.centrify.com/reference#post_servermanage-updateaccount

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"ActiveSessions": 0,
			"Host": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"CanManage": true,
			"Status": "Unknown",
			"User": "testaccount",
			"AccountRights": "0000000000000000000000000000000000000000000000000000011101111011",
			"WorkflowApprover": ...
			"NeedsPasswordReset": "NotNeeded",
			"HealthError": "NoManagementChannelAvailable",
			"SessionType": "Rdp",
			"MissingPassword": false,
			"WorkflowApprovers": ...
			"Description": "Test account",
			"FQDN": "192.168.2.3",
			"IsManaged": false,
			"Name": "Windows 01",
			"Rights": "Login, Naked, Manage, Owner, Delete, UpdatePassword, RotatePassword, View, FileTransfer",
			"DesktopApps": [],
			"CredentialType": "Password",
			"ComputerClass": "Windows",
			"ActiveCheckouts": 0,
			"Healthy": "Unknown",
			"UserDisplayName": "testaccount (Windows 01)",
			"IsFavorite": false,
			"NeedPassword": false,
			"IsAdministrativeAccount": false,
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"updateChallenges": true,
			"WorkflowEnabled": null
		}

		Respond result
		{
			"success": true,
			"Result": {
				"PVID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

	Delete Account
	https://developer.centrify.com/reference#post_servermanage-deleteaccount

		Check if Account can be deleted
		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"RRFormat": true
		}

		Respond result
		{
			"success": true,
			"Result": {
				"can": true,
				"why": null
			},
			"Message": null,
			"MessageID": null,
			"Exception": null,
			"ErrorID": null,
			"ErrorCode": null,
			"IsSoftError": false,
			"InnerExceptions": null
		}

		Deleting account
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

		Upate password
		https://developer.centrify.com/reference#post_servermanage-updatepassword

		Request body format
		{
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
		}

		Respond result
		{
			"success": true,
			"Result": {
				"PasswordCheckoutDefaultProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"PasswordCheckoutRules": {
					"_UniqueKey": "Condition",
					"_Value": [
						{
							"Conditions": [
								{
									"Prop": "IdentityCookie",
									"Op": "OpNotExists"
								}
							],
							"ProfileId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
						}
					],
					"Enabled": true,
					"_Type": "RowSet"
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

	Set Admin Account
		Request body format
		{
			"Systems": [
				"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
			],
			"PVID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

/*
{
    "ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
    "ActiveSessions": 0,
    "Host": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
    "ProxyUser": "root",
    "CanManage": false,
    "User": "user1",
    "PasswordResetRetryCount": 0,
    "WorkflowApproversList": "[
	{
		\"DisplayName\":\"admin\",
		\"ObjectType\":\"User\",
		\"DistinguishedName\":\"admin@example.com\",
		\"DirectoryServiceUuid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
		\"SystemName\":\"admin@example.com\",
		\"ServiceInstance\":\"CDS\",
		\"Locked\":false,
		\"InternalName\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
		\"StatusEnum\":\"Active\",
		\"ServiceInstanceLocalized\":\"Centrify Directory\",
		\"ServiceType\":\"CDS\",
		\"EMail\":\"admin@demo.lab\",
		\"Status\":\"Active\",
		\"Enabled\":true,\
		"Name\":\"admin@example.com\",
		\"Guid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
		\"Type\":\"User\",
		\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
		\"Principal\":\"admin@example.com\",
		\"PType\":\"User\",
		\"OptionsSelector\":true
	},
	{
		\"ReadOnly\":false,
		\"Description\":\"AD accounts that are granted local administrator access to non-domain joined machines.\",
		\"DirectoryServiceUuid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
		\"RoleType\":\"PrincipalList\",
		\"_ID\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
		\"Name\":\"LAB Cloud Local Admins\",
		\"Guid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
		\"Type\":\"Role\",
		\"ObjectType\":\"Role\",
		\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
		\"Principal\":\"LAB Cloud Local Admins\",
		\"PType\":\"Role\"
	},
	{
		\"Type\":\"Manager\",
		\"NoManagerAction\":\"useBackup\",
		\"BackupApprover\":
		{
			\"ReadOnly\":false,
			\"Description\":\"Machines and users who are enforced MFA for direct access without going through PAS.\",
			\"DirectoryServiceUuid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"RoleType\":\"PrincipalList\",
			\"_ID\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"Name\":\"LAB MFA Machines & Users\",
			\"Guid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"Type\":\"Role\",
			\"ObjectType\":\"Role\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
			\"Principal\":\"LAB MFA Machines & Users\",
			\"PType\":\"Role\"
		},
		\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\"
	}
	]",
    "AccountRights": "0000000000000000000000000000000000000000000000000000011101111011",
    "NeedsPasswordReset": "NotNeeded",
    "HealthError": "OK",
    "SessionType": "Ssh",
    "MissingPassword": false,
    "WorkflowApprovers": "[
		{
			\"DisplayName\":\"mspadmin\",
			\"ObjectType\":\"User\",
			\"DistinguishedName\":\"admin@example.com\",
			\"DirectoryServiceUuid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"SystemName\":\"admin@example.com\",
			\"ServiceInstance\":\"CDS\",
			\"Locked\":false,
			\"InternalName\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",\
			"StatusEnum\":\"Active\",
			\"ServiceInstanceLocalized\":\"Centrify Directory\",
			\"ServiceType\":\"CDS\",
			\"EMail\":\"admin@demo.lab\",
			\"Status\":\"Active\",
			\"Enabled\":true,
			\"Name\":\"admin@example.com\",
			\"Guid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"Type\":\"User\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
			\"Principal\":\"admin@example.com\",
			\"PType\":\"User\",
			\"OptionsSelector\":true
		},
		{
			\"ReadOnly\":false,
			\"Description\":\"AD accounts that are granted local administrator access to non-domain joined machines.\",
			\"DirectoryServiceUuid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"RoleType\":\"PrincipalList\",
			\"_ID\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"Name\":\"LAB Cloud Local Admins\",
			\"Guid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
			\"Type\":\"Role\",
			\"ObjectType\":\"Role\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
			\"Principal\":\"LAB Cloud Local Admins\",
			\"PType\":\"Role\"
		},
		{
			\"Type\":\"Manager\",
			\"NoManagerAction\":\"useBackup\",
			\"BackupApprover\":
			{
				\"ReadOnly\":false,
				\"Description\":\"Machines and users who are enforced MFA for direct access without going through PAS.\",
				\"DirectoryServiceUuid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
				\"RoleType\":\"PrincipalList\",
				\"_ID\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
				\"Name\":\"LAB MFA Machines & Users\",
				\"Guid\":\"xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx\",
				\"Type\":\"Role\",
				\"ObjectType\":\"Role\",
				\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128\",
				\"Principal\":\"LAB MFA Machines & Users\",
				\"PType\":\"Role\"
			},
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\"
		}]",
	"WorkflowApprovers": "[
		{
			\"Type\":\"Manager\",
			\"NoManagerAction\":\"approve\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
			\"OptionsSelector\":
		}
	]",
	"WorkflowApprovers": "[
		{
			\"Type\":\"Manager\",
			\"NoManagerAction\":\"deny\",
			\"Type-generated-field\":\"/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/user_icon_sml.png?_ver=1598087128\",
			\"OptionsSelector\":true
		}
	]",
    "FQDN": "192.168.18.15",
    "IsManaged": true,
    "Name": "XML File",
    "Rights": "Login, Naked, Manage, Owner, Delete, UpdatePassword, RotatePassword, View, FileTransfer",
    "UseWheel": true,
    "WorkflowDefaultOptions": "{\"GrantMin\":60}",
    "DesktopApps": [],
    "CredentialType": "Password",
    "ComputerClass": "CustomSsh",
    "ActiveCheckouts": 0,
    "WorkflowEnabled": true,
    "Healthy": "OK",
    "UserDisplayName": "user1 (XML File)",
    "IsFavorite": false,
    "NeedPassword": false,
    "IsAdministrativeAccount": false,
    "NoManagerAction": "useBackup",
    "BackupApprover": {
        "ReadOnly": false,
        "Description": "Machines and users who are enforced MFA for direct access without going through PAS.",
        "DirectoryServiceUuid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
        "RoleType": "PrincipalList",
        "_ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
        "Name": "LAB MFA Machines & Users",
        "Guid": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
        "Type": "Role",
        "ObjectType": "Role",
        "Type-generated-field": "/vfslow/lib/ui/../uibuild/compiled/centrify/production/resources/images/entities/group_icon_sml.png?_ver=1598087128",
        "Principal": "LAB MFA Machines & Users",
        "PType": "Role"
    },
    "WorkflowSent": true,
    "_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
    "updateChallenges": true,
    "Description": ""
}

Get AWS Access Keys
	Request body format
	{
		"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
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
		"Result": [
			{
				"_TableName": "accesskeys",
				"_Timestamp": "/Date(1610267674670)/",
				"Created": "/Date(1599710110000)/",
				"_PartitionKey": "ABC0751",
				"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"_entitycontext": "*",
				"AccountId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
				"AccessKeyId": "XXXXXXXXXXXXXX",
				"_metadata": {
					"Version": 1,
					"IndexingVersion": 1
				}
			}
		],
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Retrieve AWS Access Keys
	Request body format
	{
		"AccessKeyRowkey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"_TableName": "accesskeys",
			"_encryptkeyid": "XXXXXX",
			"SecretAccessKey": "XXXXXXXXXXXXXXXXXX",
			"Created": "/Date(1599710110000)/",
			"_PartitionKey": "XXXXXX",
			"ID": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_entitycontext": "W/\"datetime'2021-01-10T08%3A34%3A34.6256835Z'\"",
			"AccountId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"AccessKeyId": "XXXXXXXXXXXXXX"
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Delete AWS Access Key
	Request body format
	{
		"AccessKeyRowkey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

Verify AWS Access Key
	Request body format
	{
		"UserName": "pas_test",
		"CloudProviderId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"AccessKeyId": "XXXXXXXXXXX",
		"SecretAccessKey": "XXXXXXXXXXXXXXXXXXXX",
		"User": "pas_test",
		"AccountId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": "pas_test",
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

Add AWS Access Key
	Request body format
	{
		"AccessKeyId": "XXXXXXXXXXX",
		"SecretAccessKey": "XXXXXXXXXXXXXXXXXXXX",
		"User": "xxxxxx",
		"AccountId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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
