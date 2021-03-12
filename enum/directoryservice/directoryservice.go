package directoryservice

// DirectoryService is an enum of the various directory service
type DirectoryService int

const (
	CentrifyDirectory DirectoryService = iota
	ActiveDirectory
	FederatedDirectory
	GoogleDirectory
	LDAPDirectory
)

func (r DirectoryService) String() string {
	names := [...]string{
		"Centrify Directory",
		"Active Directory",
		"Federated Directory",
		"Google Directory",
		"LDAP Directory"}

	return names[r]
}
