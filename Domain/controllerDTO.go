package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenDTO struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	Token     string    `json:"token" bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type UserUnverifiedDTO struct {
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	OTP       string    `json:"otp" bson:"otp"`
	Password  string    `json:"password" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type UserDTO struct {
	UserID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName       string             `json:"username" bson:"username"`
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
	PopularityID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BlogID          primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	PopularityValue int                `json:"popularity_value" bson:"popularity_value"` // it should be either ascending or descending or non(zero)
	ViewCount       int                `json:"view_count" bson:"view_count"`
	Likes           []string           `json:"likes" bson:"likes"`       // list of user_ids — deduplicate before insert
	Dislikes        []string           `json:"dislikes" bson:"dislikes"` // list of user_ids — same
	Comments        []*CommentDTO      `json:"comments" bson:"comments"` // comment IDs or plain content
}

type CommentDTO struct {
	CommentID primitive.ObjectID `json:"_id," bson:"_id,"`
	OwnerID   primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	UserName  string             `json:"user_name" bson:"user_name"`
	Comment   string             `json:"comment" bson:"comment"`
}

type UpdateProfileDTO struct {
	UserName       string `json:"username" bson:"username"`
	PersonalBio    string `json:"personal_bio" bson:"personal_bio"`
	ProfilePic     string `json:"profile_pic" bson:"profile_pic"` // store as URL or base64
	Email          string `json:"email" bson:"email"`
	PhoneNum       string `json:"phone_num" bson:"phone_num"` // validate format
	TelegramHandle string `json:"telegram_handle" bson:"telegram_handle"`
	Password       string `json:"password" bson:"password"` // securely hashed
}

type AICommentDTO struct {
	UserMessage string `json:"usermessage" bson:"usermessage"`
	Comment     string `json:"comment" bson:"comment"`
}

type AIBlogDTO struct {
	UserMessage string `json:"usermessage" bson:"usermessage"`
}

type AIUserBlogDTO struct {
	BlogDTO     *BlogDTO `json:"my_blog" bson:"my_blog"`
	UserMessage string  `json:"usermessage" bson:"usermessage"`
}

// type AIBlogGenerationDTO struct{
// 	UserMessage string `json:"usermessage" bson:"usermessage"`
// }
