package server

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type VectorSearch struct {
	ctx context.Context
	rdb *redis.Client
}

func (vs *VectorSearch) initializeSchema() {
	vs.rdb.Do(vs.ctx,
		"FT.CREATE", "articles", "ON", "JSON", "PREFIX", "1", "article:", "SCHEMA",
		"$.headline", "AS", "headline", "TEXT",
		"$.thumbnail_url", "AS", "thumbnail_url", "TEXT",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
		"$.embedding", "AS", "embedding", "VECTOR", "HNSW", "6", "TYPE", "FLOAT32", "DIM", "50", "DISTANCE_METRIC", "COSINE",
		// need to store vector embedding, headline, thumbnail_url (if exists), and timestamp
		// slug is the id
		// thumbnail url is stored in article["dominantMedia"]["attachment_uuid"]
		// ex. https://snworksceo.imgix.net/dpn/f799a2b0-b72f-4a0b-a927-659bc30dff92.sized-1000x1000.png?w=800
	)

	// pipe := vs.rdb.Pipeline()
}

func (vs *VectorSearch) search(slug string) {
	/* check if article is in Redis, else pull article */
	res, err := vs.rdb.Do(vs.ctx, "JSON.GET", "article:"+slug).Slice()
	if err != nil {
		return
	}
	print(res)

}