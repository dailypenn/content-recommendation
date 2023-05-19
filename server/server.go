package server

// import (
// 	"context"
// 	"github.com/redis/go-redis/v9"
// )

// func (vs *VectorSearch) search(slug string) {
// 	/* check if article is in Redis, else pull article */
// 	res, err := vs.rdb.Do(vs.ctx, "JSON.GET", "article:"+slug).Slice()
// 	if err != nil {
// 		return
// 	}
// 	print(res)

// }