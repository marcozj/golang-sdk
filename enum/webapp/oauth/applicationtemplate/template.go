package applicationtemplate

// ApplicationTemplate is an enum of the various application template
type ApplicationTemplate int

const (
	OAuth2Client ApplicationTemplate = iota
	OAuth2Server
)

func (r ApplicationTemplate) String() string {
	names := [...]string{
		"OAuth2ServerClient",
		"OAuth2Server",
	}

	return names[r]
}
