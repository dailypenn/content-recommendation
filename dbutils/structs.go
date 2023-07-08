package dbutils

type Response struct {
	Section    Section    `json:"section"`
	Articles   []Article  `json:"articles"`
	Pagination Pagination `json:"pagination"`
}
type Tag struct {
	ID       string `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Metadata any    `json:"metadata"`
	CeoID    string `json:"ceo_id"`
}
type Section struct {
	ID          string `json:"id"`
	UUID        string `json:"uuid"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	SortOrder   any    `json:"sort_order"`
	Template    string `json:"template"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
	PublishedAt any    `json:"published_at"`
	Metadata    any    `json:"metadata"`
	CeoID       string `json:"ceo_id"`
	Tags        []Tag `json:"tags"`
}
type Author struct {
	ID       string `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Slug     string `json:"slug"`
	Bio      string `json:"bio"`
	Tagline  string `json:"tagline"`
	Metadata any    `json:"metadata"`
	CeoID    string `json:"ceo_id"`
	Status   string `json:"status"`
}
type DominantMedia struct {
	ID               string `json:"id"`
	UUID             string `json:"uuid"`
	AttachmentUUID   string `json:"attachment_uuid"`
	BaseName         string `json:"base_name"`
	SeoTitle         any    `json:"seo_title"`
	SeoDescription   any    `json:"seo_description"`
	SeoImage         any    `json:"seo_image"`
	Extension        string `json:"extension"`
	PreviewExtension string `json:"preview_extension"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	Source           any    `json:"source"`
	ClickThrough     any    `json:"click_through"`
	Type             string `json:"type"`
	Height           string `json:"height"`
	Width            string `json:"width"`
	Status           string `json:"status"`
	Weight           string `json:"weight"`
	CreatedAt        string `json:"created_at"`
	ModifiedAt       string `json:"modified_at"`
	PublishedAt      string `json:"published_at"`
	Metadata         []any  `json:"metadata"`
	Hits             string `json:"hits"`
	NormalizedTags   string `json:"normalized_tags"`
	SvgPreview       any    `json:"svg_preview"`
	CeoID            string `json:"ceo_id"`
	SstsID           any    `json:"ssts_id"`
	SstsPath         any    `json:"ssts_path"`
	Transcoded       string `json:"transcoded"`
	Authors          []any  `json:"authors"`
}
type Article struct {
	ID             string        `json:"id" bson:"_id"`
	UUID           string        `json:"uuid"`
	Slug           string        `json:"slug"`
	SeoTitle       any           `json:"seo_title"`
	SeoDescription any           `json:"seo_description"`
	SeoImage       any           `json:"seo_image"`
	Headline       string        `json:"headline"`
	Subhead        any           `json:"subhead"`
	Abstract       string        `json:"abstract"`
	Content        string        `json:"content"`
	Infobox        string        `json:"infobox"`
	Template       any           `json:"template"`
	ShortToken     string        `json:"short_token"`
	Status         string        `json:"status"`
	Weight         string        `json:"weight"`
	MediaID        string        `json:"media_id"`
	CreatedAt      string        `json:"created_at"`
	ModifiedAt     string        `json:"modified_at"`
	PublishedAt    string        `json:"published_at"`
	Metadata       []any         `json:"metadata"`
	Hits           string        `json:"hits"`
	NormalizedTags string        `json:"normalized_tags"`
	CeoID          string        `json:"ceo_id"`
	SstsID         any           `json:"ssts_id"`
	SstsPath       any           `json:"ssts_path"`
	Tags           []Tag         `json:"tags"`
	Authors        []Author      `json:"authors"`
	DominantMedia  any           `json:"dominantMedia"`
	// Section        string
	CTime		   int64
	MTime		   int64
	PTime		   int64

}
type Pagination struct {
	First    int `json:"first"`
	Last     int `json:"last"`
	Previous int `json:"previous"`
	Next     int `json:"next"`
	Total    int `json:"total"`
	Current  int `json:"current"`
}
type Document struct {
	Slug           string        `json:"slug"`
	Headline       string        `json:"headline"`
	Content        string        `json:"content"`
	DominantMedia  any           `json:"-"`
	CreatedAt	   string		 `json:"-"`
	Thumbnail_url  string		 `json:"thumbnail_url"`
	Embedding	   []float64     `json:"embedding"`
}