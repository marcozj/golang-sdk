package setsubtype

// SetSubtype is an enum of the various sub type of application
type SetSubtype int

const (
	Web SetSubtype = iota
	Desktop
)

func (r SetSubtype) String() string {
	names := [...]string{
		"Web",
		"Desktop"}

	return names[r]
}
