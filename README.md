# Spotify Playlist Combiner

The project lets you combine multiple Spotify playlists and add them all to your queue.

This requires a premium Spotify account.

This was created for personal use and has minimal error handling, this is not a "product", watch the console for progress and crashes.

### Usage

1. Create an app account in the [developers portal](https://developer.spotify.com/), use `http://localhost:8080/callback` as the Redirect URI
2. Clear your current Spotify queue
3. Run the project (either using the compiled versions releases or by using `go run .`)
4. Start listening to a random song 
5. Fill in the input with the information from Spotify
6. Visit http://localhost:8080/

Example Inputs:

```
Spotify Client ID: 
ID HERE
Spotify Client Secret: 
SECRET HERE
Redirect URI: http://localhost:8080/callback
Comma seperated spotify playlist IDs:
37i9dQZF1DX2L0iB23Enbq,37i9dQZF1DWUa8ZRTfalHk
```

### Dependencies

- [Go v1.16+](https://golang.org/dl/)
