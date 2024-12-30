package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/japannext/snooze/pkg/common/opensearch/dsl"
)

type Mappings struct {
	Properties map[string]Mapping `json:"properties"`
}

type Mapping struct {
	Type       string             `json:"type"`
	Fields     *Fields            `json:"fields"`
	Properties map[string]Mapping `json:"properties"`
}

type Fields struct {
	Keyword Keyword `json:"keyword"`
}

type Keyword struct {
	Type        string `json:"type"`
	IgnoreAbove int    `json:"ignore_above"`
}

// Do a search on all fields prefixed with the prefix.
// This is useful for "objects" that contain key/values that may change
// with future configurations, like "identity" or "labels"
func (req *SearchReq) MultiFieldSearch(ctx context.Context, prefixes []string, search string) error {
	request := opensearchapi.IndicesGetReq{
		Indices: []string{req.Index},
	}
	resp, err := client.Indices.Get(ctx, request)
	if err != nil {
		return err
	}
	idx, ok := resp.Indices[req.Index]
	if !ok {
		return fmt.Errorf("index %s not found in response", req.Index)
	}
	var mappings *Mappings
	if err := json.Unmarshal(idx.Mappings, &mappings); err != nil {
		return fmt.Errorf("failed to unmarshal index %s: %s", req.Index, err)
	}

	flats := flattenMappings(mappings.Properties)
	log.Debugf("flats: %s", flats)
	flats = filterMappings(flats, prefixes)

	req.Doc.Query.Bool.MinimumShouldMatch = 1
	for _, flat := range flats {
		req.Doc.Or(dsl.QueryItem{Match: map[string]dsl.Match{
			flat: {Query: search},
		}})
	}

	return nil
}

// Flatten the mapping structures of the index.
// e.g. labels.my-label1, labels.my-label2, ...
func flattenMappings(mappings map[string]Mapping) []string {
	var flats []string
	for prefix, mapping := range mappings {
		if len(mapping.Properties) > 0 {
			for _, key := range flattenMappings(mapping.Properties) {
				flats = append(flats, fmt.Sprintf("%s.%s", prefix, key))
			}
		}
	}
	return flats
}

// Filter flattened entries by starting prefix (e.g. identity, labels, ...)
func filterMappings(flats, prefixes []string) []string {
	var filtered []string
	m := make(map[string]bool)
	for _, prefix := range prefixes {
		m[prefix] = true
	}
	for _, flat := range flats {
		prefix := strings.Split(flat, ".")[0]
		if _, ok := m[prefix]; ok {
			filtered = append(filtered, flat)
		}
	}
	return filtered
}
