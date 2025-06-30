package models

import (
    "azugo.io/core/validation"
)

type User struct {
    Username string `json:"Username" validate:"required,max=20"`
    Password string `json:"Password" validate:"required,max=16,min=6"`
    Email    string `json:"Email" validate:"required,email"`
}

func (u *User) Validate(validate *validation.Validate) error {
    return validate.Struct(u)
}

type UserList struct {
    Users []User `json:"users"`
}