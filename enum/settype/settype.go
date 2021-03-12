package settype

// SetType is an enum of the various type of manual set
type SetType int

const (
	System SetType = iota
	Database
	Domain
	Account
	Secret
	SSHKey
	Service
	Application
	ResourceProfile
	CloudProvider
)

func (r SetType) String() string {
	names := [...]string{
		"Server",
		"VaultDatabase",
		"VaultDomain",
		"VaultAccount",
		"DataVault",
		"SshKeys",
		"Subscriptions",
		"Application",
		"ResourceProfiles",
		"CloudProviders"}

	return names[r]
}
