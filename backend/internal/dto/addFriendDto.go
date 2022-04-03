package dto

type AddFriendInputDto struct {
	FriendID      string `json:"friendID"`
	CurrentUserID string `json:"-"`
}
