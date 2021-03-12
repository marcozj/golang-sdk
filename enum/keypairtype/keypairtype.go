package keypairtype

// KeyPairType is an enum of the various type of SSHKey pair
type KeyPairType int

const (
	PublicKey KeyPairType = iota
	PrivateKey
	PuTTY
)

// String converts the KeyPairType to a string
func (keytype KeyPairType) String() string {
	switch keytype {
	case PublicKey:
		return "PublicKey"
	case PrivateKey:
		return "PrivateKey"
	case PuTTY:
		return "PPK"
	default:
		return "PrivateKey"
	}
}
