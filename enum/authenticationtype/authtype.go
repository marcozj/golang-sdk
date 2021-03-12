package authenticationtype

// AuthenticationType is an enum of the various type of Authentication type used authenticate against tenant
type AuthenticationType int

const (
	OAuth2 AuthenticationType = iota
	UsernamePassword
	DelegatedMachineCredential
)

func (r AuthenticationType) String() string {
	names := [...]string{
		"oauth",
		"unpw",
		"dmc"}

	return names[r]
}
