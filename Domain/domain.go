package domain

import "time"

const (
	ADMIN       = "ADMIN"
	UESR        = "USER"
	SUPER_ADMIN = "SUPER_ADMIN"
)

type User struct {
	UserID         string
	UserName       string
	PersonalBio    string
	ProfilePic     string
	Email          string
	PhoneNum       string
	TelegramHandle string
	Password       string
	Role           string
}

type Blog struct {
	BlogID      string
	OwnerID     string
	Title       string
	Tags        []string
	Author      string
	Description string
	LastUpdate  time.Time
}

type Popularity struct {
	PopularityID string
	BlogID       string
	ViewCount    int
	Likes        []string
	Dislikes     []string
	Comments     []string
}
