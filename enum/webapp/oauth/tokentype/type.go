package tokentype

// TokenType is an enum of the various application template
type TokenType int

const (
	JwtRS256 TokenType = iota
	Opaque
)

func (r TokenType) String() string {
	names := [...]string{
		"JwtRS256",
		"Opaque",
	}

	return names[r]
}
