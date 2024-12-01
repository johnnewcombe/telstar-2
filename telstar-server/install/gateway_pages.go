package install

import (
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/dal"
	"strings"
	"time"
)

func CreateGatewayPages(settings config.Config) error {
	var (
		err error
	)

	if globals.Debug {
		defer logger.TimeTrack(time.Now(), "CreateGatewayPages")
	}

	if err = create6Page(settings); err != nil {
		return err
	}
	if err = create61Page(settings); err != nil {
		return err
	}
	if err = create62Page(settings); err != nil {
		return err
	}
	if err = create63Page(settings); err != nil {
		return err
	}
	if err = create64Page(settings); err != nil {
		return err
	}

	return nil
}

func create6Page(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 6
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.FrameType = globals.FRAME_TYPE_INFORMATION
	frame.Content.Data = "http://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBthQIECAbUizKdSDSDR8PTL3w-UFPLy7aceXmgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgKpEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIA1SLMp1INJBw5b-2nJl5oM-Hpl74fPNB03oN_TRl5IECBAgDc_PPpl281iDpoy88qDDyyoMmXph07MuRBiy7N_dcgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIA7FBDnzJk-nBmIIcGtFQQYlaLOqVaUUNQ2YfKDpoyoGbZggDIEHDfp3dEHbLy56d-5BvzIOmjTzQZsO3f15oM-HblXIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIA7IPUi1KUGHImT6U-bFqSIM6KGrafSDp54ZUGjrt39eSBAgDIEGblv2oOmjKg2b92fLz6IOXXdu07s6DtpyZd_TL4QIECAMgQc8vLtpx5VyBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgDs0E6xUizA1bTky7-mXwgzb-SDpoyoJ1hBT4ZcfTl12rkCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIEDRBUixY0GwGqaMqDfu2ad2VB0y7MvTL46IOeXl2048q5AgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQICyRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 0}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)
}

func create61Page(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 61
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Title.Data = "[b][m-][g[   Telstar Gateway]][b][h-]\r\n"
	frame.Title.Type = "markup"
	frame.FrameType = globals.FRAME_TYPE_GATEWAY
	frame.Content.Data = "[G]You are about to connect to Cave\r\n" +
		"[G]Adventure,[F]make sure you have a map!\r\n\r\n" +
		"[W]Exiting the Gateway service will\r\n" +
		"[W]return you to Telstar.\r\n\r\n" +
		"[G]The DLE character (Ctrl P), where\r\n" +
		"[G]keyboards support it, can also be\r\n" +
		"[G]used to return to Telstar."
	frame.Content.Type = "markup"
	frame.RoutingTable = []int{610, 611, 612, 613, 614, 615, 616, 617, 618, 619, 6}
	frame.Connection.Address = "localhost"
	frame.Connection.Port = 6505
	frame.Connection.Mode = globals.ConnectionModeFullDuplex

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)
}

func create62Page(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 62
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Title.Data = "[b][m-][g[   Telstar Gateway]][b][h-]\r\n"
	frame.Title.Type = "markup"
	frame.FrameType = globals.FRAME_TYPE_GATEWAY
	frame.Content.Data = "[G]You are about to connect to[W]CCL4.\r\n" +
		"[G]Prepare yourself for Viz-type humour.\r\n\r\n\r\n" +
		"[W]Exiting the Gateway service will\r\n" +
		"[W]return you to Telstar.\r\n\r\n" +
		"[G]The DLE character (Ctrl P), where\r\n" +
		"[G]keyboards support it, can also be\r\n" +
		"[G]used to return to Telstar."
	frame.Content.Type = "markup"
	frame.RoutingTable = []int{620, 621, 622, 623, 624, 625, 626, 627, 628, 629, 6}
	frame.Connection.Address = "fish.ccl4.org"
	frame.Connection.Port = 23
	frame.Connection.Mode = globals.ConnectionModeViewdata

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create63Page(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 63
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Title.Data = "[b][m-][g[   Telstar Gateway]][b][h-]\r\n"
	frame.Title.Type = "markup"
	frame.FrameType = globals.FRAME_TYPE_GATEWAY
	frame.Content.Data = "[G]You are about to connect to[W]NxTel.\r\n\r\n\r\n" +
		"[W]Exiting the Gateway service will\r\n" +
		"[W]return you to Telstar.\r\n\r\n" +
		"[G]The DLE character (Ctrl P), where\r\n" +
		"[G]keyboards support it, can also be\r\n" +
		"[G]used to return to Telstar."
	frame.Content.Type = "markup"
	frame.RoutingTable = []int{630, 631, 632, 633, 634, 635, 636, 637, 638, 639, 6}
	frame.Connection.Address = "nx.nxtel.org"
	frame.Connection.Port = 23280
	frame.Connection.Mode = globals.ConnectionModeViewdata

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create64Page(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 64
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Title.Data = "[b][m-][g[   Telstar Gateway]][b][h-]\r\n"
	frame.Title.Type = "markup"
	frame.FrameType = globals.FRAME_TYPE_GATEWAY
	frame.Content.Data = "[G]You are about to connect to[W]TEEFAX.\r\n\r\n\r\n" +
		"[W]Exiting the Gateway service will\r\n" +
		"[W]return you to Telstar.\r\n\r\n" +
		"[G]The DLE character (Ctrl P), where\r\n" +
		"[G]keyboards support it, can also be\r\n" +
		"[G]used to return to Telstar."
	frame.Content.Type = "markup"
	frame.RoutingTable = []int{640, 641, 642, 643, 644, 645, 646, 647, 648, 649, 6}
	frame.Connection.Address = "pegasus.matrixnetwork.co.uk"
	frame.Connection.Port = 6502
	frame.Connection.Mode = globals.ConnectionModeViewdata

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}
