package user

import (
	"github.com/manifest-destiny/api"
)

// FindByGoogleID returns a user from the database by user's google ID.
func FindByGoogleID(conn *api.DB, u *User, gid string) error {
	return conn.DB.Get(u, `SELECT * FROM app_user WHERE google_id=$1`, gid)
}

// PersistUser adds a new user to the database.
func PersistUser(conn *api.DB, u *User) error {
	_, err := conn.DB.NamedExec(`INSERT INTO app_user (google_id, email, full_name, alias, picture, show_picture, locale, country) VALUES (:google_id, :email, :full_name, :alias, :picture, :show_picture, :locale, :country);`, u)
	return err
}

// UpdateUser updates an existing user to the database.
func UpdateUser(conn *api.DB, u *User) error {
	_, err := conn.DB.NamedExec(`UPDATE app_user SET (google_id, email, full_name, alias, picture, show_picture, locale, country) = (:google_id, :email, :full_name, :alias, :picture, :show_picture, :locale, :country) WHERE id = :id;`, u)
	return err
}
