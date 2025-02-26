package diff

import (
	"github.com/rjmohammad/kin-openapi/openapi3"
)

// ContentDiff describes the changes between content properties each containing media type objects: https://swagger.io/specification/#media-type-object
type ContentDiff struct {
	MediaTypeAdded    StringList         `json:"mediaTypeAdded,omitempty" yaml:"mediaTypeAdded,omitempty"`
	MediaTypeDeleted  StringList         `json:"mediaTypeDeleted,omitempty" yaml:"mediaTypeDeleted,omitempty"`
	MediaTypeModified ModifiedMediaTypes `json:"mediaTypeModified,omitempty" yaml:"mediaTypeModified,omitempty"`
}

// ModifiedMediaTypes is map of media type names to their respective diffs
type ModifiedMediaTypes map[string]*MediaTypeDiff

func newContentDiff() *ContentDiff {
	return &ContentDiff{
		MediaTypeAdded:    StringList{},
		MediaTypeDeleted:  StringList{},
		MediaTypeModified: ModifiedMediaTypes{},
	}
}

// Empty indicates whether a change was found in this element
func (diff *ContentDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.MediaTypeAdded) == 0 &&
		len(diff.MediaTypeDeleted) == 0 &&
		len(diff.MediaTypeModified) == 0
}

func (diff *ContentDiff) removeNonBreaking(state *state) {
	if diff.Empty() {
		return
	}

	switch state.direction {
	case directionRequest:
		// if this is part of the request, then media-types can be deleted without breaking the client
		diff.MediaTypeDeleted = nil
	case directionResponse:
		// if this is part of the response, then media-types can be added without breaking the client
		diff.MediaTypeAdded = nil
	}
}

func getContentDiff(config *Config, state *state, content1, content2 openapi3.Content) (*ContentDiff, error) {
	diff, err := getContentDiffInternal(config, state, content1, content2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking(state)
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getContentDiffInternal(config *Config, state *state, content1, content2 openapi3.Content) (*ContentDiff, error) {

	result := newContentDiff()

	for name1, media1 := range content1 {
		if media2, ok := content2[name1]; ok {
			diff, err := getMediaTypeDiff(config, state, media1, media2)
			if err != nil {
				return nil, err
			}

			if !diff.Empty() {
				result.MediaTypeModified[name1] = diff
			}
		} else {
			result.MediaTypeDeleted = append(result.MediaTypeDeleted, name1)
		}
	}

	for name2 := range content2 {
		if _, ok := content1[name2]; !ok {
			result.MediaTypeAdded = append(result.MediaTypeAdded, name2)
		}
	}

	return result, nil
}
