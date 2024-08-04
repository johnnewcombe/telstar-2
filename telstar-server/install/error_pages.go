package install

import (
	"fmt"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/dal"
	"strings"
)

func CreateErrorPages(settings config.Config) error {
	var (
		err error
	)
	if err = create9901Page(settings); err != nil {
		return err
	}
	if err = create9902Page(settings); err != nil {
		return err
	}
	if err = create9903Page(settings); err != nil {
		return err
	}
	return nil
}

func create9901Page(settings config.Config) error {

	var (
		err     error
		frame   types.Frame
		primary bool
	)

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	frame.PID.PageNumber = 9901
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Content.Type = "markup"

	var sb strings.Builder
	sb.WriteString("[D]Unexpected Error\r\n\n")
	sb.WriteString("[r][l-]\r\n\n")
	sb.WriteString("[Y]The error message was:\r\n\n")
	sb.WriteString("[W][ERROR]\r\n\n\n\n")
	sb.WriteString("[Y]It would help if this issue could be\r\n")
	sb.WriteString("[Y]reported on the[W]TELSTAR[Y]group at:\r\n\n")
	sb.WriteString("[W]    https://groups.io/g/telstar\r\n\n")
	sb.WriteString("[r][l-]\r\n")
	//sb.WriteString("[G]   PRESS[W]1[G]TO LOG IN,[W]" + string(globals.HASH) + "[G]TO CONTINUE")

	frame.Content.Data = sb.String()
	frame.FrameType = globals.FRAME_TYPE_INFORMATION
	frame.Cursor = false
	frame.RoutingTable = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	frame.NavMessage = fmt.Sprintf("[R][n][Y]PRESS[W]%s[Y]TO CONTINUE,[W]", string(globals.HASH))
	frame.NavMessageNotFound = NavNotFoundBlue

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

	return err
}
func create9902Page(settings config.Config) error {

	var (
		err     error
		frame   types.Frame
		primary bool
	)

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	frame.PID.PageNumber = 9902
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Content.Type = "edittf"

	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBANp2adSLNQRaVKfSQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQFEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgBIIKDn559Mu1Bl5ct_JBvx4-vLLkWIOmjTzQaeaDJ1yoECAEg6b0GHcgy-OmXluw7EHDlvz8sO1Bmw6dnXllXIKmjKgQIASDNh07OvLKg0YeaDFly7kHLLw38umXIg6b0HTRlQIECBAgBIOfnn0y7UGHJt07tPPpyw9N_JcgobMuHnlQYcePLw6IECAEg39eSDDw37N-fTl5oM2_kg6aMqDTux793bLu05d2PKuQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAEgocsvPmHvgum9Byy9OvLcg6b0HTRlQbcOncg25d3VcgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAUSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.FrameType = globals.FRAME_TYPE_INITIAL
	frame.Cursor = false
	frame.RoutingTable = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	frame.NavMessage = fmt.Sprintf("[B][n][Y]PRESS[W]%s[Y]TO CONTINUE,[W]", string(globals.HASH))
	frame.NavMessageNotFound = NavNotFoundBlue

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

	return err
}

func create9903Page(settings config.Config) error {

	var (
		err     error
		frame   types.Frame
		primary bool
	)

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	frame.PID.PageNumber = 9903
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Content.Type = "edittf"

	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQDY8GpFrwbKCnFpVpMOKgk06dWKgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIEBRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAEgqaMqDPh6Ze-Hygx7927Lj6ad-5Bow80GXdky5FyBAgQIASCpo080Gnmgy-OGXH0y5EGnMgwoOeXl2048qDvh5oECBAgBIMu7JlyIOvPTuzoMKCJMioFEPpy2IKClBj0YeWHH0y8liAEg39NGXl3088qDpo080G3Lh3c0HTRh6IOmjLyyoECBAgQIASDbh8oMWVBh3INPPn1yoO-npoQdNGVBHw9MvfD5Qb-SBAgBII-Hpl74fKDnl5dtOPKuQUNmXDzyoMOPHl4dEG_ryQIECAEgw8N-zfn05eaDNv5IOmjKg07se_d2y7tOXdjyrkCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgKJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIASChyy8-Ye-C6b0HLL068tyDpvQdNGVBtw6dyDbl3dVyBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAokSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.FrameType = globals.FRAME_TYPE_INITIAL
	frame.Cursor = false
	frame.RoutingTable = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	frame.NavMessage = fmt.Sprintf("[B][n][Y]PRESS[W]%s[Y]TO CONTINUE,[W]", string(globals.HASH))
	frame.NavMessageNotFound = NavNotFoundBlue

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

	return err
}
