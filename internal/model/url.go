package model

import "fmt"

type URL struct {
	Long string
	Base string
	ID   string
}

func (u *URL) Short() string {
	return fmt.Sprintf("%s/%s", u.Base, u.ID)
}
