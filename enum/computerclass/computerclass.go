package computerclass

// ComputerClass is an enum of the various type of resource for an vaulted account
type ComputerClass int

const (
	Windows ComputerClass = iota
	Unix
	CiscoIOS
	CiscoNXOS
	JuniperJunos
	HPNonStop
	IBMi
	CheckPointGaia
	PaloAltoPANOS
	F5BIGIP
	CiscoAsyncOS
	VMwareVMkernel
	GenericSSH
	CustomSSH
)

func (r ComputerClass) String() string {
	names := [...]string{
		"Windows",
		"Unix",
		"CiscoIOS",
		"CiscoNXOS",
		"JuniperJunos",
		"HpNonStopOS",
		"IBMi",
		"CheckPointGaia",
		"PaloAltoNetworksPANOS",
		"F5NetworksBIGIP",
		"CiscoAsyncOS",
		"VMwareVMkernel",
		"GenericSsh",
		"CustomSsh"}

	return names[r]
}
