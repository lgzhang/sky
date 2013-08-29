package query

const (
	TypeAssignment   = "assignment"
	TypeCondition    = "condition"
	TypeExit         = "exit"
	TypeSelection    = "selection"
	TypeTemporalLoop = "temporal_loop"
	TypeEventLoop    = "event_loop"
)

type Statement interface {
	QueryElement
	FunctionName(init bool) string
	MergeFunctionName() string
	Serialize() map[string]interface{}
	Deserialize(map[string]interface{}) error
	CodegenAggregateFunction(init bool) (string, error)
	CodegenMergeFunction(map[string]interface{}) (string, error)
	Defactorize(data interface{}) error
	RequiresInitialization() bool
	Variables() []*Variable
	String() string
}