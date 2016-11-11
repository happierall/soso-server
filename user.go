package soso

import (
	"errors"
	"strconv"
)

type User struct {
	ID     string
	Token  string
	IsAuth bool
}

func (u User) IntID() (int64, error) {
	if len(u.ID) > 0 {
		return strconv.ParseInt(u.ID, 10, 64)
	}
	return -1, errors.New("User.ID does not defined")
}
