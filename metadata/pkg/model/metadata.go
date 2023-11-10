package model

// Metadata defines the metadata for stocks.
type Metadata struct {
	ID          string `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Industry    string `json:"industry"`
	Description string `json:"description"`
}
