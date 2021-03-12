package enum

// Weekday defines weekday data type
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (day Weekday) String() string {
	names := [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday"}

	if day < Sunday || day > Saturday {
		return "Unknown"
	}
	// return the name of a Weekday constant from the names array above.
	return names[day]
}
