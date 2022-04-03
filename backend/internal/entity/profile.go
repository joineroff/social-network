package entity

type Profile struct {
	User            *User
	IsFriend        bool
	FriendsQuantity int
}
