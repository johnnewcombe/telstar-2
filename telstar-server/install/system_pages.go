package install

import (
	"github.com/johnnewcombe/telstar-library/convert"
	"github.com/johnnewcombe/telstar-library/globals"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar/config"
	"github.com/johnnewcombe/telstar/dal"
	"strings"
)

const (
	NavSelect        = "[B][n][Y]Select item or[W]*page# : [_+]"
	NavContinue      = "[B][n][Y]Press # to continue :[W]"
	NavNotFoundBlue  = "[B][n][Y]Page not Found :[W]"
	NavNotFoundGreen = " [G]Page not Found :[W]"
)

func CreateSystemPages(settings config.Config) error {

	var (
		err error
	)

	if err = create9a(settings); err != nil {
		return err
	}
	if err = create90a(settings); err != nil {
		return err
	}
	if err = create91a(settings); err != nil {
		return err
	}
	if err = create91b(settings); err != nil {
		return err
	}
	if err = create91c(settings); err != nil {
		return err
	}
	if err = create91d(settings); err != nil {
		return err
	}
	if err = create91e(settings); err != nil {
		return err
	}
	if err = create91f(settings); err != nil {
		return err
	}
	if err = create94a(settings); err != nil {
		return err
	}
	if err = create96a(settings); err != nil {
		return err
	}
	if err = create98a(settings); err != nil {
		return err
	}
	if err = create99a(settings); err != nil {
		return err
	}
	//if err = create990a(settings); err != nil {
	//	return err
	//}

	// error pages
	if err = create9901Page(settings); err != nil {
		return err
	}

	// experiment page
	if err = create101Page(settings); err != nil {
		return err
	}

	return nil
}

func CreateSystemRedirectPages(settings config.Config) error {

	var (
		err error
	)
	if err = create0a(settings); err != nil {
		return err
	}

	return nil
}

func create0a(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
	)
	frame.PID.PageNumber = 0
	frame.PID.FrameId = "a"
	frame.Redirect.PageNumber = 9
	frame.Redirect.FrameId = "a"
	frame.Visible = true

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create9a(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)
	frame.PID.PageNumber = 9
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.FrameType = globals.FRAME_TYPE_INFORMATION
	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgc4UCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQICixYsWLFixYsWLFixYsWLFixYsWLFixYsWLFixYsWLFixYgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAxQIIs6PJnRYtKTOjoKkWnUQUIMeLTQIECBAgQIECBAgQIEDRAgp2adSLNQSZ0aegQIECBAgQIECBAgQIECBAgQIECBAgQNkCCPBqRa8GygpxaVaTDi00CBAgQIECBAgQIECBAgQIECBA4QIIc-dUgw6iCJFqQZMymgQIECBAgQIECBAgQIECBAgQIEDlAgmT46CfOQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQMECCZPjoJ8aMgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQICixYsWLFixYsWLFixYsWLFixYsWLFixYsWLFixYsWLFixYgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAFUQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgBVAMUNMQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIAVQDFDTAlMPUQIECBAgQIECBAgQIECBAgQIECBAgQIECAFUAxQ0wJTD1A0EHSQIECBAgQIECBAgQIECBAgQIECBAgQIECANMCUw9QNBB0kCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAPUDQQdJAgQIECBAgQIECBAgQIECAosWLFixYsWLFixYsQIECAHSQIECAosWLFixYsWLFixYsWIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{90, 91, 92, 93, 94, 95, 6, 97, 98, 990, 0}
	frame.NavMessage = NavSelect
	frame.NavMessageNotFound = NavNotFoundBlue

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)
}

