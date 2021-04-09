package platform

import (
	"fmt"

	logger "github.com/marcozj/golang-sdk/logging"
	"github.com/marcozj/golang-sdk/restapi"
)

type SamlWebApp struct {
	WebApp
	apiGetSpMetadataFromUrl string

	//TemplateName     string `json:"TemplateName,omitempty" schema:"template_name,omitempty"`     // "Generic SAML", "AWSConsoleSAML", "ClouderaSAML", "CloudLock SAML", "ConfluenceServerSAML", "Dome9Saml", "GitHubEnterpriseSAML", "JIRACloudSAML", "JIRAServerSAML", "PaloAltoNetworksSAML", "SplunkOnPremSAML", "SumoLogicSAML"
	CorpIdentifier   string `json:"CorpIdentifier,omitempty" schema:"corp_identifier,omitempty"` // Used for AWS (AWS Account ID), JIRACloudSAML (Jira Cloud Subdomain)
	AdditionalField1 string `json:"AdditionalField1,omitempty" schema:"app_entity_id,omitempty"` // Used for ClouderaSAML (Cloudera Entity ID), JIRACloudSAML (SP Entity ID)
	ServiceName      string `json:"ServiceName,omitempty" schema:"application_id,omitempty"`
	IdpMetadataUrl   string `json:"IdpMetadataUrl,omitempty" schema:"idp_metadata_url,omitempty"`
	// Trust menu
	SpMetadataUrl         string `json:"SpMetadataUrl,omitempty" schema:"sp_metadata_url,omitempty"`
	SpConfigMethod        int    `json:"SpConfigMethod" schema:"sp_config_method"`
	SpMetadataXml         string `json:"SpMetadataXml,omitempty" schema:"sp_metadata_xml,omitempty"`
	Audience              string `json:"Audience,omitempty" schema:"sp_entity_id,omitempty"`                  // SP Entity ID / Issuer / Audience
	ACS_Url               string `json:"Url,omitempty" schema:"acs_url,omitempty"`                            // Assertion Consumer Service (ACS) URL
	RecipientSameAsAcsUrl bool   `json:"RecipientSameAsAcsUrl" schema:"recipient_sameas_acs_url"`             // Recipient same as ACS URL
	Recipient             string `json:"Recipient,omitempty" schema:"recipient,omitempty"`                    // Recipient
	WantAssertionsSigned  bool   `json:"WantAssertionsSigned" schema:"sign_assertion"`                        // Sign Assertion
	NameIDFormat          string `json:"NameIDFormat,omitempty" schema:"name_id_format,omitempty"`            // NameID Format
	SpSingleLogoutUrl     string `json:"SpSingleLogoutUrl,omitempty" schema:"sp_single_logout_url,omitempty"` // Single Logout URL
	EncryptAssertion      bool   `json:"EncryptAssertion,omitempty" schema:"encrypt_assertion,omitempty"`     // Encrypt SAML Response Assertion
	//EncryptionThumbprint string
	RelayState        string `json:"RelayState,omitempty" schema:"relay_state,omitempty"`                // Relay State
	AuthnContextClass string `json:"AuthnContextClass,omitempty" schema:"authn_context_class,omitempty"` // Authentication Context Class
	// SAML Response menu
	SamlAttributes     []SamlAttribute `json:"SamlAttributes,omitempty" schema:"saml_attribute,omitempty"` // SAML Response attributes
	SamlResponseScript string          `json:"Script,omitempty" schema:"saml_response_script,omitempty"`   // SAML Response Custom Logic
	SamlScript         string          `json:"SamlScript,omitempty" schema:"saml_script,omitempty"`
}

type SamlAttribute struct {
	Name  string `json:"Name,omitempty" schema:"name,omitempty"`
	Value string `json:"Value,omitempty" schema:"value,omitempty"`
}

func NewSamlWebApp(c *restapi.RestClient) *SamlWebApp {
	webapp := newWebpp(c)
	s := SamlWebApp{}
	s.WebApp = *webapp
	s.apiGetSpMetadataFromUrl = "/saasManage/GetSpMetadataFromUrl"

	return &s
}

func (o *SamlWebApp) Read() error {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return fmt.Errorf(errormsg)
	}
	var queryArg = make(map[string]interface{})
	queryArg["_RowKey"] = o.ID

	logger.Debugf("Generated Map for Read(): %+v", queryArg)
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

	return nil
}

