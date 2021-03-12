package cmdparamtype

// CmdParmType is an enum of the various application template
type CmdParmType int

const (
	Integer CmdParmType = iota
	Date
	String
	Account
	CloudProivder
	Database
	Device
	Domain
	ResourceProfile
	Role
	Secret
	Service
	SSHKey
	System
	User
)

func (r CmdParmType) String() string {
	names := [...]string{
		"int",
		"date",
		"string",
		"VaultAccount",
		"CloudProviders",
		"VaultDatabase",
		"Device",
		"VaultDomain",
		"ResourceProfile",
		"Role",
		"DataVault",
		"Subscriptions",
		"SshKeys",
		"Server",
		"User"}

	return names[r]
}
