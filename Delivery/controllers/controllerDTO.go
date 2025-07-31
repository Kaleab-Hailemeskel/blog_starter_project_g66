package controllers

import (
	domain "blog_starter_project_g66/Domain"
	"time"
)

const (
	ADMIN       = "ADMIN"
	UESR        = "USER"
	SUPER_ADMIN = "SUPER_ADMIN"
)

type UserDTO struct {
	UserID         string `json:"user_id" bson:"user_id"`
	UserName       string `json:"user_name" bson:"user_name"`
	PersonalBio    string `json:"personal_bio" bson:"personal_bio"`
	ProfilePic     string `json:"profile_pic" bson:"profile_pic"` // store as URL or base64
	Email          string `json:"email" bson:"email"`
	PhoneNum       string `json:"phone_num" bson:"phone_num"` // validate format
	TelegramHandle string `json:"telegram_handle" bson:"telegram_handle"`
	Password       string `json:"password" bson:"password"` // securely hashed
	Role           string `json:"role" bson:"role"`         // values: "admin", "user", "super_admin"
}
type BlogDTO struct {
	BlogID      string    `json:"blog_id" bson:"blog_id"`
	OwnerID     string    `json:"owner_id" bson:"owner_id"` // references UserID
	Title       string    `json:"title" bson:"title"`
	Tags        []string  `json:"tags" bson:"tags"`
	Author      string    `json:"author" bson:"author"` // could be redundant with OwnerID
	Description string    `json:"description" bson:"description"`
	LastUpdate  time.Time `json:"last_update_time" bson:"last_update_time"`
}
type PopularityDTO struct {
	PopularityID string   `json:"popularity_id" bson:"popularity_id"`
	BlogID       string   `json:"blog_id" bson:"blog_id"`
	ViewCount    int      `json:"view_count" bson:"view_count"`
	Likes        []string `json:"likes" bson:"likes"`       // list of user_ids — deduplicate before insert
	Dislikes     []string `json:"dislikes" bson:"dislikes"` // list of user_ids — same
	Comments     []string `json:"comments" bson:"comments"` // comment IDs or plain content
}

func ChangeToDTOUserfunc(domainUser *domain.User) *UserDTO {
	return &UserDTO{
		UserID:         domainUser.UserID,
		UserName:       domainUser.UserName,
		PersonalBio:    domainUser.PersonalBio,
		ProfilePic:     domainUser.ProfilePic,
		Email:          domainUser.Email,
		PhoneNum:       domainUser.PhoneNum,
		TelegramHandle: domainUser.TelegramHandle,
		Password:       domainUser.Password,
		Role:           domainUser.Role,
	}
}
func ChangeToDomainUser(udto *UserDTO) *domain.User {
	return &domain.User{
		UserID:         udto.UserID,
		UserName:       udto.UserName,
		PersonalBio:    udto.PersonalBio,
		ProfilePic:     udto.ProfilePic,
		Email:          udto.Email,
		PhoneNum:       udto.PhoneNum,
		TelegramHandle: udto.TelegramHandle,
		Password:       udto.Password,
		Role:           udto.Role,
	}
}

func ChangeToDomainBlog(bdto *BlogDTO) *domain.Blog {
	return &domain.Blog{
		BlogID:      bdto.BlogID,
		OwnerID:     bdto.OwnerID,
		Title:       bdto.Title,
		Tags:        bdto.Tags,
		Author:      bdto.Author,
		Description: bdto.Description,
		LastUpdate:  bdto.LastUpdate,
	}
}
func ChangeToDTOBlog(bdto *domain.Blog) *BlogDTO {
	return &BlogDTO{
		BlogID:      bdto.BlogID,
		OwnerID:     bdto.OwnerID,
		Title:       bdto.Title,
		Tags:        bdto.Tags,
		Author:      bdto.Author,
		Description: bdto.Description,
		LastUpdate:  bdto.LastUpdate,
	}
}
func ChangeToDomainPopularity(pdto *PopularityDTO) *domain.Popularity {
	return &domain.Popularity{
		PopularityID: pdto.PopularityID,
		BlogID:       pdto.BlogID,
		ViewCount:    pdto.ViewCount,
		Likes:        pdto.Likes,
		Dislikes:     pdto.Dislikes,
		Comments:     pdto.Comments,
	}
}
func ChangeToDTOPopularity(pdto *domain.Popularity) *PopularityDTO {
	return &PopularityDTO{
		PopularityID: pdto.PopularityID,
		BlogID:       pdto.BlogID,
		ViewCount:    pdto.ViewCount,
		Likes:        pdto.Likes,
		Dislikes:     pdto.Dislikes,
		Comments:     pdto.Comments,
	}
}