// Update function updates an existing WebApp and returns a map that contains update result
func (o *SamlWebApp) Update() (*restapi.GenericMapResponse, error) {
	if o.ID == "" {
		errormsg := fmt.Sprintf("Missing ID for %s", GetVarType(0))
		logger.Errorf(errormsg)
		return nil, fmt.Errorf(errormsg)
	}

	err := o.processSpMetaData()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}

	err = o.processWorkflow()
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
	queryArg["_RowKey"] = o.ID

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

func (o *SamlWebApp) processWorkflow() error {
	// Resolve guid of each approver
	if o.WorkflowEnabled && o.WorkflowApproverList != nil {
		err := ResolveWorkflowApprovers(o.client, o.WorkflowApproverList)
		if err != nil {
			return err
		}
		// Due to historical reason, WorkflowSettings attribute is not in json format rather it is in string so need to perform conversion
		// Convert approvers from struct to string so that it can be assigned to the actual attribute used for privision.
		wfApprovers := FlattenWorkflowApprovers(o.WorkflowApproverList)
		o.WorkflowSettings = "{\"WorkflowApprover\":" + wfApprovers + "}"
	}
	return nil
}

func (o *SamlWebApp) processSpMetaData() error {

	if o.SpMetadataUrl != "" {
		var queryArg = make(map[string]interface{})
		queryArg["_RowKey"] = o.ID
		queryArg["Url"] = o.SpMetadataUrl
		resp, err := o.client.CallGenericMapAPI(o.apiGetSpMetadataFromUrl, queryArg)
		if err != nil {
			logger.Errorf(err.Error())
			return err
		}

		if !resp.Success {
			errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
			logger.Errorf(errmsg)
			return fmt.Errorf(errmsg)
		} else {
			mapToStruct(o, resp.Result)
		}
	}

	if o.Recipient != "" {
		o.RecipientSameAsAcsUrl = false
	} else {
		o.RecipientSameAsAcsUrl = true
	}

	return nil
}

// GetIDByName returns vault object ID by name
func (o *SamlWebApp) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("%s name must be provided", GetVarType(o))
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("error retrieving %s: %s", GetVarType(o), err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves vault object from tenant by name
func (o *SamlWebApp) GetByName() error {
	if o.ID == "" {
		_, err := o.GetIDByName()
		if err != nil {
			logger.Errorf(err.Error())
			return fmt.Errorf("failed to find ID of %s %s. %v", GetVarType(o), o.Name, err)
		}
	}

	err := o.Read()
	if err != nil {
		return err
	}
	return nil
}

// Query function returns a single WebApp object in map format
func (o *SamlWebApp) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Application WHERE 1=1 AND AppType='Web' AND WebAppType='Saml'"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.ServiceName != "" {
		query += " AND ServiceName='" + o.ServiceName + "'"
	}
	if o.CorpIdentifier != "" {
		query += " AND CorpIdentifier='" + o.CorpIdentifier + "'"
	}
	if o.AdditionalField1 != "" {
		query += " AND AdditionalField1='" + o.AdditionalField1 + "'"
	}
	if o.Audience != "" {
		query += " AND Audience='" + o.Audience + "'"
	}

	return queryVaultObject(o.client, query)
}

