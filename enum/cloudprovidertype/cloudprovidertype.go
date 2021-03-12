package cloudprovidertype

// CloudProviderType is an enum of the various type of cloud providers
type CloudProviderType int

const (
	AWS CloudProviderType = iota
)

func (r CloudProviderType) String() string {
	names := [...]string{
		"Aws"}

	return names[r]
}
