package redisutil


// IndexFieldInfo holds information about a single indexed field
type IndexFieldInfo struct {
	Path   string `json:"path"`
	Alias  string `json:"alias"`
	Type   string `json:"type"`
}

// IndexInfo holds information about a RediSearch index and its fields
type IndexInfo struct {
	Name   string           `json:"name"`
	Fields []IndexFieldInfo `json:"fields"`
}

// GetIndexesAndFields lists RediSearch indexes and their fields from static knowledge (for healthz)
func GetIndexesAndFields() ([]IndexInfo, error) {
	// These are hardcoded to match the FT.CREATE statements in CreateCustomerIndex and CreateEventIndex
	return []IndexInfo{
		{
			Name: "customerIdx",
			Fields: []IndexFieldInfo{
				{Path: "$.primaryIdentifiers.email", Alias: "email", Type: "TEXT"},
				{Path: "$.primaryIdentifiers.phone", Alias: "phone", Type: "TEXT"},
				{Path: "$.primaryIdentifiers.visitor_id", Alias: "visitor_id", Type: "TEXT"},
			},
		},
		{
			Name: "eventIdx",
			Fields: []IndexFieldInfo{
				{Path: "$.identifiers.visitor_id", Alias: "visitor_id", Type: "TEXT"},
				{Path: "$.identifiers.call_id", Alias: "call_id", Type: "TEXT"},
				{Path: "$.identifiers.chat_id", Alias: "chat_id", Type: "TEXT"},
				{Path: "$.identifiers.external_id", Alias: "external_id", Type: "TEXT"},
				{Path: "$.identifiers.lead_id", Alias: "lead_id", Type: "TEXT"},
				{Path: "$.identifiers.tickets_id", Alias: "tickets_id", Type: "TEXT"},
			},
		},
	}, nil
}
