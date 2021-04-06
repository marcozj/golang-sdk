package clientidtype

// ClientIDType is an enum of the various application template
type ClientIDType int

const (
	AnythingOrList ClientIDType = iota
	Confidential
)

func (r ClientIDType) String() string {
	names := [...]string{
		"list",
		"confidential",
	}

	return names[r]
}
