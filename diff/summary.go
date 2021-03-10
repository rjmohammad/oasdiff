package diff

import "reflect"

// Summary summarizes the changes between two OpenAPI specifications
type Summary struct {
	Diff       bool                              `json:"diff" yaml:"diff,omitempty"`
	Components map[ComponentName]*SummaryDetails `json:"components,omitempty" yaml:"components,omitempty"`
}

func newSummary() *Summary {
	return &Summary{
		Components: map[ComponentName]*SummaryDetails{},
	}
}

// SummaryDetails summarizes the changes between equivalent parts of the two OpenAPI specifications: paths, schemas, parameters, headers, responses etc.
type SummaryDetails struct {
	Added    int `json:"added,omitempty" yaml:"added,omitempty"`    // number of added items
	Deleted  int `json:"deleted,omitempty" yaml:"deleted,omitempty"`  // number of deleted items
	Modified int `json:"modified,omitempty" yaml:"modified,omitempty"` // number of modified items
}

type componentWithSummary interface {
	getSummary() *SummaryDetails
}

// ComponentName is the key type of the summary map
type ComponentName string

// Component constants are the keys in the summary map
const (
	PathsComponent           ComponentName = "paths"
	SecurityComponent        ComponentName = "security"
	ServersComponent         ComponentName = "servers"
	TagsComponent            ComponentName = "tags"
	SchemasComponent         ComponentName = "schemas"
	ParametersComponent      ComponentName = "parameters"
	HeadersComponent         ComponentName = "headers"
	RequestBodiesComponent   ComponentName = "requestBodies"
	ResponsesComponent       ComponentName = "responses"
	SecuritySchemesComponent ComponentName = "securitySchemes"
	CallbacksComponent       ComponentName = "callbacks"
)

// GetSummaryDetails returns the summary for a specific component
func (summary *Summary) GetSummaryDetails(component ComponentName) SummaryDetails {
	if details, ok := summary.Components[component]; ok {
		if details != nil {
			return *details
		}
	}

	return SummaryDetails{}
}

func (summary *Summary) add(component componentWithSummary, name ComponentName) {
	if !isNilPointer(component) {
		summary.Components[name] = component.getSummary()
	}
}

func isNilPointer(i interface{}) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()
}
