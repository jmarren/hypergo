package pages

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/a-h/templ"
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

// type ComponentHandler func(rw *RW) (templ.Component, error)

func FavSongForm(rw *hypergo.RW) (templ.Component, error) {
	queryParams := rw.URL.Query().Get("errors")

	errors := strings.Split(queryParams, ",")

	return views.FavoriteSongForm(errors), nil
}

func postFavoriteSong(rw *hypergo.RW) error {
	favSong := rw.FormValue("fav-song")
	fmt.Printf("fav-song = %s\n", favSong)

	queryParams := []string{}

	if len(favSong) < 3 {
		queryParams = append(queryParams, url.QueryEscape("song must be > 3 characters long"))
	}

	if len(queryParams) > 0 {
		rw.Location("/songs/favorite?errors=" + strings.Join(queryParams, ","))
		return nil
	}

	rw.Location("/songs/blackbird")
	return nil
}

func init() {

	favoriteSong = ""
	// create the router
	SongsRouter = hypergo.NewRouter("#songs-component")

	// wrap
	SongsRouter.Wrap(hypergo.UnsafeWrapFunc(views.Songs))

	SongsRouter.GetComponent("blackbird", hypergo.SimpleComponent(views.Blackbird))

	SongsRouter.GetComponent("favorite", hypergo.NewComponent(FavSongForm))

	SongsRouter.Post("favorite", postFavoriteSong)

	YesterdayRouter := hypergo.NewRouter("#yesterday-component")
	YesterdayRouter.Wrap(hypergo.UnsafeWrapFunc(views.Yesterday))
	YesterdayRouter.GetComponent("stats", hypergo.SimpleComponent(views.YesterdayStats))
	YesterdayRouter.GetComponent("artwork", hypergo.SimpleComponent(views.YesterdayArtwork))

	SongsRouter.SubRouter("yesterday/", YesterdayRouter)

}
