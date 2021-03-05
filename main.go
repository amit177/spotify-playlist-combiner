package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const redirectURI = "http://localhost:8080/callback"

var clientID = ""
var clientSecret = ""
var playlistIds = []string{}

func main() {
	rand.Seed(time.Now().UnixNano())

	// get input from user for shit we need
	fmt.Println("Spotify Client ID: ")
	fmt.Scanln(&clientID)
	fmt.Println("Spotify Client Secret: ")
	fmt.Scanln(&clientSecret)
	fmt.Println("Redirect URI (set this in the Spotify app):", redirectURI)

	playlists := ""
	fmt.Println("Comma seperated spotify playlist IDs:")
	fmt.Scanln(&playlists)

	playlistIds = strings.Split(playlists, ",")

	// run the web that will handle the Spotify authentication
	r := gin.Default()
	r.GET("/", primaryRoute)
	r.GET("/callback", callbackRouteError)
	r.Run("127.0.0.1:8080")
}
