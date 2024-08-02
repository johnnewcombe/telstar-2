package api

import (
	"bitbucket.org/johnnewcombe/telstar-library/globals"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar/config"
	"bitbucket.org/johnnewcombe/telstar/dal"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net/http"
	"strings"
)

// TODO put error messages here
const (
	ERR_NON_API_ACCOUNT  = "this accounts does not have API access"
	ERR_INVALID_HASH     = "invalid user"
	ERR_INVALID_PAGEID   = "invalid pageId"
	ERR_PAGEID_NOT_FOUND = "frame not found"
	ERR_USER_SCOPE       = "user does not have sufficient scope to perform this task"
	ERR_INVALID_USERID   = "invalid user ID"
	ERR_USER_NOT_FOUND   = "user not found"
	MSG_USER_DELETED     = "user deleted"
	MSG_USER_UPDATED     = "user updated"
	MSG_FRAME_DELETED    = "frame deleted"
	MSG_FRAMES_PURGED    = "%d frames deleted"
	MSG_FRAME_UPDATED    = "frame updated"
	MSG_LOGIN_SUCCESS    = "login successful"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	var (
		output strings.Builder
	)
	switch r.Method {
	case "GET":
		// TODO: Consider an Http Template for this maybe
		w.WriteHeader(http.StatusOK)
		output.WriteString("<!DOCTYPE html>")
		output.WriteString("<head>")
		output.WriteString("<title>Welcome to the Telstar Viewdata API!</title>")
		output.WriteString("<style>body {width: 50em;margin: 0 auto;font-family: Tahoma, Verdana, Arial, sans-serif;}</style>")
		output.WriteString("</head>")
		output.WriteString("<body>")
		output.WriteString("<h1>Welcome to the Telstar Viewdata API!</h1>")
		output.WriteString("<p>If you see this page, the the Telstar API web server is successfully installed and\nworking.</p>")
		output.WriteString("<p>For online documentation and support please refer to <a href=\"https://glasstty.com/\" target=\"_blank\">glasstty.com</a>.</p>")
		output.WriteString("<p><em>Thank you for using Telstar.</em></p>")
		output.WriteString("</body>")

		w.Write([]byte(output.String()))

	default:
		w.WriteHeader(http.StatusNotFound)

	}
}

