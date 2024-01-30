package api

type ServerStatus struct {
	Version SystemVersion `json:"version" bson:"version"`
}

// Version of Telstar is determined by changes to the API, therefore Version is defined within the API only
type SystemVersion struct {
	Major int    `json:"major" bson:"major"`
	Minor int    `json:"minor" bson:"minor"`
	Patch int    `json:"patch" bson:"patch"`
	Info  string `json:"info" bson:"msg"`
}

var systemStatus = ServerStatus{
	version,
}
