package model

import "time"

type User struct {
	Id           int64     `json:"id,omitempty"`
	FirstName    string    `json:"firstname"`
	MiddleName   string    `json:"middlename"`
	LastName     string    `json:"lastname"`
	Mobile       string    `json:"mobile"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordhash"`
	RegisteredAt time.Time `json:"registeredat"`
	LastLogin    time.Time `json:"lastlogin"`
	Intro        string    `json:"intro"`
	Profile      string    `json:"profile"`
	Posts        *[]Post   `json:"posts,omitempty"`
}
