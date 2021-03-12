package databaseclass

// DatabaseClass is an enum of the various type of database
type DatabaseClass int

const (
	SQLServer DatabaseClass = iota
	Oracle
	SAPASE
)

func (r DatabaseClass) String() string {
	names := [...]string{
		"SQLServer",
		"Oracle",
		"SAPAse"}

	return names[r]
}
