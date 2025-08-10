package domain

import (
	"time"
)

const (
	ADMIN       = "ADMIN"
	UESR        = "USER"
	SUPER_ADMIN = "SUPER_ADMIN"
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
type RefreshToken struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}
type UserUnverified struct {
	UserName  string
	Email     string
	OTP       string
	Password  string
	Role      string
	ExpiresAt time.Time
}

type User struct {
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
	Title       string
	Tags        []string
	Author      string
	Description string
	LastUpdate  time.Time
}

type Popularity struct {
	PopularityValue int //? Lately added for the sake of calculating the popularity
	ViewCount       int
	Likes           []string
	Dislikes        []string
	Comments        []*Comment
}

type Comment struct {
	UserName string
	Comment  string
}

// I added the filter struct, b/c while filtering I was passing around 4 parameters at once so now 5 of them are in one struct it will be easy to pass arguments
type Filter struct {
	Popularity_value int // it should be either ascending or descending
	Tag              string
	AfterDate        *time.Time
	AuthorName       string
	Title            string
}
