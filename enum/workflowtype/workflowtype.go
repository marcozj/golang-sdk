package workflowtype

// WorkflowType is an enum of the various type of manual set
type WorkflowType int

const (
	AccountWorkflow WorkflowType = iota
	AgentAuthWorkflow
	SecretsWorkflow
	PrivilegeElevationWorkflow
)

func (r WorkflowType) String() string {
	names := [...]string{
		"wf",
		"agentAuthWorkflow",
		"secretsWorkflow",
		"privilegeElevationWorkflow"}

	return names[r]
}
