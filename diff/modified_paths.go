package diff

import (
	"github.com/rjmohammad/kin-openapi/openapi3"
)

// ModifiedPaths is a map of paths to their respective diffs
type ModifiedPaths map[string]*PathDiff

func (modifiedPaths ModifiedPaths) addPathDiff(config *Config, state *state, path1 string, pathItem1, pathItem2 *openapi3.PathItem) error {

	diff, err := getPathDiff(config, state, pathItem1, pathItem2)
	if err != nil {
		return err
	}

	if !diff.Empty() {
		modifiedPaths[path1] = diff
	}

	return nil
}
