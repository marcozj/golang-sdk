package accountmapping

// AccountMapping is an enum of the various application template
type AccountMapping int

const (
	ADAttribute AccountMapping = iota
	SharedAccount
	UseScript
	SetByUser
)

func (r AccountMapping) String() string {
	names := [...]string{
		"ADAttribute",
		"Fixed",
		"UseScript",
		"SetByUser"}

	return names[r]
}
