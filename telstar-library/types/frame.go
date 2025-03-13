package types

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
)

type Frame struct {
	//ID interface{} `bson:"_id,omitempty"` // this was used by some on stack overflow
	// FIXME Need to include Object ID for internal use i.e. data cleaning etc
	//       but remove it before sending json to caller
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	PID          Pid                `json:"pid" bson:"pid"`
	Visible      bool               `json:"visible" bson:"visible"`
	Carousel     bool               `json:"carousel" bson:"carousel"`
	HeaderText   string             `json:"header-text" bson:"header-text"`
	Cost         int                `json:"cost" bson:"cost"`
	DisableClear bool               `json:"disable-clear" bson:"disable-clear"`
	FrameType    string             `json:"frame-type" bson:"frame-type"`
	Redirect     Pid                `json:"redirect" bson:"redirect"`
	Content      Content            `json:"content" bson:"content"`
	Footer       Title              `json:"footer" bson:"footer"`
	Title        Title              `json:"title" bson:"title"`
	RoutingTable []int              `json:"routing-table" bson:"routing-table"`
	Connection   Connection         `json:"connection" bson:"connection"`
	AuthorIdOld  string             `json:"author-id" bson:"author-id"`
	ResponseData ResponseData       `json:"response-data" bson:"response-data"`
	//CacheId          string       `json:"cache-id" bson:"cache-id"`
	// FIXME The Cursor field doesn't seem to be being used at all, either remove it or find a use for it
	//  the cursor in the nav field is added using markup
	Cursor bool `json:"cursor" bson:"cursor"`
	// a transient page will not be stored in the database e.g. results of response pages etc.
	//TransientPage bool `json:"transient-page" bson:"transient-page"`
	// these will override the defaults specified in settings at the page level
	NavMessage         string `json:"navmessage-select" bson:"navmessage-select"`
	NavMessageNotFound string `json:"navmessage-notfound" bson:"navmessage-notfound"`
}
type Pid struct {
	PageNumber int    `json:"page-no" bson:"page-no"`
	FrameId    string `json:"frame-id" bson:"frame-id"`
}
type Content struct {
	//Data interface{} `json:"data" bson:"data"`
	Data string `json:"data" bson:"data"`
	Type string `json:"type" bson:"type"`
}
type Title struct {
	Data      string   `json:"data" bson:"data"`
	Type      string   `json:"type" bson:"type"`
	MergeData []string `json:"merge-data" bson:"merge-data"`
}
type Connection struct {
	Address          string   `json:"address" bson:"address"`
	Mode             string   `json:"mode" bson:"mode"`
	Port             int      `json:"port" bson:"port"`
	remoteConnection net.Conn // used by methods only
}
type ResponseData struct {
	Fields []ResponseField `json:"response-fields" bson:"responses"`
	Action ResponseAction  `json:"response-action" bson:"response-action"`
}
type ResponseField struct {
	//	Label    string `json:"label" bson:"label"`
	VPos     int    `json:"vpos" bson:"vpos"`
	HPos     int    `json:"hpos" bson:"hpos"`
	Required bool   `json:"required" bson:"required"`
	Length   int    `json:"length" bson:"length"`
	Type     string `json:"type" bson:"type"`
	// FIXME: this is specified as 'auto_submit' on many response tmp
	// in which case this will always return false, this needs to be fixed
	// in the database for v2
	AutoSubmit bool `json:"auto-submit" bson:"auto-submit"`
	Password   bool `json:"password" bson:"password"`

	// these two removed for v2.0 as theses are results not definitions
	//Value      string `json:"value" bson:"value"`
	//Valid      bool   `json:"valid" bson:"valid"`
}
type ResponseAction struct {
	Exec            string   `json:"exec" bson:"exec"`
	Args            []string `json:"args" bson:"args"`
	PostActionFrame Pid      `json:"post-action-frame" bson:"post-action-frame"`
	PostCancelFrame Pid      `json:"post-cancel-frame" bson:"post-cancel-frame"`
}

func (f *Frame) IsValid() bool {
	return len(f.PID.FrameId) == 1
}

func (f *Frame) GetPageId() string {
	return fmt.Sprint(f.PID.PageNumber) + f.PID.FrameId
}

func (f *Frame) GetRedirectPageId() string {
	return fmt.Sprint(f.Redirect.PageNumber) + f.Redirect.FrameId
}

// GetZeroPageRoute returns the current page appended with 0. For example if the
// current page is 293, the 'zero' page route would be page 2930.
func (f *Frame) GetZeroPageRoute() int {
	nextPage := f.PID.PageNumber * 10
	if nextPage > 999999999 {
		nextPage = 0
	}
	return nextPage
}

func (c *Connection) IsValid() bool {
	if len(c.Address) == 0 || c.Port == 0 || len(c.Mode) == 0 {
		return false
	}
	return true
}

func (c *Connection) GetUrl() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

func (c *Connection) GetRemoteConnection() net.Conn {
	return c.remoteConnection
}

func (c *Connection) SetRemoteConnection(value net.Conn) {
	c.remoteConnection = value
}

// CreateRedirect Creates a redirect from the fromPageId to this frame.
func (f *Frame) CreateRedirect(fromPageNumber int, fromFrameId string) Frame {

	// e.g.   {"pid": {"page-no": 0, "frame-id": "a"}, "redirect": {"page-no": 9, "frame-id": "a"},  "visible": True,}

	var rf Frame

	rf.PID.PageNumber = fromPageNumber
	rf.PID.FrameId = fromFrameId
	rf.Redirect.PageNumber = f.PID.PageNumber
	rf.Redirect.FrameId = f.PID.FrameId
	rf.Visible = true

	return rf
}

func (f *Frame) CreateDefaultRoutingTable() []int {

	var (
		routingTable []int
		pageNumber   int
	)

	pageNumber = f.PID.PageNumber

	// default routing
	for i := 0; i < len(routingTable); i++ {
		routingTable = append(routingTable, i+(pageNumber*10))
	}

	// sort out hash route
	pn := pageNumber
	for pn > 999 {
		pn = pn / 10
	}
	routingTable = append(routingTable, pn)

	return routingTable

}

// Load populates a Frame object from json byte data
func (f *Frame) Load(jsonBytes []byte) error {

	if !json.Valid(jsonBytes) {
		return fmt.Errorf("validating frame: invalid json")
	}

	if err := json.Unmarshal(jsonBytes, &f); err != nil {
		return fmt.Errorf("parsing json: invalid")
	}

	return nil
}

// Dump Returns a Json byte representation of the Frame
func (f *Frame) Dump() ([]byte, error) {

	var (
		data []byte
		err  error
	)

	if data, err = json.Marshal(f); err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Pid) String() string {
	return fmt.Sprintf("%d%s", p.PageNumber, p.FrameId)
}
