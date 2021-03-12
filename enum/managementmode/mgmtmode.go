package managementmode

// ManagementMode is an enum of the various type of Windows management mode
type ManagementMode int

const (
	Unknown ManagementMode = iota
	RPCOverTCP
	SMB
	WinRMOverHTTP
	WinRMOverHTTPS
	Disabled
)

func (r ManagementMode) String() string {
	names := [...]string{
		"Unknown",
		"RPCOverTCP",
		"Smb",
		"WinRMOverHttp",
		"WinRMOverHttps",
		"Disabled"}

	return names[r]
}
