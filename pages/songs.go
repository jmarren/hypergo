package pages

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/a-h/templ"

	// "github.com/go-playground/validator/v10"
	"github.com/jmarren/hypergo"
	"github.com/jmarren/hypergo/views"
)

var SongsRouter *hypergo.Router

var favoriteSong string

func SetFavorite(rw *hypergo.RW) (templ.Component, error) {
	fmt.Printf("fav-song = %s\n", rw.FormValue("fav-song"))
	// rw.Location("/songs/blackbird")
	return views.Blackbird(), nil
}

// User contains user information
type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

// use a single instance of Validate, it caches struct info
// var validate *validator.Validate
func FavSongForm(rw *hypergo.RW) templ.Component {
	queryParams := rw.URL.Query().Get("errors")

	errors := strings.Split(queryParams, ",")

	return views.FavoriteSongForm(errors)
}

func postFavoriteSong(rw *hypergo.RW) {
	favSong := rw.FormValue("fav-song")
	fmt.Printf("fav-song = %s\n", favSong)

	queryParams := []string{}

	if len(favSong) < 3 {
		queryParams = append(queryParams, url.QueryEscape("song must be > 3 characters long"))
	}

	if len(queryParams) > 0 {
		rw.Location("/songs/favorite?errors=" + strings.Join(queryParams, ","))
	}

	rw.Location("/songs/blackbird")
}
func init() {

	favoriteSong = ""
	// create the router
	SongsRouter = hypergo.NewRouter("#songs-component")

	// wrap
	SongsRouter.Wrap(hypergo.SimpleWrapper(views.Songs))

	SongsRouter.GetComponent("blackbird", hypergo.SimpleComponent(views.Blackbird))

	SongsRouter.GetComponent("favorite", FavSongForm)

	SongsRouter.Post("favorite", postFavoriteSong)

	YesterdayRouter := hypergo.NewRouter("#yesterday-component")
	YesterdayRouter.Wrap(hypergo.SimpleWrapper(views.Yesterday))
	YesterdayRouter.GetComponent("stats", hypergo.SimpleComponent(views.YesterdayStats))
	YesterdayRouter.GetComponent("artwork", hypergo.SimpleComponent(views.YesterdayArtwork))

	SongsRouter.SubRouter("yesterday/", YesterdayRouter)

}
