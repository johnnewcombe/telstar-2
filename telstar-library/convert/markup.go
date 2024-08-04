package convert

import (
	"github.com/johnnewcombe/telstar-library/globals"
	"regexp"
	"strings"
)

func MarkupToRawV(markup string) (string, error) {

	var (
		//length int
		alphaGraphics       string
		err                 error
		alphaGraphicsMarkup []string
		colour              string
	)

	markup = getControlMarkup(markup)
	alphaGraphicsMarkup = getAlphaGraphicMarkup(markup)

	for _, element := range alphaGraphicsMarkup {

		// looking for the colour indicator
		// alpha graphic syntax is set as [g[Welcometo Telstar]], g = graphic green
		re := regexp.MustCompile(`^\[[rygbcmw]\[`)

		if re.MatchString(element) {

			colour = replaceNth(element[:3], "[", "]", 2)
			colour = getControlMarkup(colour)

			// remove the colour markup from the string
			res := re.ReplaceAll([]byte(element), []byte("[["))
			elementTrimmed := string(res)

			// trim the ends otherwise these will get alpha-graphic-ed also
			elementTrimmed = strings.Trim(elementTrimmed, "[")
			elementTrimmed = strings.Trim(elementTrimmed, "]")

			// htab and vtab will always be zero as the padding will be expected in the markup string
			if alphaGraphics, err = TextToAlphagraphics(elementTrimmed, colour, 0, 0); err != nil {
				return "", err
			}

			// must move cursor down 4 rows after alphagraphics
			markup = strings.ReplaceAll(markup, element, alphaGraphics)
		}
	}

	return markup, nil
}

func getControlMarkup(markup string) string {

	// NOTE,NOTE,NOTE Any changes here MUST be reflected in utils.GetMarkupLen()
	markup = strings.ReplaceAll(markup, "[R]", globals.ALPHA_RED)
	markup = strings.ReplaceAll(markup, "[G]", globals.ALPHA_GREEN)
	markup = strings.ReplaceAll(markup, "[Y]", globals.ALPHA_YELLOW)
	markup = strings.ReplaceAll(markup, "[B]", globals.ALPHA_BLUE)
	markup = strings.ReplaceAll(markup, "[M]", globals.ALPHA_MAGENTA)
	markup = strings.ReplaceAll(markup, "[C]", globals.ALPHA_CYAN)
	markup = strings.ReplaceAll(markup, "[W]", globals.ALPHA_WHITE)
	markup = strings.ReplaceAll(markup, "[F]", globals.FLASH)
	markup = strings.ReplaceAll(markup, "[S]", globals.STEADY)
	markup = strings.ReplaceAll(markup, "[N]", globals.NORMAL_HEIGHT)
	markup = strings.ReplaceAll(markup, "[D]", globals.DOUBLE_HEIGHT)
	markup = strings.ReplaceAll(markup, "[-]", globals.END_BACKGROUND)
	markup = strings.ReplaceAll(markup, "[n]", globals.NEW_BACKGROUND)
	markup = strings.ReplaceAll(markup, "[r]", globals.MOSAIC_RED)
	markup = strings.ReplaceAll(markup, "[g]", globals.MOSAIC_GREEN)
	markup = strings.ReplaceAll(markup, "[y]", globals.MOSAIC_YELLOW)
	markup = strings.ReplaceAll(markup, "[b]", globals.MOSAIC_BLUE)
	markup = strings.ReplaceAll(markup, "[m]", globals.MOSAIC_MAGENTA)
	markup = strings.ReplaceAll(markup, "[c]", globals.MOSAIC_CYAN)
	markup = strings.ReplaceAll(markup, "[w]", globals.MOSAIC_WHITE)
	markup = strings.ReplaceAll(markup, "[h.]", globals.SEPARATOR_GRAPHIC_DOTS_HIGH)
	markup = strings.ReplaceAll(markup, "[m.]", globals.SEPARATOR_GRAPHIC_DOTS_MID)
	markup = strings.ReplaceAll(markup, "[l.]", globals.SEPARATOR_GRAPHIC_DOTS_LOW)
	markup = strings.ReplaceAll(markup, "[h-]", globals.SEPARATOR_GRAPHIC_SOLID_HIGH)
	markup = strings.ReplaceAll(markup, "[m-]", globals.SEPARATOR_GRAPHIC_SOLID_MID)
	markup = strings.ReplaceAll(markup, "[l-]", globals.SEPARATOR_GRAPHIC_SOLID_LOW)
	markup = strings.ReplaceAll(markup, "[=]", globals.SEPARATOR_GRAPHIC_SOLID_DOUBLE)
	markup = strings.ReplaceAll(markup, "[_+]", string(globals.CURON))
	markup = strings.ReplaceAll(markup, "[_-]", string(globals.CUROFF))
	markup = strings.ReplaceAll(markup, "[@]", string(globals.HOME))
	markup = strings.ReplaceAll(markup, "[H]", string(globals.HT))
	markup = strings.ReplaceAll(markup, "[V]", string(globals.VT))

	// TODO implement markup that includes spaces e.g. [12] = 12 space (0x20) characters
	//  any changes here need to be reflected in the GetMarkupLen() Util method
	//for n:=1; n<globals.COLS;n++{
	//}

	return markup

}

func getAlphaGraphicMarkup(text string) []string {

	// e.g. "[g[Welcome to Telstar]]" = graphic green alpha graphics
	re := regexp.MustCompile(`\[[rygbcmw]\[([^\[\]]*)\]\]`)
	if re.MatchString(text) {
		return re.FindAllString(text, -1)
	}
	return []string{}

}

// Replace the nth occurrence of old in s by new.
func replaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}
