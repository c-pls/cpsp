# cpsp

*Yet* another Spotify CLI tool written in Go. <br>
Compatible with ***any*** OS.


# Installation
- ### Prerequisite: [Go](https://go.dev/)
    
    - git clone [https://github.com/c-pls/cpsp.git](github.com/c-pls/cpsp.git) 
    - cd cpsp
    - go build && go install


# Connect to Spotify App
`cpsp` need to connect to Spotify's API.
    
-  Go to the [Spotify dashboard](https://developer.spotify.com/dashboard/)
- Click Create an app (You now can see your Client ID and Client Secret)

- Now click Edit Settings

- Add http://localhost:8888/callback to the Redirect URIs

- Scroll down and click Save

- Open terminal (cmd or Powershell in Windows)

- Run cpsp

- Enter your Client ID

- Enter your Client Secret

- You will be redirected to an official Spotify webpage to ask you for permissions.

- Accepting the permission

# Usage

````
cpsp play                       Resume playback
cpsp pause                      Stop playback

cpsp next                       Skip to the next song in a playlist.
cpsp prev                       Return to the previous song in a playlist.
cpsp replay                     Replays the current track from the beginning.
cpsp peek <percentage>          Jump to a specific time (in percentage) in the current song.

cpsp vol up                     Increases the volume by 10%.
cpsp vol down                   Decreases the volume by 10%.
cpsp vol <amount>               Sets the volume to an amount between 0 and 100.
cpsp vol                        Shows the current volume.

cpsp status                     Shows the current playing status

cpsp repeat                     Toggle repeat playback
cpsp repeat on                  Turn on repeat playback
cpsp repeat off                 Turn on repeat playback

cpsp shuffle                    Toggle shuffle playback
cpsp shuffle on                 Turn on shuffle playback
cpsp shuffle off                Turn on shuffle playback

cpsp lyric                      See the lyric of current song
````

# Copyright