func putLogin(w http.ResponseWriter, r *http.Request) {

	var (
		userDbEntry types.User
		userPutData types.User
		jwtToken    string
		err         error
		user        types.User
		settings    config.Config
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get the username and password from the request
	if err = json.NewDecoder(r.Body).Decode(&userPutData); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	// FIXME does this make sense for login?
	if user, err = dal.GetUser(settings.Database.Connection, userPutData.UserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// get the corresponding user from the database, for this we need to be authorised which won't be the case
	// util we have checked the users credentials, however, a user can always get details of itself so we use
	// the userId that was sent in the login data and pass this as the authenticated user
	// FIXME is the above correct or should we user dal.GetUser instead
	if userDbEntry, err = dal.GetUserByUser(settings.Database.Connection, userPutData.UserId, user); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// non api accounts cant use the API, this includes guests
	if !userDbEntry.ApiAccess {
		render.Render(w, r, ErrUnauthorizedRequest(errors.New(ERR_NON_API_ACCOUNT)))
		return
	}

	// check any hash we have (if uer doesn't exist then hash will be empty string
	if !dal.CheckPasswordHash(userPutData.Password, userDbEntry.Password) {
		render.Render(w, r, ErrUnauthorizedRequest(errors.New(ERR_INVALID_HASH)))
		return
	}

	// if we get to here then we are authorised
	userDbEntry.Authenticated = true

	// create the jwt using the
	if jwtToken, err = createJwtToken(userPutData.UserId, 15); err != nil {

		render.Render(w, r, ErrServerRequest(err))
		return
	}

	cookie := &http.Cookie{
		Name:  "token",
		Value: jwtToken,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	// send a status 200 back in the response body
	// Note that the HttpStatusCode gets set in the HttpResponse
	// the msg appears in the response body, the two are not related
	render.Render(w, r, HttpResult(200, MSG_LOGIN_SUCCESS))

}

func getFrame(w http.ResponseWriter, r *http.Request) {

	var (
		jsonFrame  []byte
		frame      types.Frame
		user       types.User
		err        error
		pageNumber int
		frameId    string
		authUserId string
		settings   config.Config
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if user, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	if pageId := chi.URLParam(r, "pageId"); len(pageId) > 0 {

		if pageNumber, frameId, err = utils.ConvertPageIdToPID(pageId); err != nil {
			render.Render(w, r, ErrTeapotRequest(err))
		}

		primary := strings.ToLower(r.URL.Query().Get("db")) == globals.DBPRIMARY
		connection := settings.Database.Connection

		if frame, err = dal.GetFrameByUser(connection, pageNumber, frameId, primary, user); err != nil {
			render.Render(w, r, ErrNotFoundRequest(err))
			return
		}

	} else {
		render.Render(w, r, ErrNotFoundRequest(errors.New(ERR_INVALID_PAGEID)))
		return
	}

	if jsonFrame, err = json.MarshalIndent(frame, "", "    "); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}
	w.Write(jsonFrame)
}

func getFrames(w http.ResponseWriter, r *http.Request) {

	var (
		frames     []types.Frame
		user       types.User
		err        error
		jsonFrames []byte
		settings   config.Config
		authUserId string
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// gget authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if user, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	primary := strings.ToLower(r.URL.Query().Get("db")) == globals.DBPRIMARY
	connection := settings.Database.Connection

	if frames, err = dal.GetFramesByUser(connection, primary, user); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	if jsonFrames, err = json.MarshalIndent(frames, "", "    "); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}
	w.Write(jsonFrames)
}

func publishFrame(w http.ResponseWriter, r *http.Request) {

	var (
		user       types.User
		err        error
		pageNumber int
		frameId    string
		authUserId string
		settings   config.Config
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		_ = render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if user, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	if pageId := chi.URLParam(r, "pageId"); len(pageId) > 0 {

		if pageNumber, frameId, err = utils.ConvertPageIdToPID(pageId); err != nil {
			render.Render(w, r, ErrTeapotRequest(err))
		}

		connection := settings.Database.Connection

		if err = dal.PublishFrameByUser(connection, pageNumber, frameId, user); err != nil {
			render.Render(w, r, ErrNotFoundRequest(err))
			return
		}

	} else {
		render.Render(w, r, ErrNotFoundRequest(errors.New(ERR_INVALID_PAGEID)))
		return
	}
	/*
		if jsonFrame, err = json.MarshalIndent(frame, "", "    "); err != nil {
			render.Render(w, r, ErrServerRequest(err))
			return
		}
		w.Write(jsonFrame)

	*/
}

// getStatus returns the API Version number and other info
func getStatus(w http.ResponseWriter, r *http.Request) {

	var (
		jsonVersion []byte
	)
	jsonVersion = []byte(version)
	w.Write(jsonVersion)
}

func addFrame(w http.ResponseWriter, r *http.Request) {
	// update (PUT) is used to add and update a frame, POST is not implemented
	render.Render(w, r, ErrServerRequest(errors.New("not implemented")))
}

func updateFrame(w http.ResponseWriter, r *http.Request) {

	var (
		frame      types.Frame
		user       types.User
		err        error
		settings   config.Config
		authUserId string
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if user, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	// get the frame from the request
	if err = json.NewDecoder(r.Body).Decode(&frame); err != nil {
		render.Render(w, r, ErrTeapotRequest(err))
		return
	}

	if !utils.IsValidRoutingTable(frame.RoutingTable) {
		frame.RoutingTable = utils.CreateDefaultRoutingTable(frame.PID.PageNumber)
	}

	primary := strings.ToLower(r.URL.Query().Get("db")) == globals.DBPRIMARY
	connection := settings.Database.Connection

	if err = dal.InsertOrReplaceFrameByUser(connection, frame, primary, user); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// send a status 200 back in the response body
	// Note that the HttpStatusCode gets set in the HttpResponse
	// the msg appears in the response body, the two are not related
	render.Render(w, r, HttpResult(200, MSG_FRAME_UPDATED))
}

func deleteFrame(w http.ResponseWriter, r *http.Request) {

	var (
		pageNumber   int
		frameId      string
		deletedCount int64
		err          error
		settings     config.Config
		authUserId   string
		user         types.User
		msg          string
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if user, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	purge := strings.ToLower(r.URL.Query().Get("purge")) == "true"

	if pageId := chi.URLParam(r, "pageId"); len(pageId) > 0 {

		if pageNumber, frameId, err = utils.ConvertPageIdToPID(pageId); err != nil {
			render.Render(w, r, ErrTeapotRequest(err))
		}

		primary := strings.ToLower(r.URL.Query().Get("db")) == globals.DBPRIMARY
		connection := settings.Database.Connection

		if purge {
			if deletedCount, err = dal.PurgeFramesByUser(connection, pageNumber, frameId, primary, user); err != nil {
				render.Render(w, r, ErrServerRequest(err))
				return
			}
			msg = fmt.Sprintf(MSG_FRAMES_PURGED, deletedCount)

		} else {
			if deletedCount, err = dal.DeleteFrameByUser(connection, pageNumber, frameId, primary, user); err != nil {
				render.Render(w, r, ErrServerRequest(err))
				return
			}
			msg = MSG_FRAME_DELETED
		}

		if deletedCount == 0 {
			render.Render(w, r, HttpResult(404, ERR_PAGEID_NOT_FOUND))
			return
		}

	} else {
		render.Render(w, r, ErrNotFoundRequest(errors.New(ERR_INVALID_PAGEID)))
		return
	}

	render.Render(w, r, HttpResult(200, msg))

}

func addUser(w http.ResponseWriter, r *http.Request) {
	// update (PUT) is used to add and update a user, POST is not implemented
	render.Render(w, r, ErrServerRequest(errors.New("not implemented")))
}

func getUser(w http.ResponseWriter, r *http.Request) {

	var (
		jsonUser   []byte
		userId     string
		user       types.User
		authUser   types.User
		err        error
		settings   config.Config
		authUserId string
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if authUser, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	if userId = chi.URLParam(r, "userId"); len(userId) > 0 {

		//if userId, err = strconv.Atoi(userIdS); err != nil {
		//	render.Render(w, r, ErrTeapotRequest(err))
		//	return
		//}
		// check userid,

		if !utils.IsValidUserId(userId) {
			render.Render(w, r, ErrTeapotRequest(err))
			return
		}

		connection := settings.Database.Connection

		// user cannot get the json of a frame that is out of scope
		// FIXME WHAT PERMISSIONS SHOULD A USER HAVE
		//if !dal.IsUserAdmin(connection, authUserId) {
		//	render.Render(w, r, ErrUnauthorizedRequest(errors.New(ERR_USER_SCOPE)))
		//}

		if user, err = dal.GetUserByUser(connection, userId, authUser); err != nil {
			render.Render(w, r, ErrNotFoundRequest(err))
			return
		}

		if jsonUser, err = json.MarshalIndent(user, "", "    "); err != nil {
			render.Render(w, r, ErrServerRequest(err))
			return
		}

	} else {
		render.Render(w, r, ErrNotFoundRequest(errors.New(ERR_INVALID_USERID)))
		return
	}

	w.Write(jsonUser)

}

/*
	func getUsers(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, ErrServerRequest(errors.New("not implemented)")))
	}
*/
func updateUser(w http.ResponseWriter, r *http.Request) {

	var (
		user       types.User
		authUser   types.User
		err        error
		settings   config.Config
		authUserId string
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// get authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if authUser, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	user.Authenticated = true

	// get the user-id, password etc. from the request
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Render(w, r, ErrTeapotRequest(err))
		return
	}

	connection := settings.Database.Connection

	// check logged on user for admin status
	if !dal.IsUserAdmin(connection, authUserId) {
		// FIXME A user should be able to update own name and password only, nothing else
		render.Render(w, r, ErrUnauthorizedRequest(errors.New(ERR_USER_SCOPE)))
	}

	if err = dal.InsertOrReplaceUserByUser(connection, user, authUser); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// send a status 200 back in the response body
	// Note that the HttpStatusCode gets set in the HttpResponse
	// the msg appears in the response body, the two are not related
	render.Render(w, r, HttpResult(200, MSG_USER_UPDATED))

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	var (
		userId string
		//		userId  int
		deletedCount int64
		err          error
		authUser     types.User
		settings     config.Config
		authUserId   string
	)

	// get settings from from context
	if settings, err = getSettingsFromContext(r); err != nil {
		render.Render(w, r, ErrServerRequest(err))
		return
	}

	// gget authorised user from from context
	if authUserId, err = getAuthUserIDFromContext(r); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
		return
	}

	// get the current logged on user
	if authUser, err = dal.GetUser(settings.Database.Connection, authUserId); err != nil {
		render.Render(w, r, ErrUnauthorizedRequest(err))
	}

	// if we have got this far them we have been authenticated
	authUser.Authenticated = true

	if userId = chi.URLParam(r, "userId"); len(userId) > 0 {

		//if userId, err = strconv.Atoi(userIdS); err != nil {
		//	render.Render(w, r, ErrTeapotRequest(err))
		//	return
		//}

		// check userid,
		if !utils.IsValidUserId(userId) {
			render.Render(w, r, ErrTeapotRequest(err))
			return
		}

		connection := settings.Database.Connection

		// get user and pass to DAL, DAL will enforce permissions
		//if !dal.IsUserAdmin(connection, authUserId) {
		//	render.Render(w, r, ErrUnauthorizedRequest(errors.New(ERR_USER_SCOPE)))
		//}

		if deletedCount, err = dal.DeleteUserByUser(connection, userId, authUser); err != nil {
			render.Render(w, r, ErrServerRequest(err))
			return
		}
		if deletedCount == 0 {
			render.Render(w, r, HttpResult(404, ERR_USER_NOT_FOUND))
			return
		}

	} else {
		render.Render(w, r, ErrNotFoundRequest(errors.New(ERR_INVALID_USERID)))
		return
	}
	// send a status 200 back in the response body
	// Note that the HttpStatusCode gets set in the HttpResponse
	// the msg appears in the response body, the two are not related
	render.Render(w, r, HttpResult(200, MSG_USER_DELETED))
}

// getAuthUserIDFromContext gets the authorised user from context
func getAuthUserIDFromContext(r *http.Request) (string, error) {
	var (
		claims map[string]interface{}
		err    error
	)
	if _, claims, err = jwtauth.FromContext(r.Context()); err != nil {
		return "", err
	}
	userId := claims["user-id"].(string)

	return userId, nil
}

// getSettingsFromContext get settings from context, if it doesn't exist we get a panic
// but the chi middleware.Recoverer() will deal with this and return a 500 server error.
func getSettingsFromContext(r *http.Request) (config.Config, error) {
	ctxData := r.Context().Value("ctx-data").(*ContextData)
	return ctxData.Settings, nil
}

func HttpResult(httpResponseCode int, msg string) render.Renderer {
	return &types.ApiResponse{
		// the HttpStatusCode gets set in the HttpResponse
		// the result Text appears in the body
		HTTPStatusCode: httpResponseCode,
		ResultText:     msg,
	}
}
