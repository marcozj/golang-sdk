package secrettype

// SecretType is an enum of the various type of secret
type SecretType int

const (
	Text SecretType = iota
	File
)

func (r SecretType) String() string {
	names := [...]string{
		"Text",
		"File"}

	return names[r]
}
