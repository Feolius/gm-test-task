package repository

type User struct {
	Id       int
	Login    string
	Password string
}

func (u *User) Empty() bool {
	return u.Id == 0
}
