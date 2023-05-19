package dbutils

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/mongo"
)

const baseURL = "https://www.thedp.com/section/"

var categories = []string{"news", "sports", "opinion"}

func UploadArticles(db *mongo.Database) error {  // rewrite this to use worker pools
	var wg sync.WaitGroup
	coll := db.Collection("articles")
	for _, category := range categories {
		wg.Add(1)
		go func(category string) {
			defer wg.Done()
			pageNumber := 1
			for {
				url := fmt.Sprintf("%s%s.json?page=%d", baseURL, category, pageNumber)
				fmt.Println(url)
				page, err := FetchPage(url)
				if err != nil {
					fmt.Printf("Error fetching page for category %s, page %d: %v\n", category, pageNumber, err)
					continue
				}
				if page.Pagination.Current == page.Pagination.Last {
					break
				}
				var articlesInterface []interface{}
				for i := range page.Articles {
					page.Articles[i].Content = cleanText(page.Articles[i].Content)
					page.Articles[i].Abstract = cleanText(page.Articles[i].Abstract)
					page.Articles[i].Infobox = cleanText(page.Articles[i].Infobox)
					page.Articles[i].Section = category
					cTime, _ := time.Parse("2006-01-02 15:04:05", page.Articles[i].CreatedAt)
					mTime, _ := time.Parse("2006-01-02 15:04:05", page.Articles[i].ModifiedAt)
					pTime, _ := time.Parse("2006-01-02 15:04:05", page.Articles[i].PublishedAt)
					page.Articles[i].CTime = cTime.Unix()
					page.Articles[i].MTime = mTime.Unix()
					page.Articles[i].PTime = pTime.Unix()
					articlesInterface = append(articlesInterface, page.Articles[i])
				}
				if _, err = coll.InsertMany(context.TODO(), articlesInterface); err != nil {
					fmt.Printf("Error uploading articles for category %s, page %d: %v\n", category, pageNumber, err)
				}
				pageNumber++
			}
		}(category)
	}
	wg.Wait()
	return nil
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
