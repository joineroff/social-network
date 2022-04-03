package dto

type RemoveFriendInputDto struct {
	FriendID      string `json:"friendID"`
	CurrentUserID string `json:"-"`
}
