package applicationtemplate

// ApplicationTemplate is an enum of the various application template
type ApplicationTemplate int

const (
	SAML ApplicationTemplate = iota
	AWSConsole
	Cloudera
	CloudLock
	ConfluenceServer
	Dome9
	GitHubEnterprise
	JIRACloud
	JIRAServer
	PaloAltoNetworks
	SplunkOnPrem
	SumoLogic
)

func (r ApplicationTemplate) String() string {
	names := [...]string{
		"Generic SAML",
		"AWSConsoleSAML",
		"ClouderaSAML",
		"CloudLock SAML",
		"ConfluenceServerSAML",
		"Dome9Saml",
		"GitHubEnterpriseSAML",
		"JIRACloudSAML",
		"JIRAServerSAML",
		"PaloAltoNetworksSAML",
		"SplunkOnPremSAML",
		"SumoLogicSAML"}

	return names[r]
}
