
package hypergo

import (
	"net/http"
	"strconv"
	"fmt"
)


type User struct {
	FirstName string
	LastName string
	Age int
}


func NewUser(r *http.Request) (*User, []error) {
	var err error
	var errs []error
	x := new(User)

	x.FirstName = r.FormValue("FirstName")

	x.LastName = r.FormValue("LastName")

	x.Age, err = strconv.Atoi(r.FormValue("Age"))
	if err != nil {
		errs = append(errs, fmt.Errorf("Age must be int"))
	}

	return x, errs

}



