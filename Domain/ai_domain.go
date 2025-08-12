package domain

import (
	"encoding/json"
	"time"
)

type AIResponse struct {
	MainResponse      json.RawMessage `json:"main_response" bson:"main_response"`
	EditorialResponse string          `json:"editorial_response" bson:"editorial_response"`
	IsNilResponse     bool            `json:"is_nil_response" bson:"is_nil_response"`
}

type AIBlogFilter struct {
	Tags       []string   `json:"tags" bson:"tags"`
	AfterDate  *time.Time `json:"after_date" bson:"after_date"`
	Title      string     `json:"title" bson:"title"`
	AuthorName string     `json:"author_name" bson:"author_name"`
}

