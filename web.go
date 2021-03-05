package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func primaryRoute(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect,
		"https://accounts.spotify.com/authorize?response_type=code&client_id="+clientID+"&redirect_uri="+redirectURI+"&scope=playlist-read-private user-modify-playback-state")
}

func callbackRouteError(c *gin.Context) {
	err := c.Query("error")
	code := c.Query("code")

	if err != "" {
		c.String(400, "Error: "+err)
		return
	}

	if code != "" {
		getAccessToken(code)
		tracks := getAllTracks()
		if len(tracks) == 0 {
			c.String(400, "fetched 0 tracks, tf?")
			return
		}
		addTracksToQueue(tracks)
		c.String(200, "done")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")

}
