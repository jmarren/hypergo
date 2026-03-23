
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
	validators []Validator
}

func NewUser(r *http.Request) (*User, []error) {
	errs := []error{}
	var err error
	x := new(User)

	x.FirstName = r.FormValue("FirstName")
      
	x.LastName = r.FormValue("LastName")
      
	x.Age, err = strconv.Atoi(r.FormValue("Age"))
        if err != nil {
		errs = append(errs, fmt.Errorf("Age must be a number"))
        }
      
      
	
      return x, errs

}






