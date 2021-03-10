package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ResponsesDiff is a diff between two sets of response objects: https://swagger.io/specification/#responses-object
type ResponsesDiff struct {
	Added    StringList        `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList        `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedResponses `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty return true if there is no diff
func (responsesDiff *ResponsesDiff) Empty() bool {
	if responsesDiff == nil {
		return true
	}

	return len(responsesDiff.Added) == 0 &&
		len(responsesDiff.Deleted) == 0 &&
		len(responsesDiff.Modified) == 0
}

// ModifiedResponses is map of response values to their respective diffs
type ModifiedResponses map[string]*ResponseDiff

func newResponsesDiff() *ResponsesDiff {
	return &ResponsesDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedResponses{},
	}
}

func getResponsesDiff(config *Config, responses1, responses2 openapi3.Responses) *ResponsesDiff {
	diff := getResponsesDiffInternal(config, responses1, responses2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getResponsesDiffInternal(config *Config, responses1, responses2 openapi3.Responses) *ResponsesDiff {

	result := newResponsesDiff()

	for responseValue1, responseRef1 := range responses1 {
		if responseRef1 != nil && responseRef1.Value != nil {
			if responseValue2, ok := responses2[responseValue1]; ok {
				if diff := diffResponseValues(config, responseRef1.Value, responseValue2.Value); !diff.Empty() {
					result.Modified[responseValue1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, responseValue1)
			}
		}
	}

	for responseValue2, responseRef2 := range responses2 {
		if responseRef2 != nil && responseRef2.Value != nil {
			if _, ok := responses1[responseValue2]; !ok {
				result.Added = append(result.Added, responseValue2)
			}
		}
	}

	return result
}

func (responsesDiff *ResponsesDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(responsesDiff.Added),
		Deleted:  len(responsesDiff.Deleted),
		Modified: len(responsesDiff.Modified),
	}
}
