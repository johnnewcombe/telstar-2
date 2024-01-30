package types

import (
	"strconv"
	"strings"
)

/*
TODO: Implement Users
TODO
====

Users are Guest if they are user globals.GUEST_USER or if they are not authenticated

// FIXME this doesn't really make sense
Do users not authenticated revert to a guest account? or simply behave are Guest users?
When a user is reverted do we transfer session to guest?
Do we need a Guest account, i.e. if current account not set assume guest.
If a user doesn't log on User is nil i.e. unauthenticated reader.


Dictionary of page numbers that determine if login is required? or by individual page
redirect to login if not authorised
Do we need admin status ? admin status derived from base page number or bool?
admin mailbox could be used for significant events rather than an admin account.

Add roles support for the future, editor, admin, reader etc.
users in roles or roles attached to user? default role = readers.
roles are hierarchical admin, editors, readers? or just a list attached to a user.
api users must be restricted to base page and below
admins are still restricted to base page so an admin of page 800 and below can be created.

Completed
=========
Same users as API but some with ApiAccess flag, base page is ignored if no api access.
Api handler prevents access to non-api accounts
No requirement to login.
User gets added to session when handler starts
session only on server, api is stateless
Every user whether logged in or not gets session id
*/

// User struct contains the user details note that password is hashed as it is inserted in the database
// also all fields are strings even though credentials are numbers, this prevents 0 (admin) being set accidentally
type User struct {
	UserId        string `json:"user-id" bson:"user-id"`
	Password      string `json:"password" bson:"password"` // four character numeric pin without leading zeros e.g. 1000 - 9999
	Name          string `json:"name" bson:"name"`         // 24 char name (check length to be used)
	BasePage      int `json:"base-page" bson:"base-page"`
	ApiAccess     bool   `json:"api-access" bson:"api-access"`
	Admin         bool   `json:"admin" bson:"admin"`
	Editor        bool   `json:"editor" bson:"editor"`
	Authenticated bool   // no json for this as it must default to false
}

func (u *User) UserIdAsInteger() int64 {
	//FIXME ?? is this needed?
	return 0
}

func (u *User) IsGuest() bool {
	// No user Id
	return len(u.UserId) ==0
}

func (u *User) IsInScope(pageNumber int) bool {

	// unauthenticated users are never in scope.
	if u.BasePage == 0 {
		return true
	}

	page := strconv.Itoa(pageNumber)
	if strings.HasPrefix(page, strconv.Itoa(u.BasePage)) {
		return true
	}
	return false
}

