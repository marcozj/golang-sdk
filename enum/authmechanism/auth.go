package authmechanism

// AuthenticationMechanism is an enum of the various type of Authentication Mechanisms used for authentication proflie
type AuthenticationMechanism int

const (
	Password AuthenticationMechanism = iota
	MobileAuthenticator
	PhoneCall
	SMS
	EmailConfirmationCode
	OATH_OTP
	Radius
	FIDO2
	SecurityQuestions
)

func (r AuthenticationMechanism) String() string {
	names := [...]string{
		"UP",
		"OTP",
		"PF",
		"SMS",
		"EMAIL",
		"OATH",
		"RADIUS",
		"U2F",
		"SQ"}

	return names[r]
}