func create90a(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)
	frame.PID.PageNumber = 90
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.FrameType = globals.FRAME_TYPE_EXIT
	//frame.Content.Data = ""

	frame.Content.Data = "[b][m-]\r\n      Thankyou for using[Y]TELSTAR.\r\n    You were connected to[Y][SERVER].\r\n\r\n[b][m-]\r\n                   \u001bBT                                    \u001bBT\u001bAE\u001bFL                                \u001bBT\u001bAE\u001bFL\u001bDS\u001bGT                            \u001bBT\u001bAE\u001bFL\u001bDS\u001bGT\u001bEA\u001bCR                            \u001bFL\u001bDS\u001bGT\u001bEA\u001bCR                                \u001bGT\u001bEA\u001bCR                                    \u001bCR                                                           \r\n"
	frame.Content.Type = globals.CONTENT_TYPE_MARKUP
	frame.RoutingTable = []int{900, 900, 900, 900, 900, 900, 900, 900, 900, 900, 9}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create91a(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 91
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Carousel = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIEEXdn07suXlp3Z0FTLz6IKGHPlQIECBAgQIECACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAYMS54fzJmix4-YCDToOLOjyZ0WLSkzo6AkcGHuZUcRHlB4dgyAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgGDP9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9_YNCho9zImSoAqBGoAp0FUy8-iChhz5UCA4MPEuZYwTAFzAFg1AgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgGDYCAAoACACeQHkBdYTJlixIkSKlSJEoUKIEBQABAAQAEABYN_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_39g4AgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgGDn9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9_YsAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIBix_f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f2LICAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAYs_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_39i0AgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgAIACAAgGLX9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9__f_3_9_Ytq-jT0yg7OXZs39w0Pzh3Ao_LLl3BZuHPl3dMIGllyBIWzrlLmkKJGTSJUycsoUqlZJYtXLzLBiyZlWjVs3IuHLp2UePXz9AgQokaBIlTJ0ChSqVoFi1cvQMGLJmgaNWzdA4cunaB49fP0ECDChoIkWNHQSJMqWgmTZ09BQo0qaCpVrV0FizatoLl29fQYMOLGgyZc2dBo06taDZt3b0HDjy5oOnXt3QePPr2g-ff38pgw4sZHJlzZyujTq1ktm3dvNcOPLmW6de3cn48-vZf59_fwZiHv3Y8uHYIjbMPPQDVCxcLf4E0-mXDk8mI-_dlFCn5a9_QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 91}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create91b(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 91
	frame.PID.FrameId = "b"
	frame.Visible = true
	frame.Carousel = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	frame.Content.Data = "http://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQICn_OmTJkSJApUqVKlSpUlcHApiTJkyZEiRIkyZMmTIkSJFmKf3JX58____RBg-fPnzog__2iBBo-fNX9qgQcOHDh8-fGiAp_zFf______2v7_____7X___oEGr__Qa0ur_________QoCn9ywK__6D__a__7dH__tf___1QIv_9AiQIkaPX_bo0aMpiKf_-Yr__oP_9r__tUH_-1____9og__yjnz5cldX9qUc-fPkp__uSv_-g__2v_-1Qf_7X__aIv____KZv_dIV__2pTN___yn__mK___z__a____5__tf_9qgVf__8o5_vWJX__alHP___Kf_7kr____6dBr_____-h__2qBBr__ymbPnSFf69CUzf__8pnz5iq9ejRoECBWvXr0KBGjQYOHBHz4cECBAgwcCily7__yqDhw4ePnzogQeP7RAg0fPiBB___9X___QYPH___-sCmb__KoP________3B__tUCDV__tUH___1f36XB-_______KOf_8qgRo__9B__v9X___QINX___Qav__B8-NNX_-jR___8pm__ymYrq______h_odX9qg1f___7q__9X_-h1f_5TBgwYM37__KOSur-_Ro0Pr_8_f_qDV_b___r__1f_6D___lFKlSpUu__8pmK_v7VAgwf______tNX9r______V__oP__-gQaPnwpm__yqDB__tUCD__QIEG___1f2qL____9X_-g____z5-__ymb__KoNX_-lQKl6VAgQav6_V__oNf___1f2qBX_____9-hKOf_8qgVJ0JRw4cOPLkqiQoESNGgVL0aFEjQoECJGjRoymDNm__ymTJkzZv3_____-3bt27du3Thw8efPnz58-fPnz58-fP___QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 91}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create91c(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
		sb      = strings.Builder{}
	)

	frame.PID.PageNumber = 91
	frame.PID.FrameId = "c"
	frame.Visible = true
	frame.Carousel = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	//                "0123456789012345678901234567890123456789"
	sb.WriteString("0---------1---------2---------3-------00")
	sb.WriteString("0---------1---------2---------3-------01")
	sb.WriteString("\b\t0---------1---------2---------3-------02")
	sb.WriteString("0---------1---------2---------3-------03\r")
	sb.WriteString("0---------1---------2---------3-------04\n")
	sb.WriteString("\v0---------1---------2---------3-------05\r\n\v")
	sb.WriteString("0---------1---------2---------3-------06\n\r\v")
	sb.WriteString("0---------1---------2---------3-------07\b\t") // Note: '\t' is a cursor movement not a character
	sb.WriteString("0---------1---------2---------3-------08\t\b")
	sb.WriteString("\x1e\v\v\v\v\v\v\v\v\v\v\v\v\v\v\v0---------1---------2---------3-------09")
	sb.WriteString("0---------1---------2---------3-------10")
	sb.WriteString("\x1E\n\n\n\n\n\n\n\n\n\n\n0---------1---------2---------3-------11")
	sb.WriteString("0---------1---------2---------3-------12\r\t\b")
	sb.WriteString("0---------1---------2---------3-------13\r\b3")
	sb.WriteString("\t\b0---------1---------2---------3-------14")
	sb.WriteString("0---------1---------2---------3-------15\n\n\r\v\v")
	sb.WriteString("0---------1---------2---------3-------16\r\n\n\v\v")
	sb.WriteString("0---------1---------2---------3-------17\r\n\r\v\r\n\v")
	sb.WriteString("0---------1---------2---------3-------18\r\r\r\n\r\r\r\n\v\v")
	sb.WriteString("0---------1---------2---------3-------19")
	sb.WriteString("0---------1---------2---------3-------20")
	sb.WriteString("0---------1---------2---------3-------21")
	sb.WriteString("0---------1---------2---------3-------22")
	sb.WriteString("0---------1---------2---------3-------23")

	frame.Content.Data = sb.String()
	frame.Content.Type = globals.CONTENT_TYPE_RAW
	frame.RoutingTable = []int{910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 91}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create91d(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 91
	frame.PID.FrameId = "d"
	frame.Visible = true
	frame.Carousel = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	//                "0123456789012345678901234567890123456789"

	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgZOWGFAgQIECAodLmkCBAgQIECBAgQIECBAgQEkCDwk_f_n9KgQIECBAgQICh0uaatWqBAgQPWqBAgQIECBASQIkOr___t0CBAgQIECBAgKHS5pq1a92vNq1bte7VugQIEBJBqa6P___8-fGiBAgQIECAodLmuvVr1adWrVq16tGqBASQJMP_______9KgQIECBAgQIAh0OgQIECBAgQIECBAgQIECA0SQfNX______79AwQIECBAgKHQ1TLz6IN-ZBzy8MPLD0y5DRJAg______3QIECBAgQIECAodDZ-WHho04-aDvp6aECBAaJIFCR_______ogQIECBAgQICh0Ni2YefTpo5b-ufQgQIDRLgwQIkm______7BAgQIECBAgKHQ3TL46LkCBAgQIEBolg-bP____YIPn_____-qBAgQIECAodDIECBAgQIECBAgNF8HL2S3____6gTp1_____9ECBAgQICh0MgQIECBAgQIECA0XTNEBJEjx70KBAlVf_____-sECBAgKHQyBAgQIECBAgQIDRf8hQIECBA0JIECBBr______1QIECAodDIECBAgQIECBAgQGi-NqgQIEDUkg5-Pnr______9ggQICh0MgQIECBAgQIEBougwN0FKfJQNSSBDr_________9PxAgKHQyBAgQIECBAXGmnidCgwcHCxMSQIKstBXiwakiLSQIECAodLoECBAgQIECBAgwb0CDA4TI0KAkgwfP___________foCh0uaQIECBAgQIECBEjWLEaBAgQICSr_v__________71CgKHQ5pAgQIECBAgQIECBAgQEkCBAgQIECv3_______8-fECAodDGkCBAgQIECBAgQIECBASYIECBBo-f_________9ugQICh0MaQIECBAgQIECBAgQIECBAgJIOH____r16_GjRoUCBAgKHQxpAgQIECBAgQIECBAgQIECAlufr0adCgQIECBAgQIECAodQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQICh1AgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 91}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create91e(settings config.Config) error {

	//TODO make this more attractive
	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 91
	frame.PID.FrameId = "e"
	frame.Visible = true
	frame.Carousel = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	//                "0123456789012345678901234567890123456789"

	frame.Content.Data = "\r\n\x1b\x4dTELSTAR\r\n*******\r\n\r\n" +
		" This frame should show the word\r\n" +
		" TELSTAR in double height characters.\r\n\r\n" +
		" The frame attempts to write asterisks\r\n" +
		" to the lower row, these should not\r\n" +
		" be visible."

	frame.Content.Type = globals.CONTENT_TYPE_RAW
	frame.RoutingTable = []int{910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 91}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}
