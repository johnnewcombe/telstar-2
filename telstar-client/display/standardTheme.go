package display

import (
	"bitbucket.org/johnnewcombe/telstar-client/constants"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

var Gray = color.Gray{Y: 24}

type StandardTheme struct{}

func (m StandardTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	//gray := color.Gray{Y: 24}

	if name == theme.ColorNameBackground {
		//if variant == theme.VariantDark {
		//	return color.White
		//} else {
		return Gray
		//}
	} else if name == theme.ColorNameForeground {
		//if variant == theme.VariantDark {
		//	return gray
		//} else {
		return color.White
		//}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m StandardTheme) Font(style fyne.TextStyle) fyne.Resource {

	if style.Monospace {
		return constants.MODE7GX2TTF
	} else {
		return theme.DefaultTheme().Font(style)
	}
}

func (m StandardTheme) Icon(name fyne.ThemeIconName) fyne.Resource {

	if name == "serial" {
		return constants.SerialIcon
	} else if name == "cloud" {
		return constants.CloudIcon
	} else{
		return theme.DefaultTheme().Icon(name)
	}
}

func (m StandardTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

