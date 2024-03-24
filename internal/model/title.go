package model

type Title struct {
	Id          int64  `json:"id"`
	ImdbId      string `json:"imdb_id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Year        int64  `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     int64  `json:"runtime"`
}