/*
Fetch web app
	Request body format
	{
		"_RowKey": "7d8a477a-e891-4ca8-b866-c606e732f661",
		"RRFormat": true,
		"Args": {
			"PageNumber": 1,
			"Limit": 1,
			"PageSize": 1,
			"Caching": -1
		}
	}

	Respond result
	{
		"success": true,
		"Result": {
			"IdpMetadataUrl": "https://xxxxx.my.centrify.net/saasManage/DownloadSAMLMetadataForApp?appkey=xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx&customerid=ABC0000",
			"RelayState": "",
			"ShowAdditionalField1": false,
			"AuthnContextClassShow": true,
			"AuthnContextClassHelpTip": "Select the Authentication Context Class that your Service Provider specifies to use. If SP does not specify one, select 'unspecified'.",
			"RecipientSameAsAcsUrlShow": true,
			"SignOutUrlHelptip": "Configure this in your SAML application if you want to automatically sign users out of the User portal when users sign out of your SAML application.",
			"AppPluginUrl": "",
			"IsGatewayAllowed": true,
			"IssuerReadOnly": false,
			"IsTestApp": false,
			"Icon": "/vfslow/lib/application/icons/GenericSaml.svg",
			"Featured": false,
			"SpConfigMethod": 1,
			"AudienceShow": true,
			"SamlVersionReadOnly": false,
			"ShowUrl": true,
			"SamlAttributeValueOptions": [
				{
					"options": [
						{
							"value": "LoginUser.CanonicalName",
							"text": "CanonicalName"
						},
						{
							"value": "LoginUser.Description",
							"text": "Description"
						},
						{
							"value": "LoginUser.DisplayName",
							"text": "DisplayName"
						},
						{
							"value": "LoginUser.EffectiveGroupDNs",
							"text": "EffectiveGroupDNs"
						},
						{
							"value": "LoginUser.EffectiveGroupNames",
							"text": "EffectiveGroupNames"
						},
						{
							"value": "LoginUser.Email",
							"text": "Email"
						},
						{
							"value": "LoginUser.FirstName",
							"text": "FirstName"
						},
						{
							"value": "LoginUser.GroupDNs",
							"text": "GroupDNs"
						},
						{
							"value": "LoginUser.GroupNames",
							"text": "GroupNames"
						},
						{
							"value": "LoginUser.GroupNames2",
							"text": "GroupNames2"
						},
						{
							"value": "LoginUser.HomeNumber",
							"text": "HomeNumber"
						},
						{
							"value": "LoginUser.LastName",
							"text": "LastName"
						},
						{
							"value": "LoginUser.MobileNumber",
							"text": "MobileNumber"
						},
						{
							"value": "LoginUser.OfficeNumber",
							"text": "OfficeNumber"
						},
						{
							"value": "LoginUser.RoleNames",
							"text": "RoleNames"
						},
						{
							"value": "LoginUser.Shortname",
							"text": "Shortname"
						},
						{
							"value": "LoginUser.Username",
							"text": "Username"
						},
						{
							"value": "LoginUser.Uuid",
							"text": "Uuid"
						}
					],
					"text": "LoginUser"
				}
			],
			"DisplayName": "SAML",
			"UseDefaultSigningCert": false,
			"SupportedAuthnContextClasses": [
				"unspecified",
				"PasswordProtectedTransport",
				"AuthenticatedTelephony",
				"InternetProtocol",
				"InternetProtocolPassword",
				"Kerberos",
				"MobileOneFactorContract",
				"MobileOneFactorUnregistered",
				"MobileTwoFactorContract",
				"MobileTwoFactorUnregistered",
				"NomadTelephony",
				"Password",
				"PersonalTelephony",
				"PGP",
				"PreviousSession",
				"SecureRemotePassword",
				"Smartcard",
				"SmartcardPKI",
				"SoftwarePKI",
				"SPKI",
				"Telephony",
				"TimeSyncToken",
				"TLSClient",
				"X509",
				"XMLDSig"
			],
			"ShowAccountIdentifierFinal": false,
			"SamlVersionLabel": "SAML Version",
			"SpMetadataUrl": "",
			"ShowIssuer": true,
			"SamlVersion": 2,
			"Issuer": "https://xxxxxxx.my.centrify.net/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"ShowInUP": true,
			"CatalogVisibility": "All",
			"_entitycontext": "W/\"datetime'2021-03-29T07%3A19%3A20.5448613Z'\"",
			"WantAssertionsSignedReadOnly": false,
			"_TableName": "application",
			"RecipientSameAsAcsUrl": true,
			"Generic": true,
			"LocalizationMappings": [...],
			"Audience": "",
			"UrlHelptip": "ACS URL is a given by your Service Provider. Enter it here.",
			"State": "Configured",
			"RecipientReadOnly": false,
			"SpMetadataXmlHelpTip": "If your Service Provider Metadata is given out in XML contents, copy and paste them here.",
			"RegistrationLinkMessage": null,
			"IssuerLabel": "IdP Entity ID / Issuer",
			"ShowCertBasedAuth": false,
			"RecipientShow": true,
			"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"HasAppPlugin": false,
			"Recipient": "",
			"SignOutUrlLabel": "Single Logout URL",
			"ProvScriptTemplateVersion": 1,
			"AudienceLabel": "SP Entity ID / Issuer / Audience",
			"IssuerHelptip": "This is your IdP Entity ID, also known as IdP Issuer. Give this to Service Provider during SAML configuration.",
			"SpMetadataUrlReadOnly": false,
			"SpMetadataUrlShow": true,
			"CertificateHelptip": "Download the certificate used to sign the SAML response to upload or configure in your SAML application.",
			"NameIDFormat": "unspecified",
			"_encryptkeyid": "ABC0000",
			"SpMetadataUrlLabel": "SP Metadata URL",
			"ShowIdentityProviderUrls": true,
			"RecipientSameAsAcsUrlReadOnly": false,
			"_PartitionKey": "ABC0000",
			"IdpConfigMethod": 1,
			"ShowAccountIdentifierInNewAppInstances": false,
			"LogoutUrl": "https://xxxxxx.my.centrify.net/applogout/appkey/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/customerid/ABC0000",
			"WantAssertionsSigned": false,
			"Reference": "CC-55325",
			"RelayStateHelpTip": "If your Service Provider specifies a Relay State value to use, enter it here.",
			"SetUseAttrInMetadata": true,
			"UrlLabel": "Assertion Consumer Service (ACS) URL",
			"SamlVersionHelpTip": "If you are configuring SAML 1.0 or 1.1, select SAML 1.x.",
			"ShowAdditionalField1Final": false,
			"SamlVersionShow": true,
			"CertificateSubjectName": "CN=Centrify Customer ABC0000 Application Signing Certificate",
			"ParentDisplayName": null,
			"_metadata": {
				"Version": 1,
				"IndexingVersion": 1
			},
			"Description": "This template enables you to provide single sign-on to a web application that uses SAML (Security Assertion Markup Language) for authentication.",
			"ErrorUrlLabel": "Single Sign On Error URL",
			"NameIDFormatReadOnly": false,
			"ErrorUrl": "https://xxxxxx.my.centrify.net/uperror?title=Error%20Signing%20In&message=Error%20encountered%20signing%20in%20to%20application%3A%20SAML&details=Error%20encountered%20signing%20in%20to%20application%3A%20SAML&customerid=ABC0000",
			"SignInUrl": "https://xxxxxx.my.centrify.net/applogin/appKey/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/customerId/ABC0000",
			"SpMetadataXmlReadOnly": false,
			"SpMetadataXmlLabel": "SP Metadata XML",
			"EncryptionCertOption": "Optional",
			"AuthRules": {
				"_UniqueKey": "Condition",
				"_SingleRow": true,
				"_Value": [],
				"Enabled": true,
				"_Type": "RowSet"
			},
			"RecipientLabel": "Recipient",
			"AppType": "Web",
			"AuthnContextClassReadOnly": false,
			"Name": "SAML",
			"Thumbprint": "F222EE02CCEA3C2D9243562F316E9C47A7A6B599",
			"UserNameStrategy": "ADAttribute", // or Fixed
			"CertInputMethod": "UploadCertDefault",
			"CertificateLabel": "Signing Certificate",
			"TemplateName": "Generic SAML",
			"SamlScript": "@Empty",
			"Handler": "cloudlib;Centrify.Saas.apphandlers.SamlAppHandler",
			"RelayStateReadOnly": false,
			"AudienceHelpTip": "SP Entity ID, also known as SP Issuer, or Audience, is a value given by your Service Provider. Enter it here.",
			"DefaultAuthProfile": "AlwaysAllowed",
			"NameIDFormatShow": true,
			"RecipientSameAsAcsUrlHelpTip": "If your Service Provider specifies a Recipient value to use, un-check this checkbox and enter your Recipient value below.",
			"SpMetadataXml": "",
			"Url": "",
			"ShowAdditionalField1InNewAppInstances": false,
			"SignInUrlLabel": "Single Sign On URL",
			"SignInUrlHelptip": "Configure this in your SAML application for SAML SP-initiated authentication.",
			"NameIDFormatLabel": "NameID Format",
			"RecipientSameAsAcsUrlLabel": "Same As ACS URL",
			"AppTypeDisplayName": "Web - SAML",
			"AuthChallengeDefinitionId": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
			"NameIDFormatHelpTip": "This is the Format attribute value in the &lt;NameID&gt; element in SAML Response. Select the NameID Format that your Service Provider specifies to use. If SP does not specify one, select &apos;unspecified&apos;.",
			"ErrorUrlHelptip": "Configure this in your SAML application if you want to redirect users to an error message from the User portal when SAML authentication fails.",
			"IdpMetadataXml": "xxxxxxxxxx",
			"ShowAccountIdentifier": false,
			"AllowUseOfDefaultSigningCertificate": true,
			"WantAssertionsSignedShow": true,
			"RelayStateShow": true,
			"RegistrationMessage": null,
			"NoQueryInIdpUrl": true,
			"_Timestamp": "/Date(1616230582295)/",
			"ProvCapable": true,
			"SpMetadataUrlHelpTip": "If your Service Provider Metadata is given out at a URL, enter it here and then click Load.",
			"UserNameArg": "userprincipalname",
			"AuthnContextClass": "unspecified",
			"AdminTag": "Other",
			"RelayStateLabel": "Relay State",
			"AuthnContextClassLabel": "Authentication Context Class",
			"ProvHandler": "Centrify.Cloud.Saas.Provisioning.Scim2.UserSync;Centrify.Cloud.Saas.Provisioning.Scim.UserSync.ScimUserSync",
			"SupportedNameIDFormats": [
				"unspecified",
				"emailAddress",
				"transient",
				"persistent",
				"entity",
				"kerberos",
				"WindowsDomainQualifiedName",
				"X509SubjectName"
			],
			"SupportsLinkedApps": true,
			"WebAppType": "Saml",
			"WantAssertionsSignedHelpTip": "Your Service Provider will specify which element to sign between Response and Assertion. If not specified or not sure, try Response first.",
			"AudienceReadOnly": false,
			"WantAssertionsSignedLabel": "Sign Response or Assertion?",
			"Category": "Other",
			"SpMetadataXmlShow": true,
			"UrlReadOnly": false,
			"ContentInspectionPolicyContentTypes": "ContentTypes=text/html;text/javascript;application/javascript;application/json;application/xml;",
			"OnPrem": true
		},
		"IsSoftError": false
	}

Create SAML web app
	Request body format
	{
		"ID": [
			"Generic SAML"
		]
	}

	Respond result
	{
		"success": true,
		"Result": [
			{
				"success": true,
				"ID": "Generic SAML",
				"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
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

Update SAML web app
	Request body format
	{
		"IdpConfigMethod": 1,
		"Issuer": "https://xxxxxx.my.centrify.net/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"Thumbprint": "F222EE02CCEA3C2D9243562F316E9C47A7A6B599",
		"SpConfigMethod": 1,
		"SpMetadataUrl": "https://nexus.microsoftonline-p.com/federationmetadata/saml20/federationmetadata.xml",
		"SpMetadataXml": "xxxxxxx",
		"Audience": "urn:federation:MicrosoftOnline",
		"Url": "https://login.microsoftonline.com/login.srf",
		"RecipientSameAsAcsUrl": true,
		"WantAssertionsSigned": true,
		"NameIDFormat": "unspecified",
		"SpSingleLogoutUrl": "",
		"EncryptAssertion": false,
		"EncryptionThumbprint": "",
		"RelayState": "",
		"AuthnContextClass": "unspecified",
		"SamlAttributes": [
			{
				"Name": "attribute1",
				"Value": "value1"
			}
		],
		"Script": "Email=\"afdadsf@afd.com\";",
		"AuthRules": null,
		"DefaultAuthProfile": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
		"PolicyScript": "\nif(!context.onPrem){\n    trace(\"Not onprem\");\n    var umod = module('User');\n    var user = umod.GetCurrentUser();\n    if(user.InRole(\"System Administrator\")){\n        trace(\"Allow System Administrator\");\n        policy.RequiredLevel = 2;\n    } else {\n        trace(\"Block non-System-Administrator\");\n        policy.Locked = true;\n    }\n}\n\n",
		"IconUri": "/vfslow/lib/application/icons/GenericSaml.svg",
		"_RowKey": "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
	}

	Respond result
	{
		"success": true,
		"Result": {
			"State": 0
		},
		"Message": null,
		"MessageID": null,
		"Exception": null,
		"ErrorID": null,
		"ErrorCode": null,
		"IsSoftError": false,
		"InnerExceptions": null
	}

*/
