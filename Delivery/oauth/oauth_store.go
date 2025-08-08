package oauth

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions/cookie"
"github.com/gin-contrib/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

func InitOAuth() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
	store := cookie.NewStore([]byte("randomString"))
store.Options(sessions.Options{
    Path:     "/",
    MaxAge:   86400 * 30,
    HttpOnly: true,
    Secure:   false,
    SameSite: http.SameSiteLaxMode,

})
	gothic.Store = store
	

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)
}
