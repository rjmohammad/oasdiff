package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ContentDiff is the diff between two content objects each containing a media type object: https://swagger.io/specification/#media-type-object
type ContentDiff struct {
	MediaTypeAdded   bool            `json:"mediaTypeAdded,omitempty" yaml:"mediaTypeAdded,omitempty"`
	MediaTypeDeleted bool            `json:"mediaTypeDeleted,omitempty" yaml:"mediaTypeDeleted,omitempty"`
	MediaTypeDiff    bool            `json:"mediaType,omitempty" yaml:"mediaType,omitempty"`
	ExtensionsDiff   *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	SchemaDiff       *SchemaDiff     `json:"schema,omitempty" yaml:"schema,omitempty"`
	ExampleDiff      *ValueDiff      `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples
	EncodingsDiff *EncodingsDiff `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

func newContentDiff() *ContentDiff {
	return &ContentDiff{}
}

// Empty return true if there is no diff
func (contentDiff *ContentDiff) Empty() bool {
	return contentDiff == nil || *contentDiff == ContentDiff{}
}

func getContentDiff(config *Config, content1, content2 openapi3.Content) *ContentDiff {
	diff := getContentDiffInternal(config, content1, content2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getContentDiffInternal(config *Config, content1, content2 openapi3.Content) *ContentDiff {

	if len(content1) == 0 && len(content2) == 0 {
		return nil
	}

	result := newContentDiff()

	if len(content1) == 0 && len(content2) != 0 {
		result.MediaTypeAdded = true
		return result
	}

	if len(content1) != 0 && len(content2) == 0 {
		result.MediaTypeDeleted = true
		return result
	}

	mediaType1, mediaTypeValue1, err := getMediaType(content1)
	if err != nil {
		return nil
	}

	mediaType2, mediaTypeValue2, err := getMediaType(content2)
	if err != nil {
		return nil
	}

	if mediaType1 != mediaType2 {
		result.MediaTypeDiff = true
		return result
	}

	result.ExtensionsDiff = getExtensionsDiff(config, mediaTypeValue1.ExtensionProps, mediaTypeValue2.ExtensionProps)
	result.SchemaDiff = getSchemaDiff(config, mediaTypeValue1.Schema, mediaTypeValue2.Schema)

	if config.IncludeExamples {
		result.ExampleDiff = getValueDiff(mediaTypeValue1.Example, mediaTypeValue1.Example)
	}

	result.EncodingsDiff = getEncodingsDiff(config, mediaTypeValue1.Encoding, mediaTypeValue2.Encoding)

	return result
}

// getMediaType returns the single MediaType entry in the content map
func getMediaType(content openapi3.Content) (string, *openapi3.MediaType, error) {

	var mediaType string
	var mediaTypeValue *openapi3.MediaType

	if len(content) != 1 {
		return "", nil, fmt.Errorf("content map has more than one value - this shouldn't happen. %+v", content)
	}

	for k, v := range content {
		mediaType = k
		mediaTypeValue = v
	}

	return mediaType, mediaTypeValue, nil
}
