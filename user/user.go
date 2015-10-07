package user

import (
	"strings"
)

// User User struct with db and json tags.
type User struct {
	ID          int64  `db:"id" json:"-"`
	GoogleID    string `db:"google_id" json:"-"`
	GoogleHash  string `db:"google_fields_hash" json:"-"`
	Name        string `db:"full_name"`
	Email       string `db:"email"`
	ShowPicture bool   `db:"show_picture"`
	Picture     string `db:"picture"`
	Alias       string `db:"alias"`
	Locale      string `db:"locale"`
	Country     string `db:"country"`
}

func (u *User) setCountryFromLocale() {
	locSplit := strings.SplitN(u.Locale, "-", 2)
	if len(locSplit) == 2 {
		u.Country = strings.ToLower(locSplit[1])
	} else {
		u.Country = ""
	}
}
