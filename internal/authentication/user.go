package authentication

// User
type User struct {
	ID       int
	Language string
}

func (u *User) IsEmpty() bool {
	return u.ID == 0
}
