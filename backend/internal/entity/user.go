package entity

import "strings"

const (
	userGenderFemail = "Femail"
	userGenderMail   = "Mail"
)

type User struct {
	Login     string   `json:"login"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Gender    int8     `json:"gender"`
	Interests []string `json:"interests"`
	City      string   `json:"city"`

	ID       string `json:"-"`
	Password string `json:"-"`
}

func (e *User) GenderAsStirng() string {
	if e.Gender == 1 {
		return userGenderFemail
	}

	return userGenderMail
}

func (e *User) InterestsAsStirng() string {
	return strings.Join(e.Interests, ", ")
}
