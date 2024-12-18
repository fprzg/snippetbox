package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
	//"fmt"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

const emailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var EmailRX = regexp.MustCompile(emailRegex)

func (v *Validator) Valid() bool {
	//fmt.Printf("\t\t%d errors\n", len(v.FieldErrors))
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
		//fmt.Printf("\t\tadded error '%s'\n", key)
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) CheckPassword(password, field string) {
	v.CheckField(NotBlank(password), field, "This field can't be empty")
	v.CheckField(MinChars(password, 8), field, "This field must be at least 8 characters long")
}

func (v *Validator) CheckEmail(email, field string) {
	v.CheckField(NotBlank(email), field, "This field cannot be blank")
	v.CheckField(Matches(email, EmailRX), field, "This field must be a valid email address")
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func InsideIntervalChars(min int, value string, max int) bool {
	return MinChars(value, min) && MaxChars(value, max)
}

func Equal[T comparable](value T, compareTo T) bool {
	return value == compareTo
}

func NotEqual[T comparable](value T, compareTo T) bool {
	return value != compareTo
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
