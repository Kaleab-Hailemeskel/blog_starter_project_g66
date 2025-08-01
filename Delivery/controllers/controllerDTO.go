package controllers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ADMIN       = "ADMIN"
	UESR        = "USER"
	SUPER_ADMIN = "SUPER_ADMIN"
)

type UserUnverifiedDTO struct{

	UserName string `json: "username" bson: "username"`
	Email string `json: "email" bson: "email"`
	Password string `json: "password" bson: "password"`
	OTP string `json:"otp" bson: "otp"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`

}
type UserDTO struct {
	UserID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName       string             `json:"user_name" bson:"user_name"`
	PersonalBio    string             `json:"personal_bio" bson:"personal_bio"`
	ProfilePic     string             `json:"profile_pic" bson:"profile_pic"` // store as URL or base64
	Email          string             `json:"email" bson:"email"`
	PhoneNum       string             `json:"phone_num" bson:"phone_num"` // validate format
	TelegramHandle string             `json:"telegram_handle" bson:"telegram_handle"`
	Password       string             `json:"password" bson:"password"` // securely hashed
	Role           string             `json:"role" bson:"role"`         // values: "admin", "user", "super_admin"
}
type BlogDTO struct {
	BlogID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OwnerID     primitive.ObjectID `json:"owner_id" bson:"owner_id"` // references UserID
	Title       string             `json:"title" bson:"title"`
	Tags        []string           `json:"tags" bson:"tags"`
	Author      string             `json:"author" bson:"author"` // could be redundant with OwnerID
	Description string             `json:"description" bson:"description"`
	LastUpdate  time.Time          `json:"last_update_time" bson:"last_update_time"`
}
type PopularityDTO struct {
	PopularityID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BlogID       primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	ViewCount    int                `json:"view_count" bson:"view_count"`
	Likes        []string           `json:"likes" bson:"likes"`       // list of user_ids — deduplicate before insert
	Dislikes     []string           `json:"dislikes" bson:"dislikes"` // list of user_ids — same
	Comments     []*CommentDTO      `json:"comments" bson:"comments"` // comment IDs or plain content
}

type CommentDTO struct {
	UserID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"user_name" bson:"user_name"`
	Comment  string             `json:"comment" bson:"comment"`
}

func changeToDomainVerification(udto *UserUnverifiedDTO) *domain.UserUnverified{
    return &domain.UserUnverified{
        Email:     udto.Email,
        OTP:       udto.OTP,
        ExpiresAt: udto.ExpiresAt,
    }
}