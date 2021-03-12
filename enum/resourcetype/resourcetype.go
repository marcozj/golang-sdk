package resourcetype

// ResourceType is an enum of the various type of ResourceType
type ResourceType int

const (
	System ResourceType = iota
	Database
	Domain
	CloudProvider
)

func (r ResourceType) String() string {
	names := [...]string{
		"system",
		"database",
		"domain",
		"cloudprovider"}

	return names[r]
}
