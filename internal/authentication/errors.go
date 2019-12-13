package authentication

type NotAuthenticated struct{}

func (e NotAuthenticated) Error() string {
	return "user is not authenticated"
}
