package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
)

const redirectURL = "http://localhost:8888/callback"

var clientId string
var clientSecret string

var scopes = []string{
	spotifyauth.ScopePlaylistModifyPrivate,
	spotifyauth.ScopePlaylistModifyPrivate,
	spotifyauth.ScopePlaylistModifyPublic,
	spotifyauth.ScopePlaylistReadCollaborative,
	spotifyauth.ScopeUserReadPrivate,
	spotifyauth.ScopeUserReadEmail,
	spotifyauth.ScopeUserReadPlaybackState,
	spotifyauth.ScopeUserModifyPlaybackState,
	spotifyauth.ScopeUserReadCurrentlyPlaying,
	spotifyauth.ScopeUserLibraryModify,
	spotifyauth.ScopeUserLibraryRead,
	spotifyauth.ScopeUserReadRecentlyPlayed,
	spotifyauth.ScopeUserTopRead,
	spotifyauth.ScopeStreaming}

var auth *spotifyauth.Authenticator
var ch = make(chan *spotify.Client)

// create a random string for the state
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var state = randSeq(16)

func configViper() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := home + "/.config/cpsp/"
	err = os.MkdirAll(configPath, os.ModePerm)
	cobra.CheckErr(err)
	if _, err := os.Stat(configPath + "config.yaml"); os.IsNotExist(err) {
		_, err := os.Create(configPath + "config.yaml")
		if err != nil {
			log.Fatal(err)
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("some error")
	}
	clientId = viper.GetString("client_id")
	clientSecret = viper.GetString("client_secret")
}

func handleClientInfo() {
	if clientId == "" || clientSecret == "" {
		configData()
	}
}

func handleRefreshData() {

	configViper()
	handleClientInfo()

	if viper.GetString("refresh_token") == "" {
		GetAccessToken()
	} else {
		createAuthenticator()
		createClient()
	}
}

func createAuthenticator() {
	auth = spotifyauth.New(
		spotifyauth.WithClientID(clientId),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(scopes...),
	)
}

func GetAccessToken() {
	createAuthenticator()

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for ", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	err := exec.Command("xdg-open", url).Run()
	if err != nil {
		log.Fatal("Couldn't open browser " + err.Error())
	}
	client = <-ch

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get the token", http.StatusForbidden)
		log.Fatalln(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s", st, state)
	}
	fmt.Fprintf(w, "Login Completed!")

	writeToConfigFile(token.RefreshToken)
	newClient := spotify.New(auth.Client(r.Context(), token))
	ch <- newClient
}

func writeToConfigFile(refreshToken string) {

	viper.Set("refresh_token", refreshToken)
	viper.Set("token_type", "Bearer")
	err := viper.WriteConfig()
	if err != nil {
		log.Println("Cannot write to the config file ", err.Error())
	}
}

func configData() {
	fmt.Printf("\nEnter client_id: ")
	fmt.Scanln(&clientId)
	fmt.Printf("\nEnter client_secret: ")
	fmt.Scanln(&clientSecret)

	viper.Set("client_id", clientId)
	viper.Set("client_secret", clientSecret)
	err := viper.WriteConfig()
	if err != nil {
		log.Println("Cannot write to the config file ", err.Error())
	}
	GetAccessToken()
	fmt.Println("Done")
}
