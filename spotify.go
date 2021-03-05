package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

const spotifyAPIBase = "https://api.spotify.com/v1"

var spotifyAccessToken = ""

type spotifyTokenResponse struct {
	AccessToken      string `json:"access_token"`
	Error            string
	ErrorDescription string `json:"error_description"`
}

type spotifyPlaylistTracksResponse struct {
	Items []spotifyPlaylistItem
	Total int
	Error *spotifyRequestError
}

type spotifyRequestError struct {
	Status  int
	Message string
}

type spotifyPlaylistItem struct {
	Track *spotifyPlaylistTrack
}

type spotifyPlaylistTrack struct {
	ID string
}

// addTracksToQueue adds all of the tracks of the current listening device
func addTracksToQueue(trackList []string) {
	log.Println("Adding", len(trackList), "tracks to queue")
	bar := pb.StartNew(len(trackList))
	for _, trackID := range trackList {
		res := spotifyPOST("/me/player/queue?uri=spotify:track:"+trackID, nil)
		if len(res) > 0 {
			spotifyErr := spotifyRequestError{}
			err := json.Unmarshal([]byte(res), &spotifyErr)
			if err != nil {
				log.Println("addTracksToQueue() Error when attempting to parse an HTTP request")
				log.Fatalln(err)
				break
			}
			log.Println("addTracksToQueue() Error from API, aborting:", spotifyErr.Message)
			log.Fatalln(err)
			break
		}
		bar.Increment()
		time.Sleep(time.Millisecond * 20)
	}
	bar.Finish()
}

func getAllTracks() (trackList []string) {
	for _, playlistID := range playlistIds {
		log.Println("Getting tracks for playlist " + playlistID)
		trackIDs := []string{}
		getTracksFromPlaylist(playlistID, 0, &trackIDs)
		trackList = append(trackList, trackIDs...)
		log.Println("Fetched " + strconv.Itoa(len(trackIDs)) + " tracks from " + playlistID)
	}

	// shuffle it twice because why not
	rand.Shuffle(len(trackList), func(i, j int) { trackList[i], trackList[j] = trackList[j], trackList[i] })
	rand.Shuffle(len(trackList), func(i, j int) { trackList[i], trackList[j] = trackList[j], trackList[i] })

	return
}

func getTracksFromPlaylist(playlistID string, offset int, trackIDs *[]string) {
	// fetch the tracks from the playlist
	output := spotifyGET("/playlists/" + playlistID + "/tracks?market=US&limit=100&offset=" + strconv.Itoa(offset) + "&fields=items(track(id)),total")
	response := spotifyPlaylistTracksResponse{}
	err := json.Unmarshal([]byte(output), &response)
	if err != nil {
		log.Println("getTracksFromPlaylist() Error when attempting to parse an HTTP request")
		log.Fatalln(err)
	}

	if response.Error != nil {
		log.Println("getTracksFromPlaylist() Error from API")
		log.Fatalln(response.Error.Message)
	}

	for _, item := range response.Items {
		*trackIDs = append(*trackIDs, item.Track.ID)
	}

	// if there are more than 100 tracks, fetch an additional 100
	if response.Total > 100 && len(*trackIDs) != response.Total {
		getTracksFromPlaylist(playlistID, offset+100, trackIDs)
	}
}

// getAccessToken sends a request to spotify to receive an access token from oauth code
func getAccessToken(code string) {
	data := "grant_type=authorization_code&redirect_uri=" + redirectURI + "&code=" + code + "&client_id=" + clientID + "&client_secret=" + clientSecret
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data))
	if err != nil {
		log.Println("getAccessToken() Error when attempting to create an HTTP request")
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("getAccessToken() Error when attempting to send an HTTP request")
		log.Fatalln(err)
	}

	result := spotifyTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("getAccessToken() Error when attempting to decode an HTTP request")
		log.Fatalln(err)
	}

	if len(result.Error) > 0 {
		log.Println("getAccessToken() Error returned from API")
		log.Fatalln(result.Error + " - " + result.ErrorDescription)
	}

	spotifyAccessToken = result.AccessToken
}

// spotifyGET sends an authenticated GET request to Spotify
func spotifyGET(path string) string {
	req, err := http.NewRequest("GET", spotifyAPIBase+path, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "Bearer "+spotifyAccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("spotifyGET() Error when attempting to create an HTTP request")
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("spotifyGET() Error when attempting to send an HTTP request")
		log.Fatalln(err)
	}
	return string(body)
}

// spotifyPOST sends an authenticated POST request to Spotify
func spotifyPOST(path string, data io.Reader) string {
	req, err := http.NewRequest("POST", spotifyAPIBase+path, data)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "Bearer "+spotifyAccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("spotifyPOST() Error when attempting to create an HTTP request")
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("spotifyPOST() Error when attempting to send an HTTP request")
		log.Fatalln(err)
	}
	return string(body)
}
