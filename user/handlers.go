package user

import (
	"github.com/emicklei/go-restful"
	"github.com/manifest-destiny/api"
)

// Resource for user handler methods.
type Resource struct {
	*api.DB
	*api.TokenValidator
}

// RegisterContainer Defines user endpoints
func RegisterContainer(container *restful.Container, r *Resource) {
	ws := new(restful.WebService)

	ws.Param(ws.HeaderParameter("Authorization", "User's access token").
		DataType("string"))

	ws.ApiVersion("1.0").
		Path("/").
		Doc("Manage users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/account").To(r.account).
		Doc("get my details").
		Operation("account").
		Writes(User{}))

	container.Add(ws)
}

// account Returns user's details. If user has never signed in, then user
// is persisted to database first. If already registered then Google-managed
// fields are updated (if they have changed).
func (r *Resource) account(request *restful.Request, response *restful.Response) {
	u, err := Authenticate(request, r.TokenValidator, r.DB)
	if err != nil {
		api.WriteError(response, err)
		return
	}

	response.WriteEntity(u)
}

// Authenticate checks a user's bearer token and returns a User. User is
// also persisted or updated in the database.
func Authenticate(req *restful.Request, v *api.TokenValidator, conn *api.DB) (u *User, err error) {
	u = &User{}

	bearerToken, err := api.BearerToken(req)
	if err != nil {
		return
	}

	// Verify the token signature
	c := &api.GoogleIdentityClaims{
		Claims: &api.Claims{Validator: v.Claims},
	}
	err = v.Signatures.Verify(bearerToken, c)
	if err != nil {
		err = api.InvalidTokenErr
		return
	}
	// Validate token claims
	err = c.Valid()
	if err != nil {
		err = api.InvalidTokenErr
		return
	}

	// Find a user by google ID. If we cannot, then persist
	// the user, else if the google-managed user fields have
	// changed then update the user.
	err = FindByGoogleID(conn, u, c.Subject)
	if err != nil {
		// Persist a new user before returning
		defer func() {
			// Set application-managed fields
			u.Alias = c.GivenName
			u.ShowPicture = true
			u.setCountryFromLocale()
			err = PersistUser(conn, u)
		}()
	} else if googleAccountModified(u, c) {
		// Update the user before returning
		defer func() { err = UpdateUser(conn, u) }()
	}

	// Set user fields
	u.GoogleID = c.Subject
	u.Email = c.Email
	u.Name = c.Name
	u.Locale = c.Locale
	u.Picture = c.Picture

	return
}

// googleAccountModified checks if google-managed fields equal fields stored
// by application.
func googleAccountModified(u *User, c *api.GoogleIdentityClaims) bool {
	if u.Email != c.Email ||
		u.Name != c.Name ||
		u.Locale != c.Locale ||
		u.Picture != c.Picture {
		return true
	}

	return false
}
