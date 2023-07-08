package dbutils

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/mongo"
)

const baseURL = "https://www.thedp.com/section/"
var categories = []string{"news", "sports", "opinion"}

func worker(id int, jobs <- chan string, db *mongo.Database) {
	coll := db.Collection("articles")
	for j := range jobs {
		fmt.Printf("worker %d: processing url %s\n", id, j)
		page, err := FetchPage(j)
		if err != nil {
			fmt.Printf("worker %d: error fetching %s\n", id, j)
			continue
		}
		var articlesInterface []interface{}
		for i := range page.Articles {
			page.Articles[i].Content = cleanText(page.Articles[i].Content)
			page.Articles[i].Abstract = cleanText(page.Articles[i].Abstract)
			page.Articles[i].Infobox = cleanText(page.Articles[i].Infobox)
			cTime, _ := time.Parse("2006-01-02 15:04:05", page.Articles[i].CreatedAt)
			mTime, _ := time.Parse("2006-01-02 15:04:05", page.Articles[i].ModifiedAt)
			pTime, _ := time.Parse("2006-01-02 15:04:05", page.Articles[i].PublishedAt)
			page.Articles[i].CTime = cTime.Unix()
			page.Articles[i].MTime = mTime.Unix()
			page.Articles[i].PTime = pTime.Unix()
			articlesInterface = append(articlesInterface, page.Articles[i])
		}
		if _, err = coll.InsertMany(context.TODO(), articlesInterface); err != nil {
			fmt.Printf("worker %d: error uploading articles for url %s\n", id, j)
		}
	}
}

func UploadArticles(db *mongo.Database) {
	start := time.Now()
	jobs := make(chan string)
	for w := 1; w <= 5; w++ {
        go worker(w, jobs, db)
    } 

	for _, category := range categories {
		url := fmt.Sprintf("%s%s.json", baseURL, category)
		page, _ := FetchPage(url)
		for i := 1; i <= page.Pagination.Last; i++ {
			jobs <- fmt.Sprintf("%s%s.json?page=%d", baseURL, category, i)
		}
	}

	close(jobs)
	fmt.Println(time.Since(start))
}

func FetchPage(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed sending request to %s: %v", url, err)
	}
	defer resp.Body.Close()
	var page Response
	if err = json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("error parsing response body as JSON: %v", err)
	}
	return &page, nil
}

func cleanText(article string) string {
	strictPolicy := bluemonday.StrictPolicy()
	return html.UnescapeString(strictPolicy.Sanitize(article))
}