func create91f(settings config.Config) error {

	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 91
	frame.PID.FrameId = "f"
	frame.Visible = true
	frame.Carousel = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAagQIECBAgQRN_XFsyoJGXTn0dEEflh4aNOPmgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBkgQM0CBAgQMkCBmgQIECBsgQN0CBAgQNkCBugQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIEDAaXQIC7AYHcDS6hAXcDA7AaXwIEHAYHcDS-hAg8IECBAgQIECBAgX3ECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAxGkUKAixGB3I0ipQEXIwOxGkcKBBxGB3I0jpQIPKBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQMhpNEgJshgeCNJokBN0MDshpPEgQchgeCNJ6kCD0gQIGaAmgQE0CBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIEDMaSRoCTMYHhDSStASdjA7MaSRoEHMYHhDSWtAg9oECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIDSBA0GlEiAo0GB4Y0osQFHgwO0GlMiBB0GB4Y0psQIPiBAgbjSaBAgQIECBAgQIECBAUQIECBAgQIECBAgQIECBAgQIECA0gQNRpVCgKtRgeINKrUBV6MDtRpXKgQdRgeINK7UCD6gQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgNIEDYaWTICzYYHijSy5AWfDA7YaWzIEHYYHijS25Ag_IECBygJoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIDSBA3Gl06Au3GB4w0uvQF34wO3Gl86BB3GB4w0vvQIP6BAgQICaBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 91}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)
}

