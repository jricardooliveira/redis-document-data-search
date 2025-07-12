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
				{Path: "$.primaryIdentifiers.cmec_visitor_id", Alias: "visitor_id", Type: "TEXT"},
			},
		},
		{
			Name: "eventIdx",
			Fields: []IndexFieldInfo{
				{Path: "$.identifiers.cmec_visitor_id", Alias: "visitor_id", Type: "TEXT"},
				{Path: "$.identifiers.cmec_contact_call_id", Alias: "call_id", Type: "TEXT"},
				{Path: "$.identifiers.cmec_contact_chat_id", Alias: "chat_id", Type: "TEXT"},
				{Path: "$.identifiers.cmec_contact_external_id", Alias: "external_id", Type: "TEXT"},
				{Path: "$.identifiers.cmec_contact_form2lead_id", Alias: "form2lead_id", Type: "TEXT"},
				{Path: "$.identifiers.cmec_contact_tickets_id", Alias: "tickets_id", Type: "TEXT"},
			},
		},
	}, nil
}
