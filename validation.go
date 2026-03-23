package hypergo

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Num int

const x Num = 10

//go:generate ./build/stringer -type=Pill
type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)

type RequestValidator interface {
	Validate(req *http.Request) (map[string]ValidatedPrimitive, []error)
}

type ValidatedPrimitive interface {
	String() string
	Int() int
}

type validatedPrimitive struct {
	value string
}

func newValidatedPrimitive(s string) *validatedPrimitive {
	return &validatedPrimitive{
		value: s,
	}
}

func (v *validatedPrimitive) String() string {
	return v.value
}

func (v *validatedPrimitive) Int() int {
	out, err := strconv.Atoi(v.value)

	if err != nil {
		panic(err)
	}
	return out
}

type PrimitiveMap map[string]ValidatedPrimitive

// func (p *PrimitiveMap) Get(s string) ValidatedPrimitive {
// 	return p[s]
// }

type requestValidator struct {
	validators map[string]StringValidator
}

func NewRequestValidator() *requestValidator {
	return &requestValidator{
		validators: make(map[string]StringValidator),
	}
}

func (r *requestValidator) Use(key string, validatorFuncs ...StringValidatorFunc) *requestValidator {
	// create key if it doesn't exist
	if r.validators[key] == nil {
		r.validators[key] = newStringValidator()
	}
	// use the validatorFuncs provided
	r.validators[key].Use(validatorFuncs...)
	return r
}

func (r *requestValidator) Validate(req *http.Request) (map[string]ValidatedPrimitive, []error) {
	outMap := make(map[string]ValidatedPrimitive)

	errs := []error{}
	for key, validator := range r.validators {
		val := req.FormValue(key)
		errs = append(errs, validator.Validate(val)...)
		outMap[key] = newValidatedPrimitive(val)
	}
	return outMap, errs
}

type StringValidatorFunc func(s string) error

type stringValidator struct {
	validators []StringValidatorFunc
}

func newStringValidator() *stringValidator {
	return &stringValidator{
		validators: []StringValidatorFunc{},
	}
}

func (s *stringValidator) Use(sfunc ...StringValidatorFunc) StringValidator {
	s.validators = append(s.validators, sfunc...)
	return s
}

type StringValidator interface {
	Validate(s string) []error
	Use(s ...StringValidatorFunc) StringValidator
}

func RequireMinLen(minLength int) StringValidatorFunc {
	return func(s string) error {
		if len(s) < minLength {
			return fmt.Errorf("must be at least %d characters long", minLength)
		}
		return nil
	}
}
func RequireMaxLen(maxLen int) StringValidatorFunc {
	return func(s string) error {
		if len(s) > maxLen {
			return fmt.Errorf("must be less than %d characters long", maxLen)
		}
		return nil
	}
}

func RequireInt(str string) error {
	if _, err := strconv.Atoi(str); err != nil {
		return fmt.Errorf("must be a number")
	}
	return nil
}

func (s *stringValidator) Validate(str string) []error {

	errors := []error{}

	for _, validator := range s.validators {
		err := validator(str)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func NoWhiteSpace(s string) error {
	if strings.Contains(s, " ") {
		return fmt.Errorf("must not include whitespace")
	}

	return nil
}

func UsernameValidator(value string, validators ...StringValidatorFunc) *stringValidator {
	return &stringValidator{
		validators: validators,
	}
}

func TryValidate() {
}

// andy := UsernameValidator("andy the manly", RequireMinLen(2), RequireMaxLen(5), NoWhiteSpace)

// vMap := NewValidationArr()
//
// // vMap.Push(andy)
// errors := vMap.Validate()
//
// fmt.Printf("errors = %v\n", errors)

//
// type validationArr struct {
// 	validators []Validator
// }
//
// func (v *validationArr) Push(validator Validator) *validationArr {
// 	v.validators = append(v.validators, validator)
// 	return v
// }
//
// func NewValidationArr() *validationArr {
// 	return &validationArr{
// 		validators: []Validator{},
// 	}
// }
// func (v *validationArr) Validate() []error {
//
// 	errors := []error{}
// 	for _, validator := range v.validators {
// 		errors = append(errors, validator.Validate()...)
// 	}
// 	return errors
// }
//
//
// type Validator interface {
// 	Validate() []error
// }
