package logincredential

// HostLoginCredentials is an enum of the various application template
type HostLoginCredentials int

const (
	UserADCredential HostLoginCredentials = iota
	PromptForCredential
	SelectAlternativeAccount
	SharedAccount
)

func (r HostLoginCredentials) String() string {
	names := [...]string{
		"ADCredential",
		"SetByUser",
		"AlternativeAccount",
		"SharedAccount"}

	return names[r]
}
