package models

type SearchResult struct {
	Type        string `json:"type"`
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
	Metadata    any    `json:"metadata,omitempty"`
}

type SearchIndexItem struct {
	ID          int64  `json:"id" db:"id"`
	Type        string `json:"type" db:"type"`
	Key         string `json:"key" db:"key"`
	DisplayName string `json:"displayName" db:"display_name"`
	Metadata    any    `json:"metadata" db:"metadata"`
}
