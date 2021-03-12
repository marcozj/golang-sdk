package applicationtemplate

// ApplicationTemplate is an enum of the various application template
type ApplicationTemplate int

const (
	Generic ApplicationTemplate = iota
	SQLServerManagementStudio
	Toad
	VSphereClient
)

func (r ApplicationTemplate) String() string {
	names := [...]string{
		"GenericDesktopApplication",
		"Ssms",
		"Toad",
		"VpxClient"}

	return names[r]
}
