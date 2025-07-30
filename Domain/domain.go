package domain

import "time"

type User struct {
	UserID         string 
	UserName       string 
	PersonalBio    string 
	ProfilePic     string  // store as URL or base64
	Email          string 
	PhoneNum       string  // validate format
	TelegramHandle string 
	Password       string // securely hashed
	Role           string // values: "admin", "user", "super_admin"
}

type Blog struct {
	BlogID      string   
	OwnerID     string   // references UserID
	Title       string   
	Tags        []string 
	Author      string    // could be redundant with OwnerID
	Description string    
	LastUpdate  time.Time 
}

type Popularity struct {
	PopularityID string   
	BlogID       string   
	ViewCount    int      
	Likes        []string // list of user_ids — deduplicate before insert
	Dislikes     []string // list of user_ids — same
	Comments     []string // comment IDs or plain content
}
