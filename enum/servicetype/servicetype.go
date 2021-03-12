package servicetype

// WindowsServiceType is an enum of the various type of Windows Service
type WindowsServiceType int

const (
	WindowsService WindowsServiceType = iota
	ScheduledTask
	IISApplicationPool
)

func (r WindowsServiceType) String() string {
	names := [...]string{
		"WindowsService",
		"ScheduledTask",
		"IISApplicationPool"}

	return names[r]
}
