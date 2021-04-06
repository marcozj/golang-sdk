package applicationtemplate

// ApplicationTemplate is an enum of the various application template
type ApplicationTemplate int

const (
	Generic ApplicationTemplate = iota
)

func (r ApplicationTemplate) String() string {
	names := [...]string{
		"Generic OpenID Connect",
	}

	return names[r]
}