func create94a(settings config.Config) error {

	//TODO make this more attractive
	var (
		frame   types.Frame
		primary bool
	)
	frame.PID.PageNumber = 94
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.FrameType = globals.FRAME_TYPE_TEST
	//                "0123456789012345678901234567890123456789"

	frame.Content.Data = "[D]SYSTEM INFO\r\n\r\n[g][m.]\r\n\r\n[SYSINFO]"

	frame.Content.Type = globals.CONTENT_TYPE_MARKUP
	frame.RoutingTable = []int{940, 941, 942, 943, 944, 945, 946, 947, 948, 949, 9}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create96a(settings config.Config) error {
	return nil
}

func create98a(settings config.Config) error {
	var (
		frame   types.Frame
		primary bool
	)

	frame.PID.PageNumber = 98
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.FrameType = globals.FRAME_TYPE_INFORMATION
	frame.Content.Data = "https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgJIECBA3_69WvVjdtW_P3qxs2vfG1b8_erW3aoECBAgQICJ0ugQIEDX_q1atWpq1a8_-rU1a_9TVrz_6szlqgQIECBAgQICSBAgQfeH7V61cvLXrz96uXFr_9NevP3q9dWqBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIASCpFmU6kGkg080GFB205e-TD0woOfnn0y7UHTRh6IOHLetBIO2nJl5oMKDhsw9M2_ltQc9O3Tsw8kHTeg6aN_PKg6aMPQEg75eWVBw5Ze2XZl3dEGTry07s6DpoyoGLluw5oMO7IgQIASBi5cMOa4PTy5UGjp04c3S9fn2YefPp08rse_agQIECBAgBIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECACgqaMqDn559Mu1Bjw7kGLKgw48eXnzy5EHXnp3Z1y5AgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgBIFQeDj38tyChyy8-mXYgpT5qBRChQ0E3Tj5b1KBAgQIECACgVB4e_bt59MPJAohQoaCbpx8t6lAgQIECBAgQIECBAgQIASBUHg7umXZv4ZUCiDw4bMqCTJUoECBAgQIECBAgQIECBAgAoFQePh9ZdmzKgUQeHDZlQSZKlAgQIECBAgQIECBAgQIECAEgVB5GXDy54fKCTJQKKWnnjQT6alAgQIECBAgQIECBAgQIAKBUHhQoaCFh56caCpl2c-mHkgh7NOXd0QKNvXZ008NmVSgBIFQeFyw7snlBCw89ONBD2acu7ogUbeuzpp4bMqlAgQIECACgVB5uHpy0-EELlh3ZPKCHs05d3RAo29dnTTw2ZVKBAgQIASBUH7ZO2nL3XYfXXll56emXmu3ZeiBR3y4lKBAgQIECBAgAoFQepl2c-mHkgracvfJh6YUEPZpy7uiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA"
	frame.Content.Type = globals.CONTENT_TYPE_EDITTF
	frame.RoutingTable = []int{980, 981, 982, 983, 984, 985, 986, 987, 988, 989, 9}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	return dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

}

func create99a(settings config.Config) error {

	var (
		err     error
		frame   types.Frame
		primary bool
	)

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	frame.PID.PageNumber = 99
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Content.Type = "markup"

	var sb strings.Builder
	sb.WriteString("[W][GREETING]    [G][DATE]   [TIME]\r\n")
	sb.WriteString("   [W][NAME]\n\n\r\n")
	sb.WriteString("[b][h-]")
	sb.WriteString("[g[ Welcome to Telstar]]") // alphagraphics delineated with [[ and ]] small g = graphics green
	sb.WriteString("[b][h-]")                  // horizontal line
	sb.WriteString("\n\n\r\n")
	sb.WriteString("      [G]YOU ARE CONNECTED TO [SERVER]\r\n")
	sb.WriteString("      [G]     [DATE] [TIME]\r\n")
	sb.WriteString("\n\n\r\n")
	//sb.WriteString("[G]   PRESS[W]1[G]TO LOG IN,[W]" + string(globals.HASH) + "[G]TO CONTINUE")

	frame.Content.Data = sb.String()
	frame.FrameType = globals.FRAME_TYPE_INITIAL
	frame.Cursor = false
	frame.RoutingTable = []int{0, 990, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//frame.NavMessage = "[G]   PRESS[W]1[G]TO LOG IN,[W]" + string(globals.HASH) + "[G]TO CONTINUE"
	frame.NavMessage = "[G]   PRESS[W]" + string(globals.HASH) + "[G]TO CONTINUE"
	frame.NavMessageNotFound = NavNotFoundBlue

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

	return err
}

func create990a(settings config.Config) error {

	var (
		err     error
		frame   types.Frame
		primary bool
	)

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	frame.PID.PageNumber = 990
	frame.PID.FrameId = "a"
	frame.Visible = true
	frame.Content.Type = "markup"

	var sb strings.Builder
	sb.WriteString("[W][GREETING]    [G][DATE]   [TIME]\r\n")
	sb.WriteString("   [W][NAME]\n\n\r\n")
	sb.WriteString("[b][h-]")
	sb.WriteString("")
	sb.WriteString("[g[Welcome to Telstar]]") // alphagraphics delineated with [[ and ]] small g = graphics green
	sb.WriteString("[b][h-]")                 // horizontal line
	sb.WriteString("\n\n\r\n")
	sb.WriteString("      [G]YOU ARE CONNECTED TO [SERVER]\r\n")
	sb.WriteString("      [G]     [DATE] [TIME]\r\n")
	sb.WriteString("\n\r\n")
	sb.WriteString("  USER ID:\n\r\n")
	sb.WriteString(" PASSWORD:")

	frame.Content.Data = sb.String()
	frame.FrameType = globals.FRAME_TYPE_RESPONSE
	frame.Cursor = true
	frame.RoutingTable = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	frame.ResponseData.Fields = []types.ResponseField{
		types.ResponseField{18, 11, true, 10, "numeric", true, false},
		types.ResponseField{20, 11, true, 10, "numeric", false, true},
	}
	frame.ResponseData.Action.Exec = "telstar.login" // internal command/app
	frame.ResponseData.Action.PostActionFrame.PageNumber = 99
	frame.ResponseData.Action.PostActionFrame.FrameId = "a"

	// if these fields are left empty e.g. Page Number 0 and FrameId "" y=then
	// this will be ammended to 0a during the routing process which means that
	// these could just be left blank.
	frame.ResponseData.Action.PostCancelFrame.PageNumber = 0
	frame.ResponseData.Action.PostCancelFrame.FrameId = "a"

	frame.NavMessageNotFound = NavNotFoundBlue

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY
	err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

	return err
}

func create101Page(settings config.Config) error {

	// FIXME With merged pages that are 960 char long, rendering the navigation message
	//  causes a scroll, this shouldn't happen if HOME/VTAB is used, should it?

	return nil

	var (
		mergedData string
		err        error
		frame      types.Frame
		primary    bool
	)

	// This page is used purely to test various page options and only on the secondary
	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	//blobData, err := blob.LoadBlob("./tmp/blobs/logo.blob")
	if err != nil {
		return err
	}
	blob1, err := convert.EdittfToRawT("https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAxYMcKAbNw6dyCTuyZfCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgKJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRAxQR4s6LSgzEEmdUi0otOogkzo0-lNg1JM-cCg8unNYgQIASBBj67OnXllWIMu7pl5dMOndty7uiDDuyINu_llXLlyBAgQMkE6LXpoIM6IgrxYNSRFpAp2Hpp37sOxBh3ZECBAgQIECAEgQad3TLy3Yemnfuw7EG7L35oMO7Ig75cPTRl5LkCBAgQIEDNBCq05M6LTpoJM6NPpTYNSTPnAoe_bww7vKDDuyIECBAgBIEGLrz07svPmg3Ze_NcDOhI1SnFQTMPTLz6IKHLTjy80CBA0QR59aLSnTYs6ogkzo0-lNg1JM-cgQIECBAgQIECBAgQIEDVBUpQY0aTDQUotCfSqUwU3Dq38kHLfhyIOWXhv5dOaBAgBIEGblv2oJGnPo74fPNBF3Z9mHdkXIECBAgQIECBAgQIECBA2QR4NSLXg2UFOLSrSYcWmCnWKkWYsQVIsWNBsLEEOHMaIASBBh3ZEHTRlQQ9-zfz54diCHh7ZUGHJ2y7unXllXIECBAgQN0EifMkxINmmCqcsPbLsQYd2RBI37NOTD5QbsvfmuQIECBA4CHQc2TDpT50WogcMGCAScgoOnLTi69MqDpvQdNGnmgQIAaBBzy8u2nHlQd9PTQIqZdmXnvzdO-HllE4d2RBt38sq5AgQOQ0yfHQT40ZYgp2adSLNQSZ0aegTIKkWnUQUIMeLTQIECAokSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEFeLBqSItJBGn0osODTqIBIM6EQIEFDDnyoFTJywvoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA")
	blob2, err := convert.EdittfToRawT("https://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECAv8-dEnz58-fPiBAgQIECBAgQIECBAgQIECBAgQIECBAgQIC7X5_wMlvz581IECBAgQIECBAgQIECBAgQIECBAgQIECBAgLtf6VP47av__UgQIECBAgQIECBAgQIECBAgQIECBAgQIECAu1_tOmjp25u_6BAgQIECBAgQIECBAgQIECBAgQIECBAgQIC7X__df_7X583oECBAgQIECBAgQIECBAgQIECBAgQIECBAgLtf_9w734f__UgQIECBAgQIECBAgQIECBAgQIECBAgQIECAv158__Zam_8_SBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA")
	//blob1 := convert.Blob{blobData, 0, 0}
	//blob2 := convert.Blob{blobData,0, 11}

	if mergedData, err = convert.RawTMerge(blob1, blob2); err != nil {
		return err
	}

	primary = strings.ToLower(settings.Database.Collection) == globals.DBPRIMARY

	frame.PID.PageNumber = 101
	frame.PID.FrameId = "a"
	frame.Content.Type = "rawT"
	frame.Content.Data = mergedData
	frame.FrameType = globals.FRAME_TYPE_INFORMATION
	frame.Cursor = false
	frame.RoutingTable = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	frame.NavMessage = NavSelect
	frame.NavMessageNotFound = NavNotFoundBlue
	frame.Visible = true

	err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primary)

	return err
}
