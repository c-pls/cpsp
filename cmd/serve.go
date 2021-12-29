package cmd

import (
	"context"
	"fmt"
	"github.com/c-pls/golyrics"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(cmdStatus)
	rootCmd.AddCommand(cmdModifyVolume)
	rootCmd.AddCommand(cmdPause)
	rootCmd.AddCommand(cmdPlay)
	rootCmd.AddCommand(cmdNextTrack)
	rootCmd.AddCommand(cmdPreviousTrack)
	rootCmd.AddCommand(cmdChangePlayBackDevice)
	rootCmd.AddCommand(cmdPlayList)
	rootCmd.AddCommand(cmdLyric)
	rootCmd.AddCommand(cmdRepeat)
	rootCmd.AddCommand(cmdShuffle)
	rootCmd.AddCommand(cmdSeek)
	rootCmd.AddCommand(cmdConfig)
	cmdLyric.PersistentFlags().String("trans", "", "Which language?")
	handleRefreshData()

}

var client *spotify.Client

var cmdLyric = &cobra.Command{
	Use:   "lyric",
	Short: "See the lyric of current song",
	Run: func(cmd *cobra.Command, args []string) {
		lang, err := cmd.Flags().GetString("trans")
		cobra.CheckErr(err)
		//handleError(err)

		currentTrack := getCurrentPlaying()
		songName := strings.Split(currentTrack.Item.Name, "-")[0]
		var artist string
		for _, e := range currentTrack.Item.Artists {
			artist += e.Name + ","
		}
		artist = artist[:len(artist)-1]
		fmt.Printf("%s - %s\n \n", currentTrack.Item.Name, artist)
		if lang == "" {
			fmt.Println(golyrics.GetLyrics(songName, currentTrack.Item.Artists[0].Name))
		} else {
			res := golyrics.GetLyricsWithTranslate(songName, currentTrack.Item.Artists[0].Name, lang)
			fmt.Println(res["translation"])
		}
	},
}

var cmdConfig = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		configData()
	},
}

var cmdPlayList = &cobra.Command{
	Use:   "playlist",
	Short: "See the detail of your playlist",
	Run: func(cmd *cobra.Command, args []string) {
		userPlayList := getUserPlayList()
		for idx, e := range userPlayList.Playlists {
			fmt.Printf("%d: %s - Total Tracks: %d \n", idx+1, e.Name, e.Tracks.Total)
		}
		x := getUserInput("Which do you want")

		playListTracks := getPlayListTrack(string(userPlayList.Playlists[x-1].ID))
		for idx, e := range playListTracks.Tracks {
			fmt.Printf("%d: %s - %s\n", idx+1, e.Track.Name, e.Track.Artists[0].Name)
		}
		trackOffset := getUserInput("Choose a song")

		playSongInAlbum(userPlayList.Playlists[x-1].URI, trackOffset-1)
	},
}
var cmdStatus = &cobra.Command{
	Use:   "status",
	Short: "Check the current playing track",
	Long:  "See which you are listening to",
	Run: func(cmd *cobra.Command, args []string) {
		currentTrack := getCurrentPlaying()

		fmt.Printf("Now playing: %s - %s\n", currentTrack.Item.Name, currentTrack.Item.Artists[0].Name)
	},
}

var cmdModifyVolume = &cobra.Command{
	Use:   "vol [PERCENTAGE] <up>|<down>",
	Short: "Change the volume of current device",
	Run: func(cmd *cobra.Command, args []string) {
		currentVolume := getPlayBackState().Device.Volume
		if len(args) == 0 {
			fmt.Println(currentVolume)
			return
		}
		action := args[0]
		if strings.EqualFold(action, "up") {
			modifyVolume(currentVolume + 10)
			return
		} else if strings.EqualFold(action, "down") {
			modifyVolume(currentVolume - 10)
			return
		}
		volume, err := strconv.Atoi(action)
		handleError(err)
		modifyVolume(volume)
	},
}
var cmdPlay = &cobra.Command{
	Use:   "play",
	Short: "Resume playback ",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		play()
	},
}
var cmdPause = &cobra.Command{
	Use:   "pause",
	Short: "Pause playback",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pause()
	},
}

var cmdNextTrack = &cobra.Command{
	Use:   "next",
	Short: "Next track",
	Run: func(cmd *cobra.Command, args []string) {
		nextTrack()
	},
}

var cmdPreviousTrack = &cobra.Command{
	Use:   "prev",
	Short: "Previous track",
	Run: func(cmd *cobra.Command, args []string) {
		previousTrack()
	},
}

var cmdChangePlayBackDevice = &cobra.Command{
	Use:   "switch",
	Short: "Listen to another device",
	Run: func(cmd *cobra.Command, args []string) {
		changePlayBackDevice()
	},
}

var cmdRepeat = &cobra.Command{
	Use:   "repeat",
	Short: "Repeat current listen song",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			toggleRepeatMode("toggle")
			return
		}
		toggleRepeatMode(args[0])
	},
}

var cmdShuffle = &cobra.Command{
	Use:   "shuffle",
	Short: "Shuffle [on | off]",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			togglePlaybackShuffle("toggle")
			return
		}
		togglePlaybackShuffle(args[0])
	},
}

var cmdSeek = &cobra.Command{
	Use:   "seek",
	Short: "Seek to the percentage of the playback  ",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		percentage, err := strconv.Atoi(args[0])
		handleError(err)
		seek(percentage)
	},
}

func createClient() {
	var token = &oauth2.Token{
		//AccessToken:  viper.GetString("access_token"),
		TokenType:    viper.GetString("token_type"),
		RefreshToken: viper.GetString("refresh_token"),
		//Expiry:       expiry,
	}
	// create the client
	client = spotify.New(auth.Client(context.Background(), token))
}
