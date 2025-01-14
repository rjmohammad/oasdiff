package diff

import "github.com/rjmohammad/kin-openapi/openapi3"

type schemaPair struct {
	Schema1 *openapi3.SchemaRef
	Schema2 *openapi3.SchemaRef
}

type schemaDiffCache map[schemaPair]*SchemaDiff
