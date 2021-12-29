package cmd

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
	"log"
	"strconv"
	"time"
)

func modifyVolume(volume int) {
	currentDeviceId := getPlayBackState().Device.ID
	ops := spotify.PlayOptions{DeviceID: &currentDeviceId}

	err := client.VolumeOpt(context.Background(), volume, &ops)
	handleError(err)
}

func play() {
	err := client.Play(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func pause() {
	err := client.Pause(context.Background())
	handleError(err)
}

func nextTrack() {
	err := client.Next(context.Background())
	handleError(err)
}

func previousTrack() {
	err := client.Previous(context.Background())
	handleError(err)
}

func getCurrentUserTopArtist() *spotify.FullArtistPage {
	res, err := client.CurrentUsersTopArtists(context.Background())
	handleError(err)
	return res
}

func getUserPlayList() *spotify.SimplePlaylistPage {
	res, err := client.CurrentUsersPlaylists(context.Background())
	handleError(err)
	return res
}

func getPlayListTrack(playlistID string) *spotify.PlaylistTrackPage {
	results, err := client.GetPlaylistTracks(context.Background(), spotify.ID(playlistID))
	handleError(err)
	if results == nil {
		log.Fatal()
	}
	return results
}

func playSongInAlbum(uri spotify.URI, offsetPos int) {
	offset := spotify.PlaybackOffset{
		Position: offsetPos,
	}

	opt := &spotify.PlayOptions{
		PlaybackContext: &uri,
		PlaybackOffset:  &offset,
		PositionMs:      0,
	}
	err := client.PlayOpt(context.Background(), opt)
	handleError(err)
}

func getCurrentPlaying() *spotify.CurrentlyPlaying {
	currentTrack, err := client.PlayerCurrentlyPlaying(context.Background())
	handleError(err)
	if !currentTrack.Playing {
		log.Fatal("Your device is not active")
	}
	return currentTrack
}

func getAvailableDevice() []spotify.PlayerDevice {
	listDevices, err := client.PlayerDevices(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(listDevices) == 0 {
		log.Fatal("There is no active device")
	}
	return listDevices
}

func getPlayBackState() *spotify.PlayerState {
	state, err := client.PlayerState(context.Background())
	handleError(err)
	return state
}

func changePlayBackDevice() {
	list := getAvailableDevice()
	if list == nil {
		return
	}
	if len(list) < 2 {
		return
	}
	for idx, e := range list {
		var running string
		if e.Active {
			running = "(Running)"
		}
		fmt.Printf("%d: %s - %s %s\n", idx+1, e.Name, e.Type, running)
	}
	x := getUserInput("Chose the device you want to play on")
	err := client.TransferPlayback(context.Background(), list[x-1].ID, true)
	handleError(err)
	fmt.Println("Currently running on", list[x-1].Name, "-", list[x-1].Type)
}

func toggleRepeatMode(mode string) {
	currentDeviceState := getPlayBackState()
	id := currentDeviceState.Device.ID
	opts := spotify.PlayOptions{DeviceID: &id}
	var state string
	if mode == "toggle" {
		if currentDeviceState.RepeatState == "track" {
			state = "off"
		} else {
			state = "track"
		}
	} else if mode == "on" {
		state = "track"
	} else if mode == "off" {
		state = "off"
	}
	err := client.RepeatOpt(context.Background(), state, &opts)
	handleError(err)
}

func togglePlaybackShuffle(mode string) {
	currentDeviceState := getPlayBackState()
	id := currentDeviceState.Device.ID
	opts := spotify.PlayOptions{DeviceID: &id}
	var state bool
	if mode == "toggle" {
		state = !currentDeviceState.ShuffleState
	} else if mode == "on" {
		state = true
	} else if mode == "off" {
		state = false
	}
	err := client.ShuffleOpt(context.Background(), state, &opts)
	handleError(err)
}

func seek(percentage int) {
	currentDeviceState := getPlayBackState()
	id := currentDeviceState.Device.ID
	opts := spotify.PlayOptions{DeviceID: &id}

	position := currentDeviceState.Item.Duration * percentage / 100

	err := client.SeekOpt(context.Background(), position, &opts)
	handleError(err)

}

func getUserInput(prompt string) int {
	fmt.Println(prompt)
	var input string
	fmt.Scan(&input)
	res, err := strconv.Atoi(input)
	handleError(err)
	return res
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func convertTime(expiry string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	layout = "2006-01-02T15:04:05.999999999Z07:00"
	t, err := time.Parse(layout, expiry)
	if err != nil {
		fmt.Println(err)
	}
	return t
}
