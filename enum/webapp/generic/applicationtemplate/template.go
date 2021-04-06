package applicationtemplate

// ApplicationTemplate is an enum of the various application template
type ApplicationTemplate int

const (
	Bookmark ApplicationTemplate = iota
	BrowserExtension
	BrowserExtensionAdvanced
	NTLMBasic
	UserPassword
)

func (r ApplicationTemplate) String() string {
	names := [...]string{
		"Generic Bookmark",
		"Generic Browser Extension",
		"GenericBrowserExtensionScript",
		"GenericNTLMBasic",
		"Generic User-Password",
	}

	return names[r]
}
