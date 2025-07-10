//TO DO :

// isolate validation logic

package services

import (
	"azugo.io/core/validation"
	"example.com/project/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate = validation.New()
)

func ValidUser(u models.User) error {
	validate := validation.New()
	return u.Validate(validate)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}