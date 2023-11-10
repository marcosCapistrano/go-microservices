package model

import model "stocksapp/metadata/pkg/model"

// StockDetails includes movie metadata its aggregated
// rating.
type MovieDetails struct {
	Rating   *float64       `json:"rating"`
	Metadata model.Metadata `json:"metadata"`
}
