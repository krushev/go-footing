package services

type Role string
const (
	Admin = "admin"
	User = "user"
)

func IsValid(roles []string) bool {
	for _, r := range roles {
		if !Role(r).IsValid() {
			return false
		}
	}
	return true
}

func (r Role) IsValid() bool {
	switch r {
	case Admin, User:
		return true
	}
	return false
}
