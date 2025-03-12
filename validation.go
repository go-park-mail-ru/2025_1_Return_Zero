package main

import (
	"strings"
)

const (
	VALIDLETTERS    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	VALIDLOGINCHARS = VALIDLETTERS + "_" + "0123456789"
	VALIDEMAILCHAES = VALIDLOGINCHARS + "@#.!"
)

type configValidator struct {
	configMinLen     int
	configMaxLen     int
	validConfigChars string
}

type Validator struct {
	passwordValidator *configValidator
	usernameValidator *configValidator
	emailValidator    *configValidator
}

func NewValidator() *Validator {
	pv := &configValidator{
		configMinLen:     4,
		configMaxLen:     25,
		validConfigChars: VALIDLOGINCHARS,
	}

	uv := &configValidator{
		configMinLen:     3,
		configMaxLen:     20,
		validConfigChars: VALIDLOGINCHARS,
	}

	ev := &configValidator{
		configMinLen:     5,
		configMaxLen:     30,
		validConfigChars: VALIDEMAILCHAES,
	}

	v := &Validator{
		passwordValidator: pv,
		usernameValidator: uv,
		emailValidator:    ev,
	}

	return v
}

func (cv *configValidator) validateField(data string) bool {
	if (len(data) < cv.configMinLen) || (len(data) > cv.configMaxLen) {
		return false
	}

	hasValidLetter := false
	for _, c := range data {
		if !strings.ContainsRune(cv.validConfigChars, c) {
			return false
		}
		if !hasValidLetter && strings.ContainsRune(VALIDLETTERS, c) {
			hasValidLetter = true
		}
	}

	return hasValidLetter
}

func (ev *configValidator) validateEmail(data string) bool {
	if len(data) < ev.configMinLen || len(data) > ev.configMaxLen {
		return false
	}

	if !strings.Contains(data, "@") || !strings.Contains(data, ".") {
		return false
	}

	for _, c := range data {
		if !strings.ContainsRune(ev.validConfigChars, c) {
			return false
		}
	}

	return true
}

func ValidateData(user *User) bool {
	username := user.Username
	password := user.Password
	email := user.Email
	validator := NewValidator()

	if !validator.usernameValidator.validateField(username) {
		return false
	}

	if !validator.passwordValidator.validateField(password) {
		return false
	}

	if !validator.emailValidator.validateEmail(email) {
		return false
	}
	return true
}
