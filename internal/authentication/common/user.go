package common

// User
type User struct {
	ID int
	// contacts
	Email string
	Phone string
	// address
	Address1 string
	Address2 string
	City     string
	State    string
	Country  string
	Zip      string
	// name & dob
	FirstName string
	LastName  string
	BirthDate int
	// language
	Language string
}

func (u *User) IsEmpty() bool {
	return u.ID == 0
}

type ExternalUser struct {
	User
	Provider   string
	ExternalId string
}
