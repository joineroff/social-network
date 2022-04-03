package dto

type ProfileDto struct {
	ID        string   `json:"id"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Gender    int8     `json:"gender"`
	Interests []string `json:"interests"`
	City      string   `json:"city"`
}
