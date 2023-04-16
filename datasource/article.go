package datasource

type Article struct {
	Slug           string  `json:"slug"`
	Headline       string  `json:"headline"`
	ThumbnailURL   string  `json:"thumbnail_url"`
	Timestamp      int64   `json:"timestamp"`
	Embedding      []float32 `json:"embedding"`
}