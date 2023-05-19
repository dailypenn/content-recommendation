package dbutils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateRedisConn() *VectorSearch {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	addr := os.Getenv("REDIS_ADDR")
	pw := os.Getenv("REDIS_PASSWORD")
	if addr == "" || pw == "" {
		log.Fatal("Missing environment variable(s) REDIS_ADDR and/or REDIS_PASSWORD")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
	})
	return &VectorSearch{
		ctx: ctx,
		rdb: client,
	}
}

func (vs *VectorSearch) InitializeSchema() {
	vs.rdb.Do(vs.ctx,
		"FT.CREATE", "articles", "ON", "JSON", "PREFIX", "1", "article:", "SCHEMA",
		"$.headline", "AS", "headline", "TEXT",
		"$.thumbnail_url", "AS", "thumbnail_url", "TEXT",
		"$.embedding", "AS", "embedding", "VECTOR", "HNSW", "6", "TYPE", "FLOAT32", "DIM", "50", "DISTANCE_METRIC", "COSINE",
	)
}

func (vs *VectorSearch) RefreshRecentArticles(db *mongo.Database) error {
	arr := make([]float64, 50)
	for i := range arr {
		arr[i] = 0.5
	}
	coll := db.Collection("articles")
	pipe := vs.rdb.Pipeline()
	pipe.FlushDB(vs.ctx)
	opts := options.Find().SetSort(map[string]int{"ctime": -1}).SetProjection(bson.D{{"slug", 1}, {"headline", 1}, {"dominantmedia", 1}, {"content", 1}, {"createdat", 1}}).SetLimit(1000)
	cursor, err := coll.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return fmt.Errorf("failed to execute mongo query: %v", err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("failed to unmarshal document: %v", err)
		}
		if _, isM := doc["dominantmedia"].(bson.M); isM {
			attachmentUUID := doc["dominantmedia"].(bson.M)["attachment_uuid"]
			fileExt := doc["dominantmedia"].(bson.M)["extension"]
			doc["thumbnail_url"] = fmt.Sprintf("https://snworksceo.imgix.net/dpn/%s.sized-1000x1000.%s?w=800", attachmentUUID, fileExt)
		}
		doc["slug"] = fmt.Sprintf("%s/%s/%s", doc["createdat"].(string)[:4], doc["createdat"].(string)[5:7], doc["slug"])
		doc["embedding"] = arr
		var d Document
		bb, _ := bson.Marshal(doc)
		bson.Unmarshal(bb, &d)
		jsonBytes, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("failed to marshal document to json: %v", err)
		}
		// add embedding, remove content
		pipe.Do(vs.ctx, "JSON.SET", "article:"+doc["slug"].(string), "$", jsonBytes)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	if _, err := pipe.Exec(vs.ctx); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (vs *VectorSearch) query(slug string, db *mongo.Database) {
	// TODO
}
