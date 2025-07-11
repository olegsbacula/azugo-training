package models

import (
    "azugo.io/core/validation"
)

type User struct {
    Username string `json:"Username" validate:"required,max=20"`
    Password string `json:"Password" validate:"required,max=16,min=6"`  
    Email    string `json:"Email" validate:"required,email"`// add a custom tag to check for a unique email
}

// UserList is used to store multiple users at once
// swagger:model UserList
type UserList struct {
    Users []User `json:"users"`
}

type PublicUser struct{
    Username string `json:"Username" validate:"required,max=20"`
    Email    string `json:"Email" validate:"required,email"`// add a custom tag to check for a unique email
}

type PublicUserS struct {
    PublicUsers []PublicUser `json:"PublicUsers"`
}

func (u *User) Validate(validate *validation.Validate) error {
    return validate.Struct(u)
}
