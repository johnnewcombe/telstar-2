package api

//--
// Data model objects and persistence mocks:
//--

// User data model
type User struct {
	ID       int    `json:"user-id" :"id"`        // ten digit numeric code e.g. 1000000000 - 9999999999
	Password string `json:"password" :"password"` // four character numeric pin without leading zeros e.g. 1000 - 9999
}
