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
	// Age            uint8  `validate:"gte=0,lte=130"`
	// Email          string `validate:"required,email"`
	// Gender         string `validate:"oneof=male female prefer_not_to"`
	// FavouriteColor string `validate:"iscolor"` // alias for 'hexcolor|rgb|rgba|hsl|hsla'
}

// use a single instance of Validate, it caches struct info
// var validate *validator.Validate

// type ComponentHandler func(rw *RW) (templ.Component, error)

func FavSongForm(rw *hypergo.RW) (templ.Component, error) {
	queryParams := rw.URL.Query().Get("errors")

	errors := strings.Split(queryParams, ",")

	return views.FavoriteSongForm(errors), nil
}

func postFavoriteSong(rw *hypergo.RW) error {
	favSong := rw.FormValue("fav-song")
	fmt.Printf("fav-song = %s\n", favSong)

	validator := hypergo.NewRequestValidator().Use("fav-song", hypergo.RequireMinLen(2), hypergo.RequireMaxLen(1), hypergo.RequireInt)

	vals, errs := validator.Validate(rw.Request)

	if len(errs) != 0 {
		fmt.Printf("errs = %v\n", errs)
	}

	fmt.Printf("fav-song = %d\n", vals["fav-song"].Int())

	// usernameValidator := hypergo.UsernameValidator(
	// 	rw.FormValue("fav-song"),
	// 	hypergo.RequireMaxLen(2),
	// 	hypergo.RequireMaxLen(5),
	// 	hypergo.NoWhiteSpace)

	// user := &User{
	// 	FirstName: "Badger",
	// 	// LastName:  "Smith",
	// }
	//
	// err := validate.Struct(user)
	//
	// if err != nil {
	//
	// 	// this check is only needed when your code could produce
	// 	// an invalid value for validation such as interface with nil
	// 	// value most including myself do not usually have code like this.
	// 	var invalidValidationError *validator.InvalidValidationError
	// 	if errors.As(err, &invalidValidationError) {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	//
	// 	var validateErrs validator.ValidationErrors
	// 	if errors.As(err, &validateErrs) {
	// 		for _, e := range validateErrs {
	// 			fmt.Printf("e.Namespace() = %s\n", e.Namespace())
	// 			fmt.Printf("e.Field() = %s\n", e.Field())
	// 			fmt.Printf("e.StructNamespace() = %s\n", e.StructNamespace())
	// 			fmt.Printf("e.StructField() =  %s\n", e.StructField())
	// 			fmt.Println(e.Tag())
	// 			fmt.Println(e.ActualTag())
	// 			fmt.Println(e.Kind())
	// 			fmt.Println(e.Type())
	// 			fmt.Println(e.Value())
	// 			fmt.Println(e.Param())
	// 			fmt.Println()
	// 		}
	// 	}
	// }
	//
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

	// validate = validator.New(validator.WithRequiredStructEnabled())

	favoriteSong = ""
	// create the router
	SongsRouter = hypergo.NewRouter("#songs-component")

	// wrap
	SongsRouter.Wrap(hypergo.UnsafeWrapFunc(views.Songs))

	SongsRouter.GetComponent("blackbird", hypergo.SimpleComponent(views.Blackbird))

	SongsRouter.GetComponent("favorite", hypergo.NewComponent(FavSongForm))

	SongsRouter.Post("favorite", hypergo.NewHandler(postFavoriteSong))

	YesterdayRouter := hypergo.NewRouter("#yesterday-component")
	YesterdayRouter.Wrap(hypergo.UnsafeWrapFunc(views.Yesterday))
	YesterdayRouter.GetComponent("stats", hypergo.SimpleComponent(views.YesterdayStats))
	YesterdayRouter.GetComponent("artwork", hypergo.SimpleComponent(views.YesterdayArtwork))

	SongsRouter.SubRouter("yesterday/", YesterdayRouter)

}
