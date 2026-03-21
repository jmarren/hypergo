package pages

import (
	"github.com/jmarren/hypergo"
	"github.com/jmarren/hypergo/views"
)

var UsersRouter *hypergo.Router

func init() {

	UsersRouter = hypergo.NewRouter("#users-component")
	UsersRouter.Wrap(hypergo.UnsafeWrapFunc(views.Users))

	UsersRouter.GetComponent("username", hypergo.SimpleComponent(views.Username))
	// usersRouter.Use(LoggerThree)
	UsersRouter.GetComponent("age", hypergo.SimpleComponent(views.Age))

}
