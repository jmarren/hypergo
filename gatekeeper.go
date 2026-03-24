

package hypergo

import (
	"net/http"
	"strconv"
	"fmt"
	"slices"
	"net/mail"
)


type User struct {
	FirstName string
	LastName string
	Age int
	Email string
	BirthMonth string
}


func NewUser(r *http.Request) (*User, []error) {
	var err error
	var errs []error
	x := new(User)

	x.FirstName = r.FormValue("FirstName")
	
	if len(x.FirstName) < 3 {
		errs = append(errs, fmt.Errorf("first name must be >= 3 characters long"))
	}

	
	if len(x.FirstName) > 10 {
		errs = append(errs, fmt.Errorf("first name must be <= 10 characters long"))
	}

	
	
        

	x.LastName = r.FormValue("LastName")
	
	if len(x.LastName) < 3 {
		errs = append(errs, fmt.Errorf("last name must be >= 3 characters long"))
	}

	
	if len(x.LastName) > 10 {
		errs = append(errs, fmt.Errorf("last name must be <= 10 characters long"))
	}

	
	
        


	err = NoWhiteSpace(x.LastName) 
	if err != nil {
		errs = append(errs, err)
	}




	x.Age, err = strconv.Atoi(r.FormValue("Age"))
	if err != nil {
		errs = append(errs, fmt.Errorf("age must be int"))
	}
	
	
	if x.Age < 18 {
		errs = append(errs, fmt.Errorf("age must be >= 18"))
	} 

	
	if x.Age > 100 {
		errs = append(errs, fmt.Errorf("age must be <= 100"))
	}

	
        


	x.Email = r.FormValue("Email")
	
	
	
	
    address, err := mail.ParseAddress(x.Email)
    if err != nil {
	    errs = append(errs, fmt.Errorf("invalid email provided for Email"))
    } else {
	x.Email = address.Address
    }
	

        

	x.BirthMonth = r.FormValue("BirthMonth")
	
	
	
	BirthMonthOptions := []string{ "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	
	if !slices.Contains(BirthMonthOptions, x.BirthMonth) {
		errs = append(errs, fmt.Errorf("invalid option selected for birth month"))
	}


	
        

	return x, errs

}



type Search struct {
	QueryString string
	ResultCount int
}


func NewSearch(r *http.Request) (*Search, []error) {
	var err error
	var errs []error
	x := new(Search)

	x.QueryString = r.FormValue("QueryString")
	
	if len(x.QueryString) < 1 {
		errs = append(errs, fmt.Errorf("search must be >= 1 characters long"))
	}

	
	if len(x.QueryString) > 24 {
		errs = append(errs, fmt.Errorf("search must be <= 24 characters long"))
	}

	
	
        

	x.ResultCount, err = strconv.Atoi(r.FormValue("ResultCount"))
	if err != nil {
		errs = append(errs, fmt.Errorf("result count must be int"))
	}
	
	
	
	
	ResultCountOptions := []int{ 10, 20, 50, 100 }
	
	if !slices.Contains(ResultCountOptions, x.ResultCount) {
		errs = append(errs, fmt.Errorf("invalid option selected for result count"))
	}


        


	return x, errs

}



